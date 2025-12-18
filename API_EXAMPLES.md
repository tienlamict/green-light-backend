# API Examples

Complete examples for testing the Green Light Backend API.

## Table of Contents
- [Authentication](#authentication)
- [Products](#products)
- [Categories](#categories)
- [File Upload](#file-upload)
- [Advanced Queries](#advanced-queries)

---

## Authentication

### Login

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "user_id": "018e1234-5678-7abc-def0-123456789abc",
      "email": "admin@example.com",
      "role": "admin",
      "created_at": "2025-11-06T10:00:00Z"
    }
  }
}
```

### Get User Info

**Request:**
```bash
curl http://localhost:8080/api/v1/auth/user-info \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "User info retrieved",
  "data": {
    "user_id": "018e1234-5678-7abc-def0-123456789abc",
    "email": "admin@example.com",
    "role": "admin",
    "created_at": "2025-11-06T10:00:00Z"
  }
}
```

---

## Products

### List All Products

**Request:**
```bash
curl "http://localhost:8080/api/v1/products?page=1&limit=10"
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "product_id": "018e1234-5678-7abc-def0-123456789abc",
      "name": "Wireless Bluetooth Headphones",
      "slug": "wireless-bluetooth-headphones",
      "sku": "ELEC-HP-001",
      "short_desc": "Premium wireless headphones with active noise cancellation",
      "description": "Experience superior sound quality...",
      "price": 199.99,
      "stock": 50,
      "thumbnail_url": "/uploads/headphones.jpg",
      "gallery": ["/uploads/headphones-1.jpg", "/uploads/headphones-2.jpg"],
      "category_id": "018e1234-5678-7abc-def0-category123",
      "category": {
        "category_id": "018e1234-5678-7abc-def0-category123",
        "name": "Electronics",
        "slug": "electronics",
        "description": "Electronic devices and gadgets",
        "is_active": true,
        "created_at": "2025-11-06T10:00:00Z",
        "updated_at": "2025-11-06T10:00:00Z"
      },
      "is_active": true,
      "created_at": "2025-11-06T10:00:00Z",
      "updated_at": "2025-11-06T10:00:00Z"
    }
  ],
  "meta": {
    "total": 6,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

### Get Product by ID or Slug

**Request by Slug:**
```bash
curl http://localhost:8080/api/v1/products/wireless-bluetooth-headphones
```

**Request by ID:**
```bash
curl http://localhost:8080/api/v1/products/018e1234-5678-7abc-def0-123456789abc
```

### Create Product

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mechanical Gaming Keyboard",
    "slug": "mechanical-gaming-keyboard",
    "sku": "ELEC-KB-001",
    "short_desc": "RGB mechanical keyboard with Cherry MX switches",
    "description": "Premium gaming keyboard featuring RGB backlighting, Cherry MX Red switches, and programmable keys. Perfect for gamers and typists alike.",
    "price": 149.99,
    "stock": 75,
    "thumbnail_url": "/uploads/keyboard.jpg",
    "gallery": ["/uploads/keyboard-1.jpg", "/uploads/keyboard-2.jpg"],
    "category_id": "YOUR_CATEGORY_ID",
    "is_active": true
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Product created",
  "data": {
    "product_id": "018e5678-1234-7abc-def0-987654321abc",
    "name": "Mechanical Gaming Keyboard",
    "slug": "mechanical-gaming-keyboard",
    ...
  }
}
```

### Update Product

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/products/018e5678-1234-7abc-def0-987654321abc \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "price": 129.99,
    "stock": 100,
    "is_active": true
  }'
```

### Delete Product (Admin Only)

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/products/018e5678-1234-7abc-def0-987654321abc \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Product deleted"
}
```

---

## Categories

### List All Categories

**Request:**
```bash
curl http://localhost:8080/api/v1/categories
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "category_id": "018e1234-5678-7abc-def0-category123",
      "name": "Electronics",
      "slug": "electronics",
      "description": "Electronic devices and gadgets",
      "is_active": true,
      "created_at": "2025-11-06T10:00:00Z",
      "updated_at": "2025-11-06T10:00:00Z"
    },
    {
      "category_id": "018e1234-5678-7abc-def0-category456",
      "name": "Furniture",
      "slug": "furniture",
      "description": "Home and office furniture",
      "is_active": true,
      "created_at": "2025-11-06T10:00:00Z",
      "updated_at": "2025-11-06T10:00:00Z"
    }
  ],
  "meta": {
    "total": 2,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

### Get Category by Slug

**Request:**
```bash
curl http://localhost:8080/api/v1/categories/electronics
```

### Create Category

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Accessories",
    "slug": "accessories",
    "description": "Product accessories and add-ons",
    "is_active": true
  }'
