# API Examples - Product Variants

## Tạo Product với Variants

### 1. Tạo Product đơn giản (không có variants)

```bash
POST /api/v1/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Wireless Mouse",
  "slug": "wireless-mouse",
  "sku": "MOUSE-001",
  "short_desc": "Ergonomic wireless mouse",
  "description": "High-precision wireless mouse with ergonomic design",
  "price": 29.99,
  "stock": 100,
  "thumbnail_url": "/uploads/mouse.jpg",
  "gallery": ["/uploads/mouse-1.jpg", "/uploads/mouse-2.jpg"],
  "category_id": "category-uuid-here",
  "is_active": true
}
```

### 2. Tạo Product với Variants (áo thun nhiều màu và size)

```bash
POST /api/v1/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Cotton T-Shirt",
  "slug": "cotton-t-shirt",
  "sku": "",
  "short_desc": "Premium cotton t-shirt",
  "description": "Soft, breathable cotton t-shirt available in multiple colors and sizes",
  "price": 29.99,
  "stock": 0,
  "thumbnail_url": "/uploads/tshirt.jpg",
  "gallery": ["/uploads/tshirt-1.jpg", "/uploads/tshirt-2.jpg"],
  "category_id": "category-uuid-here",
  "is_active": true,
  "variants": [
    {
      "sku": "TSHIRT-RED-S",
      "name": "Red - Small",
      "attributes": {
        "color": "red",
        "size": "S"
      },
      "price": null,
      "stock": 10,
      "is_active": true
    },
    {
      "sku": "TSHIRT-RED-M",
      "name": "Red - Medium",
      "attributes": {
        "color": "red",
        "size": "M"
      },
      "price": null,
      "stock": 15,
      "is_active": true
    },
    {
      "sku": "TSHIRT-RED-L",
      "name": "Red - Large",
      "attributes": {
        "color": "red",
        "size": "L"
      },
      "price": null,
      "stock": 20,
      "is_active": true
    },
    {
      "sku": "TSHIRT-BLUE-S",
      "name": "Blue - Small",
      "attributes": {
        "color": "blue",
        "size": "S"
      },
      "price": null,
      "stock": 8,
      "is_active": true
    },
    {
      "sku": "TSHIRT-BLUE-M",
      "name": "Blue - Medium",
      "attributes": {
        "color": "blue",
        "size": "M"
      },
      "price": null,
      "stock": 12,
      "is_active": true
    },
    {
      "sku": "TSHIRT-BLUE-L",
      "name": "Blue - Large",
      "attributes": {
        "color": "blue",
        "size": "L"
      },
      "price": 31.99,
      "stock": 18,
      "is_active": true
    }
  ]
}
```

**Response:**

```json
{
  "success": true,
  "message": "Product created",
  "data": {
    "product_id": "01JFXXX...",
    "name": "Cotton T-Shirt",
    "slug": "cotton-t-shirt",
    "sku": "",
    "short_desc": "Premium cotton t-shirt",
    "description": "Soft, breathable cotton t-shirt available in multiple colors and sizes",
    "price": 29.99,
    "stock": 0,
    "thumbnail_url": "/uploads/tshirt.jpg",
    "gallery": ["/uploads/tshirt-1.jpg", "/uploads/tshirt-2.jpg"],
    "category_id": "category-uuid-here",
    "is_active": true,
    "variants": [
      {
        "variant_id": "01JFYYY...",
        "product_id": "01JFXXX...",
        "sku": "TSHIRT-RED-S",
        "name": "Red - Small",
        "attributes": {
          "color": "red",
          "size": "S"
        },
        "price": null,
        "stock": 10,
        "is_active": true,
        "created_at": "2025-12-19T13:00:00Z",
        "updated_at": "2025-12-19T13:00:00Z"
      },
      ...
    ],
    "created_at": "2025-12-19T13:00:00Z",
    "updated_at": "2025-12-19T13:00:00Z"
  }
}
```

## Quản lý Variants

### 3. Lấy danh sách variants của product

```bash
GET /api/v1/products/{product_id}/variants
```

**Response:**

```json
{
  "success": true,
  "message": "Variants retrieved",
  "data": [
    {
      "variant_id": "01JFYYY...",
      "product_id": "01JFXXX...",
      "sku": "TSHIRT-RED-S",
      "name": "Red - Small",
      "attributes": {
        "color": "red",
        "size": "S"
      },
      "price": null,
      "stock": 10,
      "is_active": true,
      "created_at": "2025-12-19T13:00:00Z",
      "updated_at": "2025-12-19T13:00:00Z"
    },
    ...
  ]
}
```

### 4. Tạo variant mới cho product đã tồn tại

