# Schema Changes - January 2025

## Tổng quan
Document này mô tả các thay đổi về database schema và API được thực hiện vào tháng 1/2025.

## Thay đổi 1: Xóa trường `gallery` khỏi bảng `products`

### Lý do
- Trường `gallery` (JSON array) đã được thay thế hoàn toàn bởi bảng `product_images`
- Việc quản lý ảnh qua bảng riêng (`product_images`) mang lại nhiều lợi ích:
  - Linh hoạt hơn trong việc quản lý metadata (is_main, sort_order, status)
  - Hỗ trợ cả ảnh product-level và variant-specific
  - Dễ dàng query và filter
  - Tích hợp tốt hơn với MinIO storage

### Thay đổi Database
**Migration file**: `migrations/002_remove_gallery_add_icon.sql`

```sql
ALTER TABLE products DROP COLUMN gallery;
```

### Thay đổi API

#### Product Request (Create/Update)
**Trước:**
```json
{
  "name": "Product Name",
  "gallery": ["url1", "url2"],
  ...
}
```

**Sau:**
```json
{
  "name": "Product Name",
  ...
}
```

**Note**: Để thêm ảnh cho product, sử dụng Product Images API:
- `POST /api/v1/products/{id}/images/presign` - Get presigned URL
- `POST /api/v1/products/{id}/images` - Confirm upload

#### Product Response
**Trước:**
```json
{
  "product_id": "...",
  "name": "...",
  "gallery": ["url1", "url2"],
  ...
}
```

**Sau:**
```json
{
  "product_id": "...",
  "name": "...",
  "images": [
    {
      "image_id": "...",
      "url": "...",
      "is_main": true,
      "sort_order": 0
    }
  ],
  ...
}
```

### Breaking Changes
⚠️ **Breaking Change**: Field `gallery` đã bị xóa khỏi:
- `CreateProductRequest`
- `UpdateProductRequest`
- `ProductResponse`

### Migration Guide
1. **Backend**: Chạy migration `002_remove_gallery_add_icon.sql`
2. **Frontend**: 
   - Xóa code xử lý field `gallery`
   - Sử dụng field `images` thay thế
   - Sử dụng Product Images API để upload/manage ảnh

---

## Thay đổi 2: Thêm trường `icon_url` vào bảng `categories`

### Lý do
- Cho phép mỗi category có một icon/logo riêng
- Icon được lưu trên MinIO tương tự như `thumbnail_url` của product
- Hữu ích cho UI/UX khi hiển thị danh sách categories

### Thay đổi Database
**Migration file**: `migrations/002_remove_gallery_add_icon.sql`

```sql
ALTER TABLE categories 
ADD COLUMN icon_url VARCHAR(500) NULL 
COMMENT 'Icon image URL stored in MinIO' 
AFTER description;
```

### Thay đổi API

#### Category Request (Create)
**Trước:**
```json
{
  "name": "Electronics",
  "slug": "electronics",
  "description": "Electronic devices",
  "is_active": true
}
```

**Sau:**
```json
{
  "name": "Electronics",
  "slug": "electronics",
  "description": "Electronic devices",
  "icon_url": "http://localhost:9000/greenlight/categories/electronics-icon.png",
  "is_active": true
}
```

#### Category Request (Update)
```json
{
  "icon_url": "http://localhost:9000/greenlight/categories/new-icon.png"
}
```

#### Category Response
**Trước:**
```json
{
  "category_id": "...",
  "name": "Electronics",
  "slug": "electronics",
  "description": "...",
  "is_active": true,
  "created_at": "...",
  "updated_at": "..."
}
```

**Sau:**
```json
{
  "category_id": "...",
  "name": "Electronics",
  "slug": "electronics",
  "description": "...",
  "icon_url": "http://localhost:9000/greenlight/categories/electronics-icon.png",
  "is_active": true,
  "created_at": "...",
  "updated_at": "..."
}
```

### Breaking Changes
✅ **Backward Compatible**: Field `icon_url` là optional, không ảnh hưởng đến code cũ

### Migration Guide
1. **Backend**: Chạy migration `002_remove_gallery_add_icon.sql`
2. **Frontend**: 
   - Thêm field `icon_url` vào form tạo/sửa category
   - Hiển thị icon trong danh sách categories (nếu có)
   - Sử dụng MinIO upload flow tương tự như product thumbnail

---

## Cách Upload Icon cho Category

### Bước 1: Upload file lên MinIO
Sử dụng MinIO client hoặc presigned URL (tương tự product images)

### Bước 2: Lưu URL vào category
```bash
POST /api/v1/categories
{
  "name": "Electronics",
  "icon_url": "http://localhost:9000/greenlight/categories/electronics-icon.png"
}
```

hoặc

```bash
PUT /api/v1/categories/{id}
{
  "icon_url": "http://localhost:9000/greenlight/categories/electronics-icon.png"
}
```

---

## Summary

| Thay đổi | Type | Breaking | Migration Required |
|----------|------|----------|-------------------|
| Xóa `products.gallery` | Database + API | ✅ Yes | Yes |
| Thêm `categories.icon_url` | Database + API | ❌ No | Yes |

## Files Changed

### Database
- `migrations/001_create_tables.sql` - Updated schema
- `migrations/002_remove_gallery_add_icon.sql` - Migration script

### Domain Models
- `internal/domain/product.go` - Removed Gallery type and field
- `internal/domain/category.go` - Added IconURL field

### DTOs
- `internal/transport/http/dto/product_dto.go` - Removed gallery from all DTOs
- `internal/transport/http/dto/category_dto.go` - Added icon_url to all DTOs

### Use Cases
- `internal/usecase/product_usecase.go` - Removed gallery handling
- `internal/usecase/category_usecase.go` - Added icon_url handling

### Handlers
- `internal/transport/http/handler/product_handler.go` - Removed gallery
- `internal/transport/http/handler/category_handler.go` - Added icon_url

### Documentation
- `API_EXAMPLES.md` - Updated examples
- `docs/SCHEMA_CHANGES_2025_01.md` - This document

## Rollback Plan

Nếu cần rollback:

```sql
-- Rollback: Add gallery back to products
ALTER TABLE products 
ADD COLUMN gallery JSON 
AFTER thumbnail_url;

-- Rollback: Remove icon_url from categories
ALTER TABLE categories 
DROP COLUMN icon_url;
```

**Note**: Rollback sẽ mất dữ liệu icon_url đã lưu. Backup trước khi rollback!

