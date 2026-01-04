package storage

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"green-light-backend/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient wraps MinIO client with helper methods
type MinIOClient struct {
	client        *minio.Client // Client for internal operations
	presignClient *minio.Client // Client for generating presigned URLs (uses public endpoint)
	bucket        string
	publicURL     string
}

// PresignedURLResponse contains upload and public URLs
type PresignedURLResponse struct {
	UploadURL string `json:"upload_url"` // Presigned PUT URL for upload
	PublicURL string `json:"public_url"` // Public URL to access the object
	ObjectKey string `json:"object_key"` // MinIO object key
	ExpiresAt string `json:"expires_at"` // Expiration time
}

// NewMinIOClient creates a new MinIO client
func NewMinIOClient(cfg *config.MinIOConfig) (*MinIOClient, error) {
	// Initialize MinIO client for internal operations
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Determine presign endpoint
	// If PresignEndpoint is set, use it for generating presigned URLs
	// This allows backend to use internal endpoint (minio:9000) for operations
	// but generate presigned URLs with an endpoint accessible from frontend
	presignEndpoint := cfg.PresignEndpoint
	if presignEndpoint == "" {
		// If not set, extract host from PublicURL
		publicURLParsed, err := url.Parse(cfg.PublicURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse public URL: %w", err)
		}
		presignEndpoint = publicURLParsed.Host
	}

	// Initialize presign client with presign endpoint
	// This client is used ONLY to generate presigned URLs with correct hostname
	presignClient, err := minio.New(presignEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create presign MinIO client: %w", err)
	}

	// Check if bucket exists, create if not
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Set bucket policy to public-read (optional, for public access)
	// Note: For production, you might want more restrictive policies
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, cfg.BucketName)

	err = client.SetBucketPolicy(ctx, cfg.BucketName, policy)
	if err != nil {
		// Log warning but don't fail - policy might already be set
		fmt.Printf("Warning: failed to set bucket policy: %v\n", err)
	}

	return &MinIOClient{
		client:        client,
		presignClient: presignClient,
		bucket:        cfg.BucketName,
		publicURL:     cfg.PublicURL,
	}, nil
}

// GeneratePresignedUploadURL generates a presigned PUT URL for direct upload
// objectKey: full path in MinIO (e.g., "products/2024/12/product-id/uuid.webp")
// contentType: MIME type (e.g., "image/webp", "image/jpeg", "image/png")
// maxSize: maximum file size in bytes (e.g., 7340032 for 7MB)
func (m *MinIOClient) GeneratePresignedUploadURL(
	ctx context.Context,
	objectKey string,
	contentType string,
	maxSize int64,
) (*PresignedURLResponse, error) {
	// Validate content type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("invalid content type: %s (allowed: jpeg, png, webp)", contentType)
	}

	// Set expiration time (10 minutes)
	expiry := 10 * time.Minute
	expiresAt := time.Now().Add(expiry)

	// Generate presigned PUT URL using presignClient
	// The presignClient is configured with an endpoint that frontend can access
	presignedURL, err := m.presignClient.PresignedPutObject(
		ctx,
		m.bucket,
		objectKey,
		expiry,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Build public URL for accessing the uploaded object
	publicURL := fmt.Sprintf("%s/%s/%s", m.publicURL, m.bucket, objectKey)

	return &PresignedURLResponse{
		UploadURL: presignedURL.String(),
		PublicURL: publicURL,
		ObjectKey: objectKey,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

// DeleteObject deletes an object from MinIO
func (m *MinIOClient) DeleteObject(ctx context.Context, objectKey string) error {
	err := m.client.RemoveObject(ctx, m.bucket, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// ObjectExists checks if an object exists in MinIO
func (m *MinIOClient) ObjectExists(ctx context.Context, objectKey string) (bool, error) {
	_, err := m.client.StatObject(ctx, m.bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		// Check if error is "object not found"
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}
	return true, nil
}

// GenerateObjectKey generates a unique object key for product images
// Format: products/{yyyy}/{mm}/{product_id}/{uuid}.{ext}
func GenerateObjectKey(productID, uuid, ext string) string {
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")

	// Clean extension (remove dot if present)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	return filepath.Join("products", year, month, productID, fmt.Sprintf("%s.%s", uuid, ext))
}

// GenerateCategoryIconObjectKey generates a unique object key for category icons
// Format: categories/{yyyy}/{mm}/{category_id}/{uuid}.{ext}
func GenerateCategoryIconObjectKey(categoryID, uuid, ext string) string {
	now := time.Now()
	year := now.Format("2006")
	month := now.Format("01")

	// Clean extension (remove dot if present)
	if len(ext) > 0 && ext[0] == '.' {
		ext = ext[1:]
	}

	return filepath.Join("categories", year, month, categoryID, fmt.Sprintf("%s.%s", uuid, ext))
}

// GetObjectKeyFromURL extracts object key from public URL
func (m *MinIOClient) GetObjectKeyFromURL(publicURL string) (string, error) {
	// Parse URL
	u, err := url.Parse(publicURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Extract path and remove leading slash and bucket name
	path := u.Path
	prefix := fmt.Sprintf("/%s/", m.bucket)
	if len(path) > len(prefix) && path[:len(prefix)] == prefix {
		return path[len(prefix):], nil
	}

	return "", fmt.Errorf("invalid public URL format")
}

// PublicURL returns the base public URL for MinIO
func (m *MinIOClient) PublicURL() string {
	return fmt.Sprintf("%s/%s", m.publicURL, m.bucket)
}
