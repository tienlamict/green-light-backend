# Variant Images Support - Implementation Guide

## Tổng quan

Đã triển khai thành công tính năng hỗ trợ **nhiều ảnh cho mỗi variant**. Mỗi variant giờ đây có thể có một hoặc nhiều ảnh riêng, độc lập với ảnh của product.

## Kiến trúc Database

### Bảng `product_images`

```sql
CREATE TABLE `product_images` (
  `image_id` VARCHAR(36) NOT NULL,
  `product_id` VARCHAR(36) NOT NULL,
  `variant_id` VARCHAR(36) NULL,  -- ⭐ NULL = product image, NOT NULL = variant image
  `url` VARCHAR(500) NOT NULL,
  `object_key` VARCHAR(500) NOT NULL,
  `is_main` BOOLEAN NOT NULL DEFAULT FALSE,
  `sort_order` INT NOT NULL DEFAULT 0,
  `status` VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
  PRIMARY KEY (`image_id`),
  INDEX `idx_images_product_id` (`product_id`),
  INDEX `idx_images_variant_id` (`variant_id`),
  CONSTRAINT `fk_images_product` FOREIGN KEY (`product_id`) REFERENCES `products`(`product_id`),
  CONSTRAINT `fk_images_variant` FOREIGN KEY (`variant_id`) REFERENCES `product_variants`(`variant_id`)
);
```

### Cách hoạt động

1. **Product-level images** (ảnh chung cho tất cả variants):
   - `product_id`: có giá trị
   - `variant_id`: **NULL**
   - Ví dụ: Ảnh tổng quan sản phẩm, lifestyle shots

2. **Variant-specific images** (ảnh riêng cho từng variant):
   - `product_id`: có giá trị (để dễ query)
   - `variant_id`: **có giá trị**
   - Ví dụ: Ảnh màu đỏ cho variant màu đỏ, ảnh size L cho variant size L

## API Response Structure

### GET /api/v1/products (List Products)

```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "product_id": "01JGXXX...",
      "name": "iPhone 15 Pro",
      "variants": [
        {
          "variant_id": "01JGYYY...",
          "sku": "IP15P-128-BLK",
          "name": "128GB - Black",
          "price": 999.99,
          "stock": 50,
          "images": [
            {
              "image_id": "img-001",
              "url": "https://cdn.example.com/iphone-15-black-front.jpg",
              "is_main": true,
              "sort_order": 0
            },
            {
              "image_id": "img-002",
              "url": "https://cdn.example.com/iphone-15-black-back.jpg",
              "is_main": false,
              "sort_order": 1
            }
          ]
        },
        {
          "variant_id": "01JGZZZ...",
          "sku": "IP15P-128-WHT",
          "name": "128GB - White",
          "price": 999.99,
          "stock": 30,
          "images": [
            {
              "image_id": "img-003",
              "url": "https://cdn.example.com/iphone-15-white-front.jpg",
              "is_main": true,
              "sort_order": 0
            }
          ]
        }
      ]
    }
  ]
}
```

## Cách sử dụng

### 1. Thêm ảnh cho Product (shared images)

```bash
POST /api/v1/products/{product_id}/images
Content-Type: application/json

{
  "url": "https://cdn.example.com/product-main.jpg",
  "object_key": "products/abc123/main.jpg",
  "is_main": true,
  "sort_order": 0
  // variant_id: không truyền hoặc null
}
```

### 2. Thêm ảnh cho Variant (variant-specific images)

```bash
POST /api/v1/products/{product_id}/images
Content-Type: application/json

{
  "url": "https://cdn.example.com/variant-red.jpg",
  "object_key": "products/abc123/variants/red.jpg",
  "variant_id": "variant-001",  // ⭐ Chỉ định variant_id
  "is_main": true,
  "sort_order": 0
}
```

### 3. Query ảnh

#### a. Lấy ảnh của product (không bao gồm variant images):
```go
images, err := imageRepo.GetByProductID(ctx, productID)
// WHERE product_id = ? AND variant_id IS NULL
```

#### b. Lấy ảnh của một variant cụ thể:
```go
images, err := imageRepo.GetByVariantID(ctx, variantID)
// WHERE variant_id = ?
```

#### c. Lấy TẤT CẢ ảnh (product + tất cả variants):
```go
images, err := imageRepo.GetAllImagesByProductID(ctx, productID)
// WHERE product_id = ?
```

## Code Implementation

### Domain Model

