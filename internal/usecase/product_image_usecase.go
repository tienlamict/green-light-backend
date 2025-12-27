package usecase

import (
	"context"
	"fmt"

	"green-light-backend/internal/domain"
	"green-light-backend/pkg/storage"
	"green-light-backend/pkg/utils"
)

// ProductImageUseCase handles business logic for product images
type ProductImageUseCase struct {
	imageRepo   domain.ProductImageRepository
	productRepo domain.ProductRepository
	variantRepo domain.ProductVariantRepository
	minioClient *storage.MinIOClient
}

// NewProductImageUseCase creates a new product image use case
func NewProductImageUseCase(
	imageRepo domain.ProductImageRepository,
	productRepo domain.ProductRepository,
	variantRepo domain.ProductVariantRepository,
	minioClient *storage.MinIOClient,
) *ProductImageUseCase {
	return &ProductImageUseCase{
		imageRepo:   imageRepo,
		productRepo: productRepo,
		variantRepo: variantRepo,
		minioClient: minioClient,
	}
}

// GeneratePresignedUploadURLInput is the input for generating presigned URL
type GeneratePresignedUploadURLInput struct {
	ProductID   string
	ContentType string
	Extension   string
}

// GeneratePresignedUploadURL generates a presigned URL for direct upload to MinIO
func (uc *ProductImageUseCase) GeneratePresignedUploadURL(
	ctx context.Context,
	input GeneratePresignedUploadURLInput,
) (*storage.PresignedURLResponse, error) {
	// Verify product exists
	product, err := uc.productRepo.GetByID(ctx, input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Generate unique UUID for the image (using v7 for time-ordered consistency)
	imageUUID := utils.GenerateUUIDv7()

	// Generate object key: products/{yyyy}/{mm}/{product_id}/{uuid}.{ext}
	objectKey := storage.GenerateObjectKey(input.ProductID, imageUUID, input.Extension)

	// Generate presigned URL (max 7MB)
	maxSize := int64(7 * 1024 * 1024) // 7MB
	presignedResp, err := uc.minioClient.GeneratePresignedUploadURL(
		ctx,
		objectKey,
		input.ContentType,
		maxSize,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedResp, nil
}

// ConfirmImageUploadInput is the input for confirming image upload
type ConfirmImageUploadInput struct {
	ProductID string
	VariantID *string // Optional: if provided, image belongs to variant
	ObjectKey string
	IsMain    bool
	SortOrder int
}

// ConfirmImageUpload confirms that an image was uploaded and saves metadata to DB
func (uc *ProductImageUseCase) ConfirmImageUpload(
	ctx context.Context,
	input ConfirmImageUploadInput,
) (*domain.ProductImage, error) {
	// Verify product exists
	product, err := uc.productRepo.GetByID(ctx, input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// If variant_id is provided, verify variant exists and belongs to product
	if input.VariantID != nil && *input.VariantID != "" {
		variant, err := uc.variantRepo.GetByID(ctx, *input.VariantID)
		if err != nil {
			return nil, fmt.Errorf("variant not found: %w", err)
		}
		if variant == nil {
			return nil, fmt.Errorf("variant not found")
		}
		// Verify variant belongs to the product
		if variant.ProductID != input.ProductID {
			return nil, fmt.Errorf("variant does not belong to this product")
		}
	}

	// Verify object exists in MinIO
	exists, err := uc.minioClient.ObjectExists(ctx, input.ObjectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to verify object existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("object not found in storage")
	}

	// Build public URL
	publicURL := fmt.Sprintf("%s/%s", uc.minioClient.PublicURL(), input.ObjectKey)

	// Create image record
	image := &domain.ProductImage{
		ImageID:   utils.GenerateUUIDv7(),
		ProductID: input.ProductID,
		VariantID: input.VariantID, // Set variant_id if provided (NULL for product-level images)
		URL:       publicURL,
		ObjectKey: input.ObjectKey,
		IsMain:    input.IsMain,
		SortOrder: input.SortOrder,
		Status:    "ACTIVE",
	}

	// If this is the main image, unset other main images
	if input.IsMain {
		if err := uc.imageRepo.SetMainImage(ctx, input.ProductID, image.ImageID); err != nil {
			return nil, fmt.Errorf("failed to set main image: %w", err)
		}
	}

	// Save to database
	if err := uc.imageRepo.Create(ctx, image); err != nil {
		return nil, fmt.Errorf("failed to save image metadata: %w", err)
	}

	return image, nil
}

// GetProductImages returns all images for a product
func (uc *ProductImageUseCase) GetProductImages(ctx context.Context, productID string) ([]*domain.ProductImage, error) {
	// Verify product exists
	product, err := uc.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	images, err := uc.imageRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product images: %w", err)
	}

	return images, nil
}

// DeleteProductImage deletes an image from both MinIO and database
func (uc *ProductImageUseCase) DeleteProductImage(ctx context.Context, imageID string) error {
	// Get image metadata
	image, err := uc.imageRepo.GetByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("image not found: %w", err)
	}

	// Delete from MinIO
	if err := uc.minioClient.DeleteObject(ctx, image.ObjectKey); err != nil {
		// Log error but continue to delete from DB
		fmt.Printf("Warning: failed to delete object from MinIO: %v\n", err)
	}

	// Delete from database
	if err := uc.imageRepo.Delete(ctx, imageID); err != nil {
		return fmt.Errorf("failed to delete image metadata: %w", err)
	}

	return nil
}

// UpdateImageMetadata updates image metadata (is_main, sort_order)
type UpdateImageMetadataInput struct {
	ImageID   string
	IsMain    *bool
	SortOrder *int
}

func (uc *ProductImageUseCase) UpdateImageMetadata(
	ctx context.Context,
	input UpdateImageMetadataInput,
) (*domain.ProductImage, error) {
	// Get existing image
	image, err := uc.imageRepo.GetByID(ctx, input.ImageID)
	if err != nil {
		return nil, fmt.Errorf("image not found: %w", err)
	}

	// Update fields if provided
	if input.IsMain != nil {
		if *input.IsMain {
			// Set as main image (will unset others)
			if err := uc.imageRepo.SetMainImage(ctx, image.ProductID, input.ImageID); err != nil {
				return nil, fmt.Errorf("failed to set main image: %w", err)
			}
		}
		image.IsMain = *input.IsMain
	}

	if input.SortOrder != nil {
		image.SortOrder = *input.SortOrder
	}

	// Save changes
	if err := uc.imageRepo.Update(ctx, image); err != nil {
		return nil, fmt.Errorf("failed to update image metadata: %w", err)
	}

	return image, nil
}