```

### Update Category

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/categories/CATEGORY_ID \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Updated description",
    "is_active": true
  }'
```

### Delete Category (Admin Only)

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/categories/CATEGORY_ID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## File Upload

### Upload Image

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/uploads \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/your/image.jpg"
```

**Response:**
```json
{
  "success": true,
  "message": "File uploaded successfully",
  "data": {
    "url": "/uploads/a1b2c3d4-e5f6-7890-abcd-ef1234567890_1699200000.jpg"
  }
}
```

---

## Advanced Queries

### Search Products by Keyword

**Request:**
```bash
curl "http://localhost:8080/api/v1/products?q=headphone"
```

### Filter by Category

**Request:**
```bash
curl "http://localhost:8080/api/v1/products?category=018e1234-5678-7abc-def0-category123"
```

### Filter by Price Range

**Request:**
```bash
curl "http://localhost:8080/api/v1/products?min_price=100&max_price=500"
```

### Combined Filters

**Request:**
```bash
curl "http://localhost:8080/api/v1/products?q=phone&category=CATEGORY_ID&min_price=500&max_price=1000&page=1&limit=5&sort=price ASC"
```

### Filter Active Products Only

**Request:**
```bash
curl "http://localhost:8080/api/v1/products?is_active=true"
```

### Search Categories

**Request:**
```bash
curl "http://localhost:8080/api/v1/categories?search=electron"
```

### Sort Products

**By newest:**
```bash
curl "http://localhost:8080/api/v1/products?sort=created_at DESC"
```

**By price (low to high):**
```bash
curl "http://localhost:8080/api/v1/products?sort=price ASC"
```

**By price (high to low):**
```bash
curl "http://localhost:8080/api/v1/products?sort=price DESC"
```

**By name:**
```bash
curl "http://localhost:8080/api/v1/products?sort=name ASC"
```

---

## Error Responses

### 400 Bad Request

```json
{
  "success": false,
  "error": "Key: 'CreateProductRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### 401 Unauthorized

```json
{
  "success": false,
  "error": "Authorization header is required"
}
```

### 403 Forbidden

```json
{
  "success": false,
  "error": "Insufficient permissions"
}
```

### 404 Not Found

```json
{
  "success": false,
  "error": "Product not found"
}
```

### 500 Internal Server Error

```json
{
  "success": false,
  "error": "Internal server error"
}
```

---

## Postman Collection

You can import these examples into Postman:

1. Create a new collection
2. Add a variable `{{baseUrl}}` = `http://localhost:8080`
3. Add a variable `{{token}}` = your JWT token
4. Use `{{baseUrl}}` and `{{token}}` in your requests

Example:
```
{{baseUrl}}/api/v1/products
Authorization: Bearer {{token}}
```

---

## Testing with Scripts

### Bash Script Example

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

# Login and get token
TOKEN=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }' | jq -r '.data.token')

echo "Token: $TOKEN"

# List products
curl -s "$BASE_URL/api/v1/products" | jq

# Create product
curl -s -X POST "$BASE_URL/api/v1/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "sku": "TEST-001",
    "price": 99.99,
    "stock": 10,
    "category_id": "YOUR_CATEGORY_ID",
    "is_active": true
  }' | jq
```

---

## JavaScript/Fetch Example

```javascript
const baseUrl = 'http://localhost:8080';
let token = '';

// Login
async function login() {
  const response = await fetch(`${baseUrl}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: 'admin@example.com',
      password: 'admin123'
    })
  });
  
  const data = await response.json();
  token = data.data.token;
  console.log('Logged in, token:', token);
}

// Get products
async function getProducts() {
  const response = await fetch(`${baseUrl}/api/v1/products`);
  const data = await response.json();
  console.log('Products:', data.data);
}

// Create product
async function createProduct() {
  const response = await fetch(`${baseUrl}/api/v1/products`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify({
      name: 'New Product',
      sku: 'NP-001',
      price: 199.99,
      stock: 50,
      category_id: 'YOUR_CATEGORY_ID',
      is_active: true
    })
  });
  
  const data = await response.json();
  console.log('Created product:', data.data);
}

// Run
login().then(getProducts).then(createProduct);
```