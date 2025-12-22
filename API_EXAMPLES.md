# API Examples - Product Variants & MinIO Image Upload

## 🔐 Authentication

**Lưu ý về Authentication:**
- ✅ **GET endpoints là PUBLIC** - Không cần authentication
- 🔒 **POST/PUT/DELETE endpoints yêu cầu authentication** - Cần Bearer token và role (admin/editor)

Để lấy token cho các protected endpoints, gọi API login:
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@example.com",
  "password": "your_password"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "user_id": "a628f74e-6a75-4bc8-a287-51264a87b6fd",
      "email": "admin@example.com",
      "role": "admin"
    }
  }
}
```

Sau đó sử dụng token này trong header `Authorization: Bearer <token>` cho các POST/PUT/DELETE requests.

---

## 📦 Tạo Product với Variants (Pricing Model: Min-Max)

### ⚠️ Lưu ý quan trọng về Pricing Model

- **Product không có field `price`** - Giá được tính từ variants
- **Product có `price_min` và `price_max`** - Tự động tính từ giá của variants
- **Mỗi variant PHẢI có `price`** - Không nullable, giá trị bắt buộc
- **Ít nhất 1 variant là bắt buộc** khi tạo product

### 1. Tạo Product với Variants (Áo thun nhiều màu và size)

```bash
POST /api/v1/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Cotton T-Shirt",
  "slug": "cotton-t-shirt",
  "sku": null,
  "short_desc": "Premium cotton t-shirt",
  "description": "Soft, breathable cotton t-shirt available in multiple colors and sizes",
  "stock": 0,
  "thumbnail_url": "http://localhost:9000/greenlight/products/2024/12/prod-123/main.webp",
  "gallery": [],
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
      "price": 29.99,
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
      "price": 29.99,
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
      "price": 31.99,
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
      "price": 29.99,
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
      "price": 29.99,
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
  "message": "Product created successfully",
  "data": {
    "product_id": "018f1234-5678-7890-abcd-ef1234567890",
    "name": "Cotton T-Shirt",
    "slug": "cotton-t-shirt",
    "sku": null,
    "short_desc": "Premium cotton t-shirt",
    "description": "Soft, breathable cotton t-shirt available in multiple colors and sizes",
    "price_min": 29.99,
    "price_max": 31.99,
    "stock": 83,
    "thumbnail_url": "http://localhost:9000/greenlight/products/2024/12/prod-123/main.webp",
    "gallery": [],
    "category_id": "category-uuid-here",
    "is_active": true,
    "variants": [
      {
        "variant_id": "018f1235-5678-7890-abcd-ef1234567890",
        "product_id": "018f1234-5678-7890-abcd-ef1234567890",
        "sku": "TSHIRT-RED-S",
        "name": "Red - Small",
        "attributes": {
          "color": "red",
          "size": "S"
        },
        "price": 29.99,
        "stock": 10,
        "is_active": true,
        "created_at": "2024-12-21T10:00:00Z",
        "updated_at": "2024-12-21T10:00:00Z"
      },
      {
        "variant_id": "018f1236-5678-7890-abcd-ef1234567890",
        "product_id": "018f1234-5678-7890-abcd-ef1234567890",
        "sku": "TSHIRT-RED-M",
        "name": "Red - Medium",
        "attributes": {
          "color": "red",
          "size": "M"
        },
        "price": 29.99,
        "stock": 15,
        "is_active": true,
        "created_at": "2024-12-21T10:00:00Z",
        "updated_at": "2024-12-21T10:00:00Z"
      },
      {
        "variant_id": "018f1237-5678-7890-abcd-ef1234567890",
        "product_id": "018f1234-5678-7890-abcd-ef1234567890",
        "sku": "TSHIRT-RED-L",
        "name": "Red - Large",
        "attributes": {
          "color": "red",
          "size": "L"
        },
        "price": 31.99,
        "stock": 20,
        "is_active": true,
        "created_at": "2024-12-21T10:00:00Z",
        "updated_at": "2024-12-21T10:00:00Z"
      }
      // ... more variants
    ],
    "created_at": "2024-12-21T10:00:00Z",
    "updated_at": "2024-12-21T10:00:00Z"
  }
}
```

**Giải thích:**
- `price_min: 29.99` - Giá thấp nhất trong các variants
- `price_max: 31.99` - Giá cao nhất trong các variants
- `stock: 83` - Tổng stock của tất cả variants (10+15+20+8+12+18)
- Variants có `price` bắt buộc (không nullable)

## 🖼️ Upload Ảnh Sản Phẩm với MinIO

### 2. Upload ảnh sản phẩm (3 bước)

#### Bước 1: Yêu cầu Presigned URL

```bash
POST /api/v1/products/{product_id}/images/presign
Authorization: Bearer <token>
Content-Type: application/json

