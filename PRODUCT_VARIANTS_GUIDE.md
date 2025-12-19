# Product Variants Guide - Nhiều SKU (Stock Keeping Unit) cho 1 Sản Phẩm

## Vấn đề

Trong thực tế, một sản phẩm có thể có nhiều SKU khác nhau:
- **Ví dụ:** Áo thun có nhiều màu và size
  - SKU: `TSHIRT-RED-L` (Red, Large)
  - SKU: `TSHIRT-RED-M` (Red, Medium)
  - SKU: `TSHIRT-BLUE-L` (Blue, Large)
  - SKU: `TSHIRT-BLUE-M` (Blue, Medium)

## Giải pháp: Product Variants

Đã thêm bảng `product_variants` để quản lý các biến thể của sản phẩm.

### Schema mới

```
products (bảng chính)
├── product_id (PK)
├── name: "T-Shirt Basic"
├── slug: "t-shirt-basic"
├── sku: NULL (optional - có thể để trống nếu dùng variants)
└── price: 29.99 (giá cơ bản)

product_variants (biến thể)
├── variant_id (PK)
├── product_id (FK → products)
├── sku: "TSHIRT-RED-L" (unique)
├── name: "Red - Large"
├── attributes: {"color": "red", "size": "L"}
├── price: NULL (dùng giá của product) hoặc 31.99 (override)
└── stock: 50
```

## Cách sử dụng

### Option 1: Product không có variants (SKU đơn giản)

```json
{
  "name": "Wireless Mouse",
  "slug": "wireless-mouse",
  "sku": "MOUSE-001",  // SKU trực tiếp trong product
  "price": 29.99,
  "stock": 100
}
```

### Option 2: Product có nhiều variants

```json
{
  "name": "T-Shirt Basic",
  "slug": "t-shirt-basic",
  "sku": null,  // Không có SKU ở product level
  "price": 29.99,  // Giá cơ bản
  "stock": 0,  // Stock tổng = sum của variants
  "variants": [
    {
      "sku": "TSHIRT-RED-L",
      "name": "Red - Large",
      "attributes": {"color": "red", "size": "L"},
      "price": null,  // Dùng giá của product
      "stock": 20
    },
    {
      "sku": "TSHIRT-RED-M",
      "name": "Red - Medium",
      "attributes": {"color": "red", "size": "M"},
      "price": null,
      "stock": 15
    },
    {
      "sku": "TSHIRT-BLUE-L",
      "name": "Blue - Large",
      "attributes": {"color": "blue", "size": "L"},
      "price": 31.99,  // Override giá (đắt hơn)
      "stock": 10
    }
  ]
}
```

## Migration

Chạy migration để thêm bảng variants:

```bash
docker-compose exec db mysql --skip-ssl -u root -proot greenlight_db < migrations/002_add_product_variants.sql
```

Hoặc reset hoàn toàn:

```bash
docker-compose down -v
docker-compose up -d
# Migration sẽ chạy tự động
```

## API Endpoints (Cần implement)

### Tạo variant

```bash
POST /api/v1/products/{product_id}/variants
{
  "sku": "TSHIRT-RED-L",
  "name": "Red - Large",
  "attributes": {
    "color": "red",
    "size": "L"
  },
  "price": null,
  "stock": 20
}
```

### Lấy variants của product

```bash
GET /api/v1/products/{product_id}/variants
```

### Lấy product theo variant SKU

```bash
GET /api/v1/products/by-sku/{sku}
```

## Ví dụ thực tế

### Áo thun có nhiều màu và size

**Product:**
- Name: "Cotton T-Shirt"
- Slug: "cotton-t-shirt"
- Base Price: $29.99

**Variants:**
1. Red - Small: SKU `TSHIRT-RED-S`, Stock: 10
2. Red - Medium: SKU `TSHIRT-RED-M`, Stock: 15
3. Red - Large: SKU `TSHIRT-RED-L`, Stock: 20
4. Blue - Small: SKU `TSHIRT-BLUE-S`, Stock: 8
5. Blue - Medium: SKU `TSHIRT-BLUE-M`, Stock: 12
6. Blue - Large: SKU `TSHIRT-BLUE-L`, Stock: 18

### Điện thoại có nhiều dung lượng

**Product:**
- Name: "Smartphone Pro"
- Slug: "smartphone-pro"
- Base Price: $999.99

**Variants:**
1. 128GB - Black: SKU `PHONE-128-BLK`, Stock: 50
2. 256GB - Black: SKU `PHONE-256-BLK`, Stock: 30, Price: $1099.99
3. 512GB - Black: SKU `PHONE-512-BLK`, Stock: 20, Price: $1299.99
4. 128GB - White: SKU `PHONE-128-WHT`, Stock: 40
5. 256GB - White: SKU `PHONE-256-WHT`, Stock: 25, Price: $1099.99

## Lợi ích

1. ✅ **Quản lý kho chính xác** - Mỗi variant có stock riêng
2. ✅ **Giá linh hoạt** - Có thể override giá cho từng variant
3. ✅ **SEO tốt** - Product chính vẫn có slug duy nhất
4. ✅ **Dễ mở rộng** - Thêm variant mới không cần tạo product mới
5. ✅ **Tìm kiếm tốt** - Tìm theo SKU variant hoặc attributes

## Lưu ý

- **SKU ở product level** là optional nếu dùng variants
- **SKU ở variant level** là required và unique
- **Stock tổng** của product = sum(stock của tất cả variants)
- **Price** của variant có thể override price của product
- **Attributes** lưu dạng JSON để flexible (color, size, material, etc.)

## Next Steps

1. ✅ Migration đã tạo
2. ⏳ Cần implement Repository cho ProductVariant
3. ⏳ Cần implement UseCase cho ProductVariant
4. ⏳ Cần implement Handler cho ProductVariant API
5. ⏳ Cần update Product response để include variants

Bạn có muốn tôi implement đầy đủ Product Variants API không?

