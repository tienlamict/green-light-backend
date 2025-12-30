# Product API - Variants Integration Update

## Tổng quan
Đã cập nhật tất cả các API endpoint của Product để trả về thông tin đầy đủ về các variant của từng sản phẩm.

## Các API đã được cập nhật

### 1. **GET /api/v1/products** - List Products
**Mô tả**: Lấy danh sách sản phẩm với phân trang

**Thay đổi**:
- ✅ Mỗi product trong danh sách giờ đây bao gồm thông tin về tất cả variants của nó
- ✅ Tự động load variants cho từng product trong list

**Response Example**:
```json
{
  "success": true,
  "message": "Success",
  "data": [
    {
      "product_id": "01JGXXX...",
      "name": "iPhone 15 Pro",
      "slug": "iphone-15-pro",
      "sku": "IP15P",
      "price_min": 999.99,
      "price_max": 1499.99,
      "stock": 100,
      "variants": [
        {
          "variant_id": "01JGYYY...",
          "sku": "IP15P-128-BLK",
          "name": "128GB - Black",
          "attributes": {
            "storage": "128GB",
            "color": "Black"
          },
          "price": 999.99,
          "stock": 50,
          "is_active": true
        },
        {
          "variant_id": "01JGZZZ...",
          "sku": "IP15P-256-BLK",
          "name": "256GB - Black",
          "attributes": {
            "storage": "256GB",
            "color": "Black"
          },
          "price": 1199.99,
          "stock": 50,
          "is_active": true
        }
      ],
      "category": {
        "category_id": "...",
        "name": "Smartphones"
      },
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": 100,
    "page": 1,
    "limit": 10,
    "total_pages": 10
  }
}
```

### 2. **GET /api/v1/products/{id_or_slug}** - Get Single Product
**Mô tả**: Lấy chi tiết một sản phẩm theo ID hoặc slug

**Thay đổi**:
- ✅ Đã có sẵn từ trước, giữ nguyên logic
- ✅ Trả về đầy đủ thông tin variants

**Response Example**:
```json
{
  "success": true,
  "message": "Product retrieved",
  "data": {
    "product_id": "01JGXXX...",
    "name": "iPhone 15 Pro",
    "slug": "iphone-15-pro",
    "sku": "IP15P",
    "short_desc": "Latest iPhone with A17 Pro chip",
    "description": "Full description...",
    "price_min": 999.99,
    "price_max": 1499.99,
    "stock": 100,
    "thumbnail_url": "https://...",
    "gallery": ["https://...", "https://..."],
    "variants": [
      {
        "variant_id": "01JGYYY...",
        "sku": "IP15P-128-BLK",
        "name": "128GB - Black",
        "attributes": {
          "storage": "128GB",
          "color": "Black"
        },
        "price": 999.99,
        "stock": 50,
        "is_active": true,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "category": {
      "category_id": "...",
      "name": "Smartphones",
      "slug": "smartphones"
    },
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 3. **POST /api/v1/products** - Create Product
**Mô tả**: Tạo sản phẩm mới (yêu cầu ít nhất 1 variant)

**Thay đổi**:
- ✅ Sau khi tạo product, tự động load và trả về danh sách variants đã tạo
- ✅ Response bao gồm đầy đủ thông tin variants

**Request Example**:
```json
{
  "name": "iPhone 15 Pro",
  "slug": "iphone-15-pro",
  "sku": "IP15P",
  "short_desc": "Latest iPhone",
  "description": "Full description...",
  "stock": 0,
  "thumbnail_url": "https://...",
  "gallery": ["https://..."],
  "category_id": "01JG...",
  "is_active": true,
  "variants": [
    {
      "sku": "IP15P-128-BLK",
      "name": "128GB - Black",
      "attributes": {
        "storage": "128GB",
        "color": "Black"
      },
      "price": 999.99,
      "stock": 50,
      "is_active": true
    },
    {
      "sku": "IP15P-256-BLK",
      "name": "256GB - Black",
      "attributes": {
        "storage": "256GB",
        "color": "Black"
      },
      "price": 1199.99,
      "stock": 50,
      "is_active": true
    }
  ]
}
```

**Response**: Tương tự GET single product, bao gồm đầy đủ variants

### 4. **PUT /api/v1/products/{id_or_slug}** - Update Product and Variants
**Mô tả**: Cập nhật thông tin sản phẩm và variants (nếu có)

**Thay đổi**:
- ✅ Sau khi update, tự động load và trả về danh sách variants hiện tại
- ✅ Response bao gồm đầy đủ thông tin variants
- ✅ **MỚI**: Có thể update variants cùng lúc với product trong cùng 1 request

**Request Example - Chỉ update product**:
```json
{
  "name": "iPhone 15 Pro Updated",
  "short_desc": "Updated description"
}
```

**Request Example - Update cả product và variants**:
```json
{
  "name": "iPhone 15 Pro Updated",
  "short_desc": "Updated description",
  "variants": [
    {
      "variant_id": "01JGYYY...",
      "price": 899.99,
      "stock": 45
    },
    {
      "variant_id": "01JGZZZ...",
      "name": "256GB - Black (Updated)",
      "price": 1099.99,
      "stock": 55,
      "is_active": true
    }
  ]
}
```

**Lưu ý**:
- Khi update variants qua API này, **bắt buộc** phải có `variant_id` cho mỗi variant
- Chỉ update các field được cung cấp (partial update)
- Nếu không muốn update variants, có thể bỏ qua field `variants` hoặc để array rỗng
- Vẫn có thể sử dụng API variant riêng nếu muốn:
  - `PUT /api/v1/products/{product_id}/variants/{variant_id}` (không cần variant_id trong body)
  - `POST /api/v1/products/{product_id}/variants` (tạo variant mới)
  - `DELETE /api/v1/products/{product_id}/variants/{variant_id}` (xóa variant)

### 5. **DELETE /api/v1/products/{id_or_slug}** - Delete Product
**Mô tả**: Xóa sản phẩm

**Thay đổi**: Không có thay đổi (không cần trả về variants khi xóa)

## Chi tiết kỹ thuật

### Files đã thay đổi:
1. **internal/transport/http/handler/product_handler.go**
   - Updated `List()` method: Thêm logic load variants cho mỗi product
   - Updated `Update()` method: Thêm logic load variants sau khi update
   - `Get()` và `Create()` đã có sẵn logic load variants

### Cấu trúc Response:
```go
type ProductResponse struct {
    ProductID    string            `json:"product_id"`
    Name         string            `json:"name"`
    Slug         string            `json:"slug"`
    SKU          string            `json:"sku"`
    ShortDesc    string            `json:"short_desc"`
    Description  string            `json:"description"`
    PriceMin     *float64          `json:"price_min"`    // Min price từ variants
    PriceMax     *float64          `json:"price_max"`    // Max price từ variants
    Stock        int               `json:"stock"`
    ThumbnailURL string            `json:"thumbnail_url"`
    Gallery      []string          `json:"gallery"`
    CategoryID   string            `json:"category_id"`
    Category     *CategoryResponse `json:"category,omitempty"`
    IsActive     bool              `json:"is_active"`
    Variants     []VariantResponse `json:"variants,omitempty"` // ✅ Luôn có variants
    CreatedAt    string            `json:"created_at"`
    UpdatedAt    string            `json:"updated_at"`
}