{
  "content_type": "image/webp",
  "extension": "webp"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Presigned URL generated successfully",
  "data": {
    "upload_url": "http://localhost:9000/greenlight/products/2024/12/prod-123/uuid.webp?X-Amz-Algorithm=...",
    "public_url": "http://localhost:9000/greenlight/products/2024/12/prod-123/uuid.webp",
    "object_key": "products/2024/12/prod-123/uuid.webp",
    "expires_at": "2024-12-21T10:10:00Z"
  }
}
```

#### Bước 2: Upload ảnh trực tiếp lên MinIO

```bash
PUT {upload_url}
Content-Type: image/webp

[Binary image data]
```

**JavaScript Example:**

```javascript
const fileInput = document.getElementById('imageInput');
const file = fileInput.files[0];

// Upload directly to MinIO
const uploadResponse = await fetch(presignData.upload_url, {
  method: 'PUT',
  headers: {
    'Content-Type': 'image/webp'
  },
  body: file
});

if (uploadResponse.ok) {
  console.log('Image uploaded successfully!');
}
```

#### Bước 3: Confirm upload với Backend

```bash
POST /api/v1/products/{product_id}/images
Authorization: Bearer <token>
Content-Type: application/json

{
  "object_key": "products/2024/12/prod-123/uuid.webp",
  "is_main": true,
  "sort_order": 0
}
```

**Response:**

```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "data": {
    "image_id": "018f1238-5678-7890-abcd-ef1234567890",
    "product_id": "018f1234-5678-7890-abcd-ef1234567890",
    "url": "http://localhost:9000/greenlight/products/2024/12/prod-123/uuid.webp",
    "is_main": true,
    "sort_order": 0,
    "status": "ACTIVE",
    "created_at": "2024-12-21T10:05:00Z"
  }
}
```

#### Complete Flow Example (JavaScript)

```javascript
async function uploadProductImage(productId, imageFile, isMain = false) {
  // Step 1: Get presigned URL
  const presignRes = await fetch(`/api/v1/products/${productId}/images/presign`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      content_type: imageFile.type,
      extension: imageFile.name.split('.').pop()
    })
  });
  
  const { data: presignData } = await presignRes.json();
  
  // Step 2: Upload to MinIO
  await fetch(presignData.upload_url, {
    method: 'PUT',
    headers: { 'Content-Type': imageFile.type },
    body: imageFile
  });
  
  // Step 3: Confirm upload
  const confirmRes = await fetch(`/api/v1/products/${productId}/images`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      object_key: presignData.object_key,
      is_main: isMain,
      sort_order: 0
    })
  });
  
  return await confirmRes.json();
}
```

### 3. List ảnh của sản phẩm

```bash
GET /api/v1/products/{product_id}/images
```

**Response:**

```json
{
  "success": true,
  "message": "Product images retrieved successfully",
  "data": [
    {
      "image_id": "018f1238-5678-7890-abcd-ef1234567890",
      "product_id": "018f1234-5678-7890-abcd-ef1234567890",
      "url": "http://localhost:9000/greenlight/products/2024/12/prod-123/main.webp",
      "is_main": true,
      "sort_order": 0,
      "status": "ACTIVE",
      "created_at": "2024-12-21T10:05:00Z"
    },
    {
      "image_id": "018f1239-5678-7890-abcd-ef1234567890",
      "product_id": "018f1234-5678-7890-abcd-ef1234567890",
      "url": "http://localhost:9000/greenlight/products/2024/12/prod-123/gallery-1.webp",
      "is_main": false,
      "sort_order": 1,
      "status": "ACTIVE",
      "created_at": "2024-12-21T10:06:00Z"
    }
  ]
}
```

### 4. Update ảnh metadata

```bash
PATCH /api/v1/products/{product_id}/images/{image_id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "is_main": true,
  "sort_order": 2
}
```

### 5. Xóa ảnh

```bash
DELETE /api/v1/products/{product_id}/images/{image_id}
Authorization: Bearer <token>
```

**Lưu ý:** Xóa ảnh sẽ xóa cả object trong MinIO và metadata trong database.

---

## 🔧 Quản lý Variants

### 6. Lấy danh sách variants của product

```bash
GET /api/v1/products/{product_id}/variants
```

**Response:**

```json
{
  "success": true,
  "message": "Variants retrieved successfully",
  "data": [
    {
      "variant_id": "018f1235-5678-7890-abcd-ef1234567890",
      "product_id": "018f1234-5678-7890-abcd-ef1234567890",
      "sku": "TSHIRT-RED-S",
      "name": "Red - Small",
      "attributes": {
        "color": "red",
        "size": "S"
      },
      "price": 29.99,
      "stock": 10,
      "is_active": true,
      "created_at": "2024-12-21T10:00:00Z",
      "updated_at": "2024-12-21T10:00:00Z"
    },
    {
      "variant_id": "018f1236-5678-7890-abcd-ef1234567890",
      "product_id": "018f1234-5678-7890-abcd-ef1234567890",
      "sku": "TSHIRT-RED-M",
      "name": "Red - Medium",
      "attributes": {
        "color": "red",
        "size": "M"
      },
      "price": 29.99,
      "stock": 15,
      "is_active": true,
      "created_at": "2024-12-21T10:00:00Z",
      "updated_at": "2024-12-21T10:00:00Z"
    }
  ]
}
```

### 7. Tạo variant mới cho product đã tồn tại

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

**Lưu ý:** Sau khi tạo variant mới, `price_min` và `price_max` của product sẽ tự động cập nhật.

### 8. Cập nhật variant

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

### 9. Xóa variant

```bash
DELETE /api/v1/products/{product_id}/variants/{variant_id}
Authorization: Bearer <token>
```

### 10. Tìm variant theo SKU

```bash
GET /api/v1/variants/sku/TSHIRT-RED-L
```

**Response:**

```json
{
  "success": true,
  "message": "Variant retrieved successfully",
  "data": {
    "variant_id": "018f1237-5678-7890-abcd-ef1234567890",
    "product_id": "018f1234-5678-7890-abcd-ef1234567890",
    "sku": "TSHIRT-RED-L",
    "name": "Red - Large",
    "attributes": {
      "color": "red",
      "size": "L"
    },
    "price": 31.99,
    "stock": 20,
    "is_active": true,
    "created_at": "2024-12-21T10:00:00Z",
    "updated_at": "2024-12-21T10:00:00Z"
  }
}
```

### 11. List Products với Filters

#### 11.1. List tất cả products (có pagination)

```bash
GET /api/v1/products?page=1&limit=10
```

**Response:**

```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": [
    {
      "product_id": "018f1234-5678-7890-abcd-ef1234567890",
      "name": "Cotton T-Shirt",
      "slug": "cotton-t-shirt",
      "price_min": 29.99,
      "price_max": 31.99,
      "stock": 83,
      "category_id": "9d1351ee-7e24-4d43-bacd-fa87b9acbfa8",
      "category": {
        "category_id": "9d1351ee-7e24-4d43-bacd-fa87b9acbfa8",
        "name": "Clothing",
        "slug": "clothing"
      },
      "is_active": true,
      "created_at": "2024-12-21T10:00:00Z"
    }
    // ... more products
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

