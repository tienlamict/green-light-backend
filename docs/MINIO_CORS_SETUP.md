# 🔧 MinIO CORS Configuration

## Tổng quan

MinIO CORS được tự động cấu hình khi khởi động Docker Compose thông qua service `minio-setup`.

## Cách hoạt động

1. **Service `minio-setup`** sử dụng MinIO Client (`mc`) để:
   - Tạo bucket nếu chưa tồn tại
   - Cấu hình CORS từ file `minio/cors.json`

2. **Service này chỉ chạy một lần** (`restart: "no"`) sau khi MinIO đã sẵn sàng

## Cấu hình CORS

File cấu hình: `minio/cors.json`

```json
{
  "CORSRules": [
    {
      "AllowedOrigins": ["*"],
      "AllowedMethods": ["GET", "PUT", "POST", "DELETE", "HEAD"],
      "AllowedHeaders": ["*"],
      "ExposeHeaders": ["ETag", "Content-Length", "x-amz-request-id"],
      "MaxAgeSeconds": 3000
    }
  ]
}
```

### Các tham số:

- **AllowedOrigins**: Cho phép tất cả origins (`*`) hoặc chỉ định domain cụ thể
  - Ví dụ: `["http://localhost:3000", "https://example.com"]`
  
- **AllowedMethods**: Các HTTP methods được phép
  - GET: Lấy file
  - PUT: Upload file (presigned URL)
  - POST: Upload multipart
  - DELETE: Xóa file
  - HEAD: Lấy metadata

- **AllowedHeaders**: Headers được phép trong request
  - `["*"]` cho phép tất cả headers

- **ExposeHeaders**: Headers mà client có thể đọc từ response

- **MaxAgeSeconds**: Thời gian cache preflight request (giây)

## Chạy setup

### Lần đầu tiên

```bash
docker-compose up -d
```

Service `minio-setup` sẽ tự động chạy và cấu hình CORS.

### Chạy lại setup thủ công

Nếu cần cấu hình lại CORS:

```bash
# Xóa container setup cũ (nếu có)
docker-compose rm -f minio-setup

# Chạy lại setup
docker-compose up minio-setup
```

### Kiểm tra CORS đã được cấu hình

```bash
# Vào container minio-setup
docker-compose exec minio-setup mc cors get myminio/greenlight
```

Hoặc sử dụng MinIO Console:
1. Truy cập http://localhost:9001
2. Đăng nhập với `minioadmin`/`minioadmin`
3. Vào Buckets → greenlight → Management → CORS Configuration

## Troubleshooting

### Lỗi: "MinIO is not ready"

- Đợi thêm một chút, MinIO có thể cần thời gian khởi động
- Kiểm tra healthcheck của MinIO: `docker-compose ps minio`

### Lỗi: "Bucket already exists"

- Đây không phải lỗi, chỉ là thông báo bucket đã tồn tại

### CORS không hoạt động

1. Kiểm tra file `minio/cors.json` có đúng format không
2. Chạy lại setup: `docker-compose up minio-setup`
3. Kiểm tra logs: `docker-compose logs minio-setup`
4. Xác nhận CORS trong MinIO Console

## Cấu hình cho Production

Trong production, nên thay đổi `AllowedOrigins` từ `["*"]` thành danh sách domain cụ thể:

```json
{
  "CORSRules": [
    {
      "AllowedOrigins": [
        "https://yourdomain.com",
        "https://www.yourdomain.com"
      ],
      "AllowedMethods": ["GET", "PUT", "POST", "DELETE", "HEAD"],
      "AllowedHeaders": ["*"],
      "ExposeHeaders": ["ETag", "Content-Length", "x-amz-request-id"],
      "MaxAgeSeconds": 3000
    }
  ]
}
```

## Liên quan

- [MINIO_README.md](MINIO_README.md) - Tài liệu tổng quan về MinIO
- [MINIO_QUICK_START.md](MINIO_QUICK_START.md) - Hướng dẫn nhanh

