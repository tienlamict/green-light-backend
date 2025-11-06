# Quick Reference Card

## 🚀 Getting Started (30 seconds)

```bash
# 1. Start with Docker
docker-compose up -d

# 2. Wait for services (30 sec)
sleep 30

# 3. Seed database (run from project root)
docker-compose exec -w /root api go run scripts/seed.go

# 4. Done! API running at http://localhost:8080
```

---

## 🔑 Default Credentials

```
Email: admin@example.com
Password: admin123
```

---

## 📝 Essential Commands

### Docker
```bash
make docker-up        # Start containers
make docker-down      # Stop containers
make docker-logs      # View logs
docker-compose exec api sh  # Shell into container
```

### Local Development
```bash
make dev              # Run locally
make build            # Build binary
make seed             # Seed database
make test             # Run all tests
```

### Database
```bash
# Auto-migrations run on startup (in correct order)
# Manual migration:
make migrate

# Reset database:
./scripts/reset_db.sh    # Linux/Mac
scripts\reset_db.bat     # Windows

# Seed data:
go run scripts/seed.go
```

---

## 🌐 Quick API Tests

### 1. Health Check
```bash
curl http://localhost:8080/healthz
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

### 3. Get Products
```bash
curl http://localhost:8080/api/v1/products
```

### 4. Create Product (use token from step 2)
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "sku": "TEST-001",
    "price": 99.99,
    "stock": 10,
    "category_id": "GET_FROM_CATEGORIES_ENDPOINT",
    "is_active": true
  }'
```

---

## 📍 Important URLs

| Resource | URL |
|----------|-----|
| API Base | http://localhost:8080 |
| Swagger Docs | http://localhost:8080/api/docs/index.html |
| Health Check | http://localhost:8080/healthz |
| Products | http://localhost:8080/api/v1/products |
| Categories | http://localhost:8080/api/v1/categories |

---

## 🔐 Authentication Flow

```bash
# 1. Get Token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.data.token')

# 2. Use Token
curl http://localhost:8080/api/v1/auth/user-info \
  -H "Authorization: Bearer $TOKEN"

# 3. Create/Update/Delete (protected routes)
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{ ... }'
```

---

## 🔍 Query Parameters

### Products
```bash
# Pagination
?page=1&limit=10

# Search
?q=headphone

# Filter by category
?category=CATEGORY_ID

# Filter by price
?min_price=100&max_price=500

# Filter by status
?is_active=true

# Sort
?sort=price ASC
?sort=created_at DESC
?sort=name ASC

# Combine all
?q=phone&category=ID&min_price=500&max_price=1000&page=1&limit=5&sort=price DESC
```

### Categories
```bash
# Search
?search=electronics

# Filter by status
?is_active=true

# Sort
?sort=name ASC
```

---

## 🗂️ Project Structure (Simplified)

```
green-light-backend/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── domain/                  # Entities & interfaces
│   ├── repository/              # Database layer
│   ├── usecase/                 # Business logic
│   └── transport/http/          # HTTP handlers
├── pkg/                         # Utilities
├── scripts/seed.go              # Database seeder
├── docker-compose.yml           # Docker setup
├── Makefile                     # Build commands
└── .env                         # Config (copy from .env.example)
```

---

## 🐛 Troubleshooting

### Database Connection Failed
```bash
# Check MySQL is running
docker-compose ps

# View logs
docker-compose logs db

# Restart
docker-compose restart db
```

### Port 8080 Already in Use
```bash
# Change port in .env
PORT=8081

# Or stop conflicting service
lsof -ti:8080 | xargs kill -9  # Mac/Linux
```

### Cannot Login
```bash
# Re-seed database
make seed

# Or with Docker
docker-compose exec api go run /root/scripts/seed.go
```

### Docker Issues
```bash
# Clean restart
docker-compose down -v
docker-compose up -d --build
sleep 30
docker-compose exec api go run /root/scripts/seed.go
```

---

## 📊 Seed Data

**Categories:**
- Electronics (slug: electronics)
- Furniture (slug: furniture)

**Products:**
- Wireless Bluetooth Headphones ($199.99)
- 4K Smart TV 55 inch ($699.99)
- Smartphone 128GB ($899.99)
- Modern Office Desk ($349.99)
- Ergonomic Office Chair ($279.99)
- Bookshelf Cabinet ($189.99)

**Users:**
- admin@example.com / admin123 (role: admin)

---

## 🎯 Common Tasks

### Get Category ID for Product Creation
```bash
# List categories
curl http://localhost:8080/api/v1/categories

# Note the category_id from response
```

### Upload Image
```bash
curl -X POST http://localhost:8080/api/v1/uploads \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@/path/to/image.jpg"

# Use returned URL in product
```

### Search Products
```bash
# By name/description/SKU
curl "http://localhost:8080/api/v1/products?q=headphone"

# By category slug (get products, then filter in app)
curl "http://localhost:8080/api/v1/products?category=CATEGORY_ID"
```

### Pagination
```bash
# Page 2, 5 items per page
curl "http://localhost:8080/api/v1/products?page=2&limit=5"
```

---

## 🧪 Run Tests

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests (requires test DB)
make test-int

# With coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🔄 Workflow Example

### Complete Product Creation Flow
```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.data.token')

# 2. Get Categories
CATEGORIES=$(curl -s http://localhost:8080/api/v1/categories)
echo $CATEGORIES | jq '.data[0].category_id'

# 3. Upload Image (optional)
IMAGE_URL=$(curl -s -X POST http://localhost:8080/api/v1/uploads \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@product.jpg" | jq -r '.data.url')

# 4. Create Product
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"New Product\",
    \"sku\": \"NP-001\",
    \"price\": 199.99,
    \"stock\": 50,
    \"thumbnail_url\": \"$IMAGE_URL\",
    \"category_id\": \"PASTE_CATEGORY_ID_HERE\",
    \"is_active\": true
  }"

# 5. Verify
curl http://localhost:8080/api/v1/products | jq
```

---

## 🛠️ Environment Variables

```bash
# Required
PORT=8080
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=root
DB_NAME=greenlight_db
JWT_SECRET=change-me

# Optional (have defaults)
JWT_EXPIRE_HOURS=24
CORS_ORIGINS=*
UPLOAD_DIR=./uploads
LOG_LEVEL=info
```

---

## 📚 Documentation Files

- **README.md** - Overview
- **GETTING_STARTED.md** - Detailed setup
- **API_EXAMPLES.md** - API usage examples
- **PROJECT_SUMMARY.md** - Complete project details
- **QUICK_REFERENCE.md** - This file

---

## ⚡ One-Line Commands

```bash
# Full setup
docker-compose up -d && sleep 30 && docker-compose exec -w /root api go run scripts/seed.go

# Get token
curl -s -X POST localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"admin@example.com","password":"admin123"}' | jq -r '.data.token'

# Test auth
curl localhost:8080/api/v1/auth/user-info -H "Authorization: Bearer $(curl -s -X POST localhost:8080/api/v1/auth/login -H 'Content-Type: application/json' -d '{"email":"admin@example.com","password":"admin123"}' | jq -r '.data.token')"
```

---

## 📱 Postman/Insomnia Import

Use these base settings:
- **Base URL:** `{{baseUrl}}` = http://localhost:8080
- **Auth Token:** `{{token}}` = Bearer YOUR_JWT_TOKEN
- **Header:** Authorization: Bearer {{token}}

---

## ✅ Checklist

**First Time Setup:**
- [ ] Clone repository
- [ ] Copy .env.example to .env (if not present)
- [ ] Run `docker-compose up -d`
- [ ] Wait 30 seconds
- [ ] Seed database
- [ ] Test health endpoint
- [ ] Login and save token
- [ ] View Swagger docs

**Daily Development:**
- [ ] Start Docker: `docker-compose up -d`
- [ ] Check logs: `docker-compose logs -f api`
- [ ] Make changes
- [ ] Test: `make test`
- [ ] Stop: `docker-compose down`

---

**Made with ❤️ using Go, Gin, GORM, and Clean Architecture**