```go
// internal/domain/product_image.go
type ProductImage struct {
    ImageID    string    `gorm:"primaryKey;type:varchar(36)"`
    ProductID  string    `gorm:"type:varchar(36);not null;index"`
    VariantID  *string   `gorm:"type:varchar(36);index"` // ⭐ Nullable pointer
    URL        string    `gorm:"type:varchar(500);not null"`
    ObjectKey  string    `gorm:"type:varchar(500);not null"`
    IsMain     bool      `gorm:"default:false;not null"`
    SortOrder  int       `gorm:"default:0;not null"`
    Status     string    `gorm:"type:varchar(20);default:'ACTIVE';not null"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    UpdatedAt  time.Time `gorm:"autoUpdateTime"`

    // Relationships
    Product *Product        `gorm:"foreignKey:ProductID"`
    Variant *ProductVariant `gorm:"foreignKey:VariantID"`
}

// internal/domain/product_variant.go
type ProductVariant struct {
    VariantID  string
    ProductID  string
    SKU        string
    Name       string
    Price      float64
    Stock      int
    // ...
    
    // Relationships
    Product *Product        `gorm:"foreignKey:ProductID"`
    Images  []*ProductImage `gorm:"foreignKey:VariantID"` // ⭐ Variant images
}
```

### Repository Interface

```go
type ProductImageRepository interface {
    Create(ctx context.Context, image *ProductImage) error
    GetByID(ctx context.Context, imageID string) (*ProductImage, error)
    
    // Get product-level images only (variant_id IS NULL)
    GetByProductID(ctx context.Context, productID string) ([]*ProductImage, error)
    
    // Get images for a specific variant
    GetByVariantID(ctx context.Context, variantID string) ([]*ProductImage, error)
    
    // Get all images (product + all variants)
    GetAllImagesByProductID(ctx context.Context, productID string) ([]*ProductImage, error)
    
    Update(ctx context.Context, image *ProductImage) error
    Delete(ctx context.Context, imageID string) error
    SetMainImage(ctx context.Context, productID, imageID string) error
    DeleteOrphanImages(ctx context.Context, olderThan time.Duration) error
}
```

### DTO Response

```go
// internal/transport/http/dto/product_variant_dto.go
type VariantResponse struct {
    VariantID  string            `json:"variant_id"`
    ProductID  string            `json:"product_id"`
    SKU        string            `json:"sku"`
    Name       string            `json:"name"`
    Price      float64           `json:"price"`
    Stock      int               `json:"stock"`
    Images     []ImageResponse   `json:"images,omitempty"` // ⭐ Variant images
    // ...
}

type ImageResponse struct {
    ImageID   string `json:"image_id"`
    URL       string `json:"url"`
    IsMain    bool   `json:"is_main"`
    SortOrder int    `json:"sort_order"`
}
```

### Handler Helper Function

```go
// internal/transport/http/handler/product_handler.go
func (h *ProductHandler) loadVariantsWithImages(c *gin.Context, productID string) []dto.VariantResponse {
    variants, err := h.variantUseCase.GetByProductID(c.Request.Context(), productID)
    if err != nil || len(variants) == 0 {
        return []dto.VariantResponse{}
    }

    // Load images for each variant
    for _, variant := range variants {
        images, err := h.imageRepo.GetByVariantID(c.Request.Context(), variant.VariantID)
        if err == nil && len(images) > 0 {
            variant.Images = images
        }
    }

    return dto.ToVariantListResponse(variants)
}
```

## Use Cases

### 1. Sản phẩm có nhiều màu sắc

**Product**: iPhone 15 Pro
- Variant 1: Black (có 3 ảnh màu đen)
- Variant 2: White (có 3 ảnh màu trắng)
- Variant 3: Blue (có 3 ảnh màu xanh)

Mỗi variant có ảnh riêng để khách hàng thấy chính xác màu sắc.

### 2. Sản phẩm có nhiều size

**Product**: Áo thun
- Variant 1: Size S (ảnh người mặc size S)
- Variant 2: Size M (ảnh người mặc size M)
- Variant 3: Size L (ảnh người mặc size L)

### 3. Kết hợp nhiều thuộc tính

**Product**: Giày thể thao
- Variant 1: Red - Size 40 (ảnh giày đỏ size 40)
- Variant 2: Red - Size 41 (ảnh giày đỏ size 41)
- Variant 3: Blue - Size 40 (ảnh giày xanh size 40)

## Migration

### Chạy migration

```bash
# Nếu database đã tồn tại, cần chạy ALTER TABLE
mysql -u root -p green_light_db < migrations/001_create_tables.sql