```bash
POST /api/v1/products/{product_id}/variants
Authorization: Bearer <token>
Content-Type: application/json

{
  "sku": "TSHIRT-GREEN-M",
  "name": "Green - Medium",
  "attributes": {
    "color": "green",
    "size": "M"
  },
  "price": 32.99,
  "stock": 15,
  "is_active": true
}
```

### 5. Cập nhật variant

```bash
PUT /api/v1/products/{product_id}/variants/{variant_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "stock": 25,
  "price": 34.99,
  "is_active": true
}
```

### 6. Xóa variant

```bash
DELETE /api/v1/products/{product_id}/variants/{variant_id}
Authorization: Bearer <token>
```

### 7. Tìm variant theo SKU

```bash
GET /api/v1/variants/sku/TSHIRT-RED-L
```

**Response:**

```json
{
  "success": true,
  "message": "Variant retrieved",
  "data": {
    "variant_id": "01JFYYY...",
    "product_id": "01JFXXX...",
    "sku": "TSHIRT-RED-L",
    "name": "Red - Large",
    "attributes": {
      "color": "red",
      "size": "L"
    },
    "price": null,
    "stock": 20,
    "is_active": true,
    "created_at": "2025-12-19T13:00:00Z",
    "updated_at": "2025-12-19T13:00:00Z"
  }
}
```

### 8. Lấy product kèm variants

```bash
GET /api/v1/products/{product_id}
```

**Response:**

```json
{
  "success": true,
  "message": "Product retrieved",
  "data": {
    "product_id": "01JFXXX...",
    "name": "Cotton T-Shirt",
    "slug": "cotton-t-shirt",
    "sku": "",
    "short_desc": "Premium cotton t-shirt",
    "description": "Soft, breathable cotton t-shirt available in multiple colors and sizes",
    "price": 29.99,
    "stock": 0,
    "thumbnail_url": "/uploads/tshirt.jpg",
    "gallery": ["/uploads/tshirt-1.jpg", "/uploads/tshirt-2.jpg"],
    "category_id": "category-uuid-here",
    "category": {
      "category_id": "category-uuid-here",
      "name": "Clothing",
      "slug": "clothing",
      "description": "Clothing and apparel",
      "is_active": true,
      "created_at": "2025-12-19T12:00:00Z",
      "updated_at": "2025-12-19T12:00:00Z"
    },
    "is_active": true,
    "variants": [
      {
        "variant_id": "01JFYYY...",
        "product_id": "01JFXXX...",
        "sku": "TSHIRT-RED-S",
        "name": "Red - Small",
        "attributes": {
          "color": "red",
          "size": "S"
        },
        "price": null,
        "stock": 10,
        "is_active": true,
        "created_at": "2025-12-19T13:00:00Z",
        "updated_at": "2025-12-19T13:00:00Z"
      },
      ...
    ],
    "created_at": "2025-12-19T13:00:00Z",
    "updated_at": "2025-12-19T13:00:00Z"
  }
}
```

## Ví dụ khác

### Điện thoại với nhiều dung lượng

```json
{
  "name": "Smartphone Pro",
  "slug": "smartphone-pro",
  "sku": "",
  "short_desc": "Latest flagship smartphone",
  "description": "Premium smartphone with advanced features",
  "price": 999.99,
  "stock": 0,
  "thumbnail_url": "/uploads/phone.jpg",
  "gallery": ["/uploads/phone-1.jpg"],
  "category_id": "electronics-uuid",
  "is_active": true,
  "variants": [
    {
      "sku": "PHONE-128-BLK",
      "name": "128GB - Black",
      "attributes": {
        "storage": "128GB",
        "color": "black"
      },
      "price": null,
      "stock": 50,
      "is_active": true
    },
    {
      "sku": "PHONE-256-BLK",
      "name": "256GB - Black",
      "attributes": {
        "storage": "256GB",
        "color": "black"
      },
      "price": 1099.99,
      "stock": 30,
      "is_active": true
    },
    {
      "sku": "PHONE-512-BLK",
      "name": "512GB - Black",
      "attributes": {
        "storage": "512GB",
        "color": "black"
      },
      "price": 1299.99,
      "stock": 20,
      "is_active": true
    }
  ]
}
```

## Lưu ý

1. **SKU ở product level** là optional nếu dùng variants
2. **SKU ở variant level** là required và unique
3. **Price** của variant có thể null (dùng giá của product) hoặc override
4. **Stock** tổng của product = sum(stock của tất cả variants)
5. **Attributes** là JSON object flexible (color, size, storage, material, etc.)
6. Khi GET product, variants sẽ tự động được include trong response