type VariantResponse struct {
    VariantID  string            `json:"variant_id"`
    ProductID  string            `json:"product_id"`
    SKU        string            `json:"sku"`
    Name       string            `json:"name"`
    Attributes map[string]string `json:"attributes"`
    Price      float64           `json:"price"`
    Stock      int               `json:"stock"`
    IsActive   bool              `json:"is_active"`
    CreatedAt  string            `json:"created_at"`
    UpdatedAt  string            `json:"updated_at"`
}
```

## Performance Considerations

### List API với nhiều products:
- Mỗi product sẽ có 1 query riêng để load variants
- Với 10 products/page → 10 queries thêm cho variants
- **Khuyến nghị**: 
  - Giữ limit nhỏ (10-20 items/page)
  - Có thể optimize sau bằng cách preload variants trong repository layer
  - Hoặc thêm query parameter `include_variants=true/false` để client có thể chọn

### Optimization ideas (future):
```go
// Option 1: Batch load variants
variantsByProductID := h.variantUseCase.GetByProductIDs(ctx, productIDs)

// Option 2: Query parameter
includeVariants := c.DefaultQuery("include_variants", "true") == "true"
if includeVariants {
    // Load variants
}

// Option 3: Repository level preload
products, err := h.productRepo.ListWithVariants(ctx, filter)
```

## Testing

### Test các API:

1. **List Products**:
```bash
curl http://localhost:8080/api/v1/products?page=1&limit=10
```

2. **Get Single Product**:
```bash
curl http://localhost:8080/api/v1/products/01JGXXX...
# hoặc
curl http://localhost:8080/api/v1/products/iphone-15-pro
```

3. **Create Product**:
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Test Product",
    "category_id": "01JG...",
    "variants": [
      {
        "sku": "TEST-001",
        "name": "Default",
        "price": 99.99,
        "stock": 10
      }
    ]
  }'
```

4. **Update Product (chỉ product)**:
```bash
curl -X PUT http://localhost:8080/api/v1/products/01JGXXX... \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Updated Name"
  }'
```

5. **Update Product và Variants (merged API)**:
```bash
curl -X PUT http://localhost:8080/api/v1/products/01JGXXX... \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Updated Name",
    "variants": [
      {
        "variant_id": "01JGYYY...",
        "price": 899.99,
        "stock": 45
      }
    ]
  }'
```

## Migration Notes

### Breaking Changes: ❌ NONE
- Tất cả các thay đổi đều backward compatible
- Response chỉ thêm field `variants`, không xóa hay thay đổi field nào
- Client cũ vẫn hoạt động bình thường, chỉ cần ignore field `variants`

### New Features: ✅
- ✅ Tất cả product APIs giờ đều trả về variants
- ✅ Không cần gọi thêm API riêng để lấy variants
- ✅ Frontend có thể hiển thị đầy đủ thông tin ngay từ list view
- ✅ **MỚI**: API Update Product giờ có thể update cả variants trong cùng 1 request

## Summary

| API Endpoint | Method | Variants Included | Can Update Variants | Status |
|-------------|--------|-------------------|---------------------|--------|
| `/api/v1/products` | GET | ✅ Yes | N/A | Updated |
| `/api/v1/products/{id_or_slug}` | GET | ✅ Yes | N/A | Already had |
| `/api/v1/products` | POST | ✅ Yes | ✅ Yes (create) | Already had |
| `/api/v1/products/{id_or_slug}` | PUT | ✅ Yes | ✅ Yes (update) | **MERGED** |
| `/api/v1/products/{id_or_slug}` | DELETE | N/A | N/A | No change |
| `/api/v1/products/{id}/variants/{variant_id}` | PUT | N/A | ✅ Yes | Still available |

**Tất cả các API đọc (GET, POST, PUT) giờ đều trả về đầy đủ thông tin variants! 🎉**

**API Update Product giờ có thể update cả product và variants trong cùng 1 request! 🚀**