# Hoặc drop và tạo lại database
DROP DATABASE IF EXISTS green_light_db;
CREATE DATABASE green_light_db;
mysql -u root -p green_light_db < migrations/001_create_tables.sql
```

### Rollback (nếu cần)

```sql
-- Remove variant_id column
ALTER TABLE `product_images` DROP FOREIGN KEY `fk_images_variant`;
ALTER TABLE `product_images` DROP INDEX `idx_images_variant_id`;
ALTER TABLE `product_images` DROP COLUMN `variant_id`;
```

## Testing

### Test Case 1: Tạo product với variants và thêm ảnh

```bash
# 1. Tạo product với variants
POST /api/v1/products
{
  "name": "iPhone 15 Pro",
  "category_id": "cat-001",
  "variants": [
    {
      "sku": "IP15P-BLK",
      "name": "Black",
      "price": 999.99,
      "stock": 50
    },
    {
      "sku": "IP15P-WHT",
      "name": "White",
      "price": 999.99,
      "stock": 30
    }
  ]
}

# Response: product_id = "prod-001"
#           variant_id[0] = "var-001" (Black)
#           variant_id[1] = "var-002" (White)

# 2. Thêm ảnh cho variant Black
POST /api/v1/products/prod-001/images
{
  "variant_id": "var-001",
  "url": "https://cdn.example.com/black-front.jpg",
  "object_key": "products/prod-001/black-front.jpg",
  "is_main": true,
  "sort_order": 0
}

# 3. Thêm ảnh cho variant White
POST /api/v1/products/prod-001/images
{
  "variant_id": "var-002",
  "url": "https://cdn.example.com/white-front.jpg",
  "object_key": "products/prod-001/white-front.jpg",
  "is_main": true,
  "sort_order": 0
}

# 4. Lấy product để xem variants với images
GET /api/v1/products/prod-001

# Response sẽ bao gồm variants với images tương ứng
```

### Test Case 2: Query images

```sql
-- Lấy ảnh của product (không bao gồm variant images)
SELECT * FROM product_images 
WHERE product_id = 'prod-001' AND variant_id IS NULL;

-- Lấy ảnh của variant Black
SELECT * FROM product_images 
WHERE variant_id = 'var-001';

-- Lấy TẤT CẢ ảnh (product + all variants)
SELECT * FROM product_images 
WHERE product_id = 'prod-001'
ORDER BY variant_id IS NULL DESC, is_main DESC, sort_order;
```

## Performance Considerations

### N+1 Query Problem

Khi list products với nhiều variants, mỗi variant cần 1 query để load images:
- 10 products × 3 variants/product = 30 queries cho variants
- 30 variants × 1 query/variant = 30 queries cho images
- **Tổng: ~60 queries**

### Optimization Options

#### Option 1: Eager Loading (Recommended)
```go
// Load variants với images trong 1 query
db.Preload("Images").Find(&variants)
```

#### Option 2: Batch Loading
```go
// Load tất cả images cho nhiều variants cùng lúc
variantIDs := []string{"var-001", "var-002", "var-003"}
images, _ := imageRepo.GetByVariantIDs(ctx, variantIDs)
// Group images by variant_id
```

#### Option 3: Caching
```go
// Cache images trong Redis với TTL
cacheKey := fmt.Sprintf("variant:%s:images", variantID)
redis.Set(cacheKey, images, 5*time.Minute)
```

#### Option 4: Query Parameter
```go
// Cho phép client chọn có load images hay không
GET /api/v1/products?include_images=false
```

## Breaking Changes

### ❌ KHÔNG có breaking changes

- Tất cả API responses đều backward compatible
- Field `images` trong `VariantResponse` là optional (`omitempty`)
- Clients cũ có thể ignore field `images`
- Database migration chỉ thêm column mới (nullable)

## Summary

✅ **Đã hoàn thành:**
1. ✅ Thêm `variant_id` vào bảng `product_images` (nullable)
2. ✅ Cập nhật domain models (ProductImage, ProductVariant)
3. ✅ Cập nhật repository với 3 methods mới
4. ✅ Cập nhật DTO responses để bao gồm images
5. ✅ Cập nhật tất cả product APIs để trả về variant images
6. ✅ Tạo helper function `loadVariantsWithImages()`
7. ✅ Không có lỗi linter

✅ **Lợi ích:**
- Mỗi variant có thể có nhiều ảnh riêng
- Linh hoạt: có thể có ảnh chung (product-level) và ảnh riêng (variant-level)
- Dễ query: 3 methods riêng biệt cho các use cases khác nhau
- Backward compatible: không breaking changes

🎉 **Hệ thống giờ đã hỗ trợ đầy đủ variant images!**

