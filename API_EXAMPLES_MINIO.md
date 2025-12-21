# MinIO API Examples

Các ví dụ thực tế về cách sử dụng MinIO APIs cho upload ảnh sản phẩm.

---

## 🔐 Authentication

Tất cả protected endpoints yêu cầu JWT token trong header:

```
Authorization: Bearer <your-jwt-token>
```

Lấy token bằng cách login:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

---

## 📸 Upload Image Flow

### Step 1: Request Presigned URL

**Request**:
```bash
curl -X POST http://localhost:8080/api/v1/products/prod-123/images/presign \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "content_type": "image/webp",
    "extension": "webp"
  }'
```

**Response**:
```json
{
  "success": true,
  "message": "Presigned URL generated successfully",
  "data": {
    "upload_url": "http://localhost:9000/greenlight/products/2024/12/prod-123/550e8400-e29b-41d4-a716-446655440000.webp?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minioadmin%2F20241221%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20241221T100000Z&X-Amz-Expires=600&X-Amz-SignedHeaders=host&X-Amz-Signature=...",
    "public_url": "http://localhost:9000/greenlight/products/2024/12/prod-123/550e8400-e29b-41d4-a716-446655440000.webp",
    "object_key": "products/2024/12/prod-123/550e8400-e29b-41d4-a716-446655440000.webp",
    "expires_at": "2024-12-21T10:10:00Z"
  }
}
```

**Allowed Content Types**:
- `image/jpeg`
- `image/jpg`
- `image/png`
- `image/webp`

**Allowed Extensions**:
- `jpg`, `jpeg`, `png`, `webp`

---

### Step 2: Upload Image to MinIO

**Request**:
```bash
curl -X PUT "http://localhost:9000/greenlight/products/2024/12/prod-123/550e8400-e29b-41d4-a716-446655440000.webp?X-Amz-Algorithm=..." \
  -H "Content-Type: image/webp" \
  --data-binary "@/path/to/image.webp"
```

**Response**: HTTP 200 OK (empty body)

**JavaScript Example**:
```javascript
const uploadImage = async (uploadUrl, imageFile) => {
  const response = await fetch(uploadUrl, {
    method: 'PUT',
    headers: {
      'Content-Type': imageFile.type
    },
    body: imageFile
  });
  
  if (!response.ok) {
    throw new Error('Upload failed');
  }
  
  return true;
};

// Usage
const fileInput = document.getElementById('imageInput');
const file = fileInput.files[0];
await uploadImage(uploadUrl, file);
```

**Python Example**:
```python
import requests

def upload_image(upload_url, image_path):
    with open(image_path, 'rb') as f:
        response = requests.put(
            upload_url,
            data=f,
            headers={'Content-Type': 'image/webp'}
        )
    return response.status_code == 200

# Usage
success = upload_image(upload_url, '/path/to/image.webp')
```

---

### Step 3: Confirm Upload

**Request**:
```bash
curl -X POST http://localhost:8080/api/v1/products/prod-123/images \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "object_key": "products/2024/12/prod-123/550e8400-e29b-41d4-a716-446655440000.webp",
    "is_main": true,
    "sort_order": 0
  }'
```

**Response**:
```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "data": {
    "image_id": "img-abc123",
    "product_id": "prod-123",
    "url": "http://localhost:9000/greenlight/products/2024/12/prod-123/550e8400-e29b-41d4-a716-446655440000.webp",
    "is_main": true,
    "sort_order": 0,
    "status": "ACTIVE",
    "created_at": "2024-12-21T10:05:00Z"
  }
}
```

---

## 📋 List Product Images

**Request**:
```bash
curl -X GET http://localhost:8080/api/v1/products/prod-123/images
```

**Response**:
```json
{
  "success": true,
  "message": "Product images retrieved successfully",
  "data": [
    {
      "image_id": "img-abc123",
      "product_id": "prod-123",
      "url": "http://localhost:9000/greenlight/products/2024/12/prod-123/image1.webp",
      "is_main": true,
      "sort_order": 0,
      "status": "ACTIVE",
      "created_at": "2024-12-21T10:05:00Z"
    },
    {
      "image_id": "img-def456",
      "product_id": "prod-123",
      "url": "http://localhost:9000/greenlight/products/2024/12/prod-123/image2.webp",
      "is_main": false,
      "sort_order": 1,
      "status": "ACTIVE",
      "created_at": "2024-12-21T10:06:00Z"
    }
  ]
}
```

---

## ✏️ Update Image Metadata

**Request**:
```bash
curl -X PATCH http://localhost:8080/api/v1/products/prod-123/images/img-abc123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "is_main": false,
    "sort_order": 2
  }'
```

**Response**:
```json
{
  "success": true,
  "message": "Image metadata updated successfully",
  "data": {
    "image_id": "img-abc123",
    "product_id": "prod-123",
    "url": "http://localhost:9000/greenlight/products/2024/12/prod-123/image1.webp",
    "is_main": false,
    "sort_order": 2,
    "status": "ACTIVE",
    "created_at": "2024-12-21T10:05:00Z"
  }
}
```