#### 11.2. List products theo Category ID

```bash
GET /api/v1/products?category=9d1351ee-7e24-4d43-bacd-fa87b9acbfa8&page=1&limit=20
```

**Query Parameters:**
- `category` (string, optional): Filter theo category ID
- `page` (int, optional): Số trang (default: 1)
- `limit` (int, optional): Số items mỗi trang (default: 10)

**Response:** Tương tự như trên, nhưng chỉ trả về products thuộc category đó.

#### 11.3. List products với Search

```bash
GET /api/v1/products?q=led&page=1&limit=10
```

**Query Parameters:**
- `q` (string, optional): Search theo name, description, hoặc SKU

**Response:** Trả về products có chứa từ khóa "led" trong name, description, hoặc SKU.

#### 11.4. List products với Price Range

```bash
GET /api/v1/products?min_price=20&max_price=50&page=1&limit=10
```

**Query Parameters:**
- `min_price` (number, optional): Giá tối thiểu (price_min >= min_price)
- `max_price` (number, optional): Giá tối đa (price_max <= max_price)

**Response:** Trả về products có price_min >= 20 và price_max <= 50.

#### 11.5. List products Active Only

```bash
GET /api/v1/products?is_active=true&page=1&limit=10
```

**Query Parameters:**
- `is_active` (boolean, optional): Filter theo trạng thái active