**Notes**:
- Chỉ cần gửi fields muốn update
- Nếu set `is_main: true`, các ảnh khác sẽ tự động set `is_main: false`

---

## 🗑️ Delete Image

**Request**:
```bash
curl -X DELETE http://localhost:8080/api/v1/products/prod-123/images/img-abc123 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Response**:
```json
{
  "success": true,
  "message": "Image deleted successfully",
  "data": null
}
```

**Notes**:
- Xóa cả object trong MinIO và metadata trong MySQL
- Nếu xóa MinIO thất bại, vẫn xóa metadata (log warning)

---

## 🔗 Complete Frontend Example (React)

```jsx
import React, { useState } from 'react';

const ProductImageUpload = ({ productId, authToken }) => {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState(null);
  const [images, setImages] = useState([]);

  const uploadImage = async (file, isMain = false) => {
    try {
      setUploading(true);
      setError(null);

      // Step 1: Get presigned URL
      const presignRes = await fetch(
        `/api/v1/products/${productId}/images/presign`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${authToken}`
          },
          body: JSON.stringify({
            content_type: file.type,
            extension: file.name.split('.').pop()
          })
        }
      );

      if (!presignRes.ok) {
        throw new Error('Failed to get presigned URL');
      }

      const { data: presignData } = await presignRes.json();

      // Step 2: Upload to MinIO
      const uploadRes = await fetch(presignData.upload_url, {
        method: 'PUT',
        headers: {
          'Content-Type': file.type
        },
        body: file
      });

      if (!uploadRes.ok) {
        throw new Error('Failed to upload image');
      }

      // Step 3: Confirm upload
      const confirmRes = await fetch(
        `/api/v1/products/${productId}/images`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${authToken}`
          },
          body: JSON.stringify({
            object_key: presignData.object_key,
            is_main: isMain,
            sort_order: images.length
          })
        }
      );

      if (!confirmRes.ok) {
        throw new Error('Failed to confirm upload');
      }

      const { data: imageData } = await confirmRes.json();
      setImages([...images, imageData]);

      alert('Image uploaded successfully!');
    } catch (err) {
      setError(err.message);
      console.error('Upload error:', err);
    } finally {
      setUploading(false);
    }
  };

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      uploadImage(file, images.length === 0); // First image is main
    }
  };

  const deleteImage = async (imageId) => {
    try {
      const res = await fetch(
        `/api/v1/products/${productId}/images/${imageId}`,
        {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${authToken}`
          }
        }
      );

      if (!res.ok) {
        throw new Error('Failed to delete image');
      }

      setImages(images.filter(img => img.image_id !== imageId));
      alert('Image deleted successfully!');
    } catch (err) {
      setError(err.message);
      console.error('Delete error:', err);
    }
  };

  return (
    <div>
      <h2>Product Images</h2>
      
      <input
        type="file"
        accept="image/jpeg,image/png,image/webp"
        onChange={handleFileChange}
        disabled={uploading}
      />
      
      {uploading && <p>Uploading...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}
      
      <div style={{ display: 'flex', gap: '10px', marginTop: '20px' }}>
        {images.map(image => (
          <div key={image.image_id} style={{ position: 'relative' }}>
            <img
              src={image.url}
              alt=""
              style={{ width: '150px', height: '150px', objectFit: 'cover' }}
            />
            {image.is_main && (
              <span style={{
                position: 'absolute',
                top: 0,
                left: 0,
                background: 'green',
                color: 'white',
                padding: '2px 5px'
              }}>
                Main
              </span>
            )}
            <button
              onClick={() => deleteImage(image.image_id)}
              style={{
                position: 'absolute',
                top: 0,
                right: 0,
                background: 'red',
                color: 'white',
                border: 'none',
                cursor: 'pointer'
              }}
            >
              ×
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ProductImageUpload;
```

---

## 🔧 Complete Backend Example (Go)

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

type PresignRequest struct {
    ContentType string `json:"content_type"`
    Extension   string `json:"extension"`
}

type PresignResponse struct {
    Success bool `json:"success"`
    Data    struct {
        UploadURL string `json:"upload_url"`
        PublicURL string `json:"public_url"`
        ObjectKey string `json:"object_key"`
        ExpiresAt string `json:"expires_at"`
    } `json:"data"`
}

type ConfirmRequest struct {
    ObjectKey string `json:"object_key"`
    IsMain    bool   `json:"is_main"`
    SortOrder int    `json:"sort_order"`
}

func uploadProductImage(productID, imagePath, token string) error {
    // Step 1: Get presigned URL
    presignReq := PresignRequest{
        ContentType: "image/jpeg",
        Extension:   "jpg",
    }
    
    presignJSON, _ := json.Marshal(presignReq)
    
    req, _ := http.NewRequest(
        "POST",
        fmt.Sprintf("http://localhost:8080/api/v1/products/%s/images/presign", productID),
        bytes.NewBuffer(presignJSON),
    )
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var presignResp PresignResponse
    json.NewDecoder(resp.Body).Decode(&presignResp)
    
    // Step 2: Upload to MinIO
    imageFile, _ := os.Open(imagePath)
    defer imageFile.Close()
    
    uploadReq, _ := http.NewRequest("PUT", presignResp.Data.UploadURL, imageFile)
    uploadReq.Header.Set("Content-Type", "image/jpeg")
    
    uploadResp, err := client.Do(uploadReq)
    if err != nil {
        return err
    }
    defer uploadResp.Body.Close()
    
    if uploadResp.StatusCode != 200 {
        return fmt.Errorf("upload failed: %d", uploadResp.StatusCode)
    }
    
    // Step 3: Confirm upload
    confirmReq := ConfirmRequest{
        ObjectKey: presignResp.Data.ObjectKey,
        IsMain:    true,
        SortOrder: 0,
    }
    
    confirmJSON, _ := json.Marshal(confirmReq)
    
    confirmHTTPReq, _ := http.NewRequest(
        "POST",
        fmt.Sprintf("http://localhost:8080/api/v1/products/%s/images", productID),
        bytes.NewBuffer(confirmJSON),
    )
    confirmHTTPReq.Header.Set("Content-Type", "application/json")
    confirmHTTPReq.Header.Set("Authorization", "Bearer "+token)
    
    confirmResp, err := client.Do(confirmHTTPReq)
    if err != nil {
        return err
    }
    defer confirmResp.Body.Close()
    
    fmt.Println("Image uploaded successfully!")
    return nil
}

func main() {
    token := "your-jwt-token"
    productID := "prod-123"
    imagePath := "./test-image.jpg"
    
    err := uploadProductImage(productID, imagePath, token)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

---

## ❌ Error Responses

### Invalid Content Type
```json
{
  "success": false,
  "error": "invalid content type: image/gif (allowed: jpeg, png, webp)"
}
```

### Product Not Found
```json
{
  "success": false,
  "error": "product not found"
}
```

### Object Not Found (after upload)
```json
{
  "success": false,
  "error": "object not found in storage"
}
```

### Unauthorized
```json
{
  "success": false,
  "error": "unauthorized"
}
```

### Presigned URL Expired
```json
{
  "success": false,
  "error": "Request has expired"
}
```

---

## 📊 Rate Limits

Currently no rate limits implemented. TODO for production:
- Max 10 presigned URLs per minute per user
- Max 50 images per product
- Max 100 uploads per day per user

---

## 🎯 Best Practices

1. **Validate file before requesting presigned URL**
   ```javascript
   const validateImage = (file) => {
     const allowedTypes = ['image/jpeg', 'image/png', 'image/webp'];
     const maxSize = 7 * 1024 * 1024; // 7MB
     
     if (!allowedTypes.includes(file.type)) {
       throw new Error('Invalid file type');
     }
     
     if (file.size > maxSize) {
       throw new Error('File too large');
     }
   };
   ```

2. **Handle presigned URL expiration**
   ```javascript
   const uploadWithRetry = async (file) => {
     try {
       await uploadImage(file);
     } catch (err) {
       if (err.message.includes('expired')) {
         // Request new presigned URL and retry
         await uploadImage(file);
       }
     }
   };
   ```

3. **Show upload progress**
   ```javascript
   const uploadWithProgress = (uploadUrl, file, onProgress) => {
     return new Promise((resolve, reject) => {
       const xhr = new XMLHttpRequest();
       
       xhr.upload.addEventListener('progress', (e) => {
         if (e.lengthComputable) {
           const percentComplete = (e.loaded / e.total) * 100;
           onProgress(percentComplete);
         }
       });
       
       xhr.addEventListener('load', () => {
         if (xhr.status === 200) {
           resolve();
         } else {
           reject(new Error('Upload failed'));
         }
       });
       
       xhr.open('PUT', uploadUrl);
       xhr.setRequestHeader('Content-Type', file.type);
       xhr.send(file);
     });
   };
   ```

4. **Optimize images before upload**
   ```javascript
   import imageCompression from 'browser-image-compression';
   
   const optimizeImage = async (file) => {
     const options = {
       maxSizeMB: 1,
       maxWidthOrHeight: 1920,
       useWebWorker: true
     };
     
     return await imageCompression(file, options);
   };
   ```

---

## 🔗 Related Endpoints

- `POST /api/v1/products` - Create product
- `GET /api/v1/products/:id_or_slug` - Get product (includes images)
- `PUT /api/v1/products/:id_or_slug` - Update product
- `DELETE /api/v1/products/:id_or_slug` - Delete product (cascade deletes images)

---

**For more details, see**: [docs/MINIO_INTEGRATION.md](docs/MINIO_INTEGRATION.md)