#### 11.6. List products với Sort

```bash
GET /api/v1/products?sort=price_min ASC&page=1&limit=10
```

**Query Parameters:**
- `sort` (string, optional): Sort order (default: "created_at DESC")
  - `created_at DESC` - Mới nhất trước
  - `created_at ASC` - Cũ nhất trước
  - `price_min ASC` - Giá thấp nhất trước
  - `price_min DESC` - Giá cao nhất trước
  - `name ASC` - Tên A-Z
  - `name DESC` - Tên Z-A

#### 11.7. List products với Multiple Filters

```bash
GET /api/v1/products?category=9d1351ee-7e24-4d43-bacd-fa87b9acbfa8&min_price=25&max_price=100&is_active=true&q=shirt&sort=price_min ASC&page=1&limit=20
```

**Ví dụ đầy đủ với tất cả filters:**

```bash
curl --location 'http://localhost:8080/api/v1/products?category=9d1351ee-7e24-4d43-bacd-fa87b9acbfa8&min_price=25&max_price=100&is_active=true&q=shirt&sort=price_min%20ASC&page=1&limit=20'
```

**Response:**

```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": [
    {
      "product_id": "018f1234-5678-7890-abcd-ef1234567890",
      "name": "Cotton T-Shirt",
      "slug": "cotton-t-shirt",
      "price_min": 29.99,
      "price_max": 31.99,
      "stock": 83,
      "category_id": "9d1351ee-7e24-4d43-bacd-fa87b9acbfa8",
      "category": {
        "category_id": "9d1351ee-7e24-4d43-bacd-fa87b9acbfa8",
        "name": "Clothing",
        "slug": "clothing",
        "description": "Clothing and apparel",
        "is_active": true
      },
      "is_active": true,
      "created_at": "2024-12-21T10:00:00Z",
      "updated_at": "2024-12-21T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

**Lưu ý:**
- ✅ **GET endpoints là PUBLIC** - Không cần authentication
- Tất cả query parameters đều optional
- Có thể kết hợp nhiều filters cùng lúc
- Response không bao gồm variants và images (phải gọi endpoint riêng)
- Category được include trong response nếu có

---

### 12. Lấy product kèm variants và images

```bash
GET /api/v1/products/{product_id}
```

**Response:**

```json
{
  "success": true,
  "message": "Product retrieved successfully",
  "data": {
    "product_id": "018f1234-5678-7890-abcd-ef1234567890",
    "name": "Cotton T-Shirt",
    "slug": "cotton-t-shirt",
    "sku": null,
    "short_desc": "Premium cotton t-shirt",
    "description": "Soft, breathable cotton t-shirt available in multiple colors and sizes",
    "price_min": 29.99,
    "price_max": 31.99,
    "stock": 83,
    "thumbnail_url": "http://localhost:9000/greenlight/products/2024/12/prod-123/main.webp",
    "gallery": [],
    "category_id": "category-uuid-here",
    "category": {
      "category_id": "category-uuid-here",
      "name": "Clothing",
      "slug": "clothing",
      "description": "Clothing and apparel",
      "is_active": true,
      "created_at": "2024-12-21T09:00:00Z",
      "updated_at": "2024-12-21T09:00:00Z"
    },
    "is_active": true,
    "variants": [
      {
        "variant_id": "018f1235-5678-7890-abcd-ef1234567890",
        "product_id": "018f1234-5678-7890-abcd-ef1234567890",
        "sku": "TSHIRT-RED-S",
        "name": "Red - Small",
        "attributes": {
          "color": "red",
          "size": "S"
        },
        "price": 29.99,
        "stock": 10,
        "is_active": true,
        "created_at": "2024-12-21T10:00:00Z",
        "updated_at": "2024-12-21T10:00:00Z"
      }
      // ... more variants
    ],
    "created_at": "2024-12-21T10:00:00Z",
    "updated_at": "2024-12-21T10:00:00Z"
  }
}
```

**Lưu ý:** Để lấy images của product, gọi endpoint riêng:
```bash
GET /api/v1/products/{product_id}/images
```

## 📱 Ví dụ khác

### 13. Điện thoại với nhiều dung lượng

```json
{
  "name": "Smartphone Pro",
  "slug": "smartphone-pro",
  "sku": null,
  "short_desc": "Latest flagship smartphone",
  "description": "Premium smartphone with advanced features",
  "stock": 0,
  "thumbnail_url": "http://localhost:9000/greenlight/products/2024/12/prod-456/main.webp",
  "gallery": [],
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
      "price": 999.99,
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

**Response sẽ có:**
- `price_min: 999.99` (từ variant 128GB)
- `price_max: 1299.99` (từ variant 512GB)

## 📝 Lưu ý quan trọng

### Pricing Model
1. ❌ **Product KHÔNG có field `price`** - Đã bị loại bỏ
2. ✅ **Product có `price_min` và `price_max`** - Tự động tính từ variants
3. ✅ **Mỗi variant PHẢI có `price`** - Bắt buộc, không nullable
4. ✅ **Ít nhất 1 variant là bắt buộc** khi tạo product

### Variants
1. **SKU ở product level** là optional (có thể null) nếu dùng variants
2. **SKU ở variant level** là required và unique
3. **Price của variant** là bắt buộc (không nullable)
4. **Stock tổng** của product = sum(stock của tất cả variants)
5. **Attributes** là JSON object flexible (color, size, storage, material, etc.)
6. Khi GET product, variants sẽ tự động được include trong response

### Images (MinIO)
1. **Upload flow**: 3 bước (presign → upload → confirm)
2. **Presigned URLs** expire sau 10 phút
3. **Allowed types**: `image/jpeg`, `image/png`, `image/webp`
4. **Max size**: 7MB
5. **Object key format**: `products/{yyyy}/{mm}/{product_id}/{uuid}.{ext}`
6. **Public URL**: `http://localhost:9000/greenlight/{object_key}`
7. **Images không tự động include** khi GET product - phải gọi endpoint riêng

### Database Triggers
- Khi tạo/cập nhật/xóa variant → `price_min` và `price_max` tự động cập nhật
- Không cần manually update giá của product

## 🔗 Related Documentation

- [MinIO Integration Guide](docs/MINIO_INTEGRATION.md) - Chi tiết về upload ảnh
- [MinIO Quick Start](docs/MINIO_QUICK_START.md) - Test scripts
- [Product Variants Guide](PRODUCT_VARIANTS_GUIDE.md) - Chi tiết về variants
- [API Examples MinIO](API_EXAMPLES_MINIO.md) - Ví dụ upload ảnh đầy đủ
