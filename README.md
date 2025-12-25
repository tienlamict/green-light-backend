# Green Light Backend - Product Showcase API

A production-ready REST API backend for a product showcase website built with Go, following Clean Architecture principles.

## 🚀 Tech Stack

- **Language:** Go 1.24.2
- **Framework:** Gin
- **ORM:** GORM
- **Database:** MySQL 8.0
- **Architecture:** Clean Architecture
- **Auth:** JWT (HS256)
- **Logging:** Uber Zap
- **Deployment:** Docker + Docker Compose

## 📁 Project Structure

```
green-light-backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/                  # Domain layer (entities & interfaces)
│   │   ├── category.go
│   │   ├── product.go
│   │   └── user.go
│   ├── repository/              # Repository layer (GORM implementations)
│   │   ├── category_repository.go
│   │   ├── product_repository.go
│   │   └── user_repository.go
│   ├── usecase/                 # Use case layer (business logic)
│   │   ├── auth_usecase.go
│   │   ├── category_usecase.go
│   │   └── product_usecase.go
│   └── transport/               # Transport layer (HTTP handlers)
│       ├── http/
│       │   ├── handler/
│       │   │   ├── auth_handler.go
│       │   │   ├── category_handler.go
│       │   │   ├── product_handler.go
│       │   │   └── upload_handler.go
│       │   ├── middleware/
│       │   │   ├── auth.go
│       │   │   ├── cors.go
│       │   │   └── logger.go
│       │   ├── dto/
│       │   │   ├── auth_dto.go
│       │   │   ├── category_dto.go
│       │   │   ├── product_dto.go
│       │   │   └── response.go
│       │   └── router.go
├── pkg/
│   ├── config/                  # Configuration
│   ├── database/                # Database connection
│   ├── logger/                  # Logger setup
│   └── utils/                   # Utilities (JWT, password, etc.)
├── migrations/                  # Database migrations
├── scripts/                     # Helper scripts
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── .env.example
```

## 🛠️ Setup & Installation

### Prerequisites

- Go 1.24.2+
- Docker & Docker Compose
- Make (optional but recommended)

### Quick Start with Docker

1. Clone the repository:
```bash
git clone <repository-url>
cd green-light-backend
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Start the application:
```bash
docker-compose up -d
```

The API will be available at `http://localhost:8080`

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your local configuration
```

3. Start MySQL:
```bash
docker-compose up -d db
```

4. Run migrations and seed data:
```bash
make migrate
make seed
```

5. Start the server:
```bash
make dev
```

## 📝 Makefile Commands

```bash
make dev          # Run application in development mode
make build        # Build binary
make migrate      # Run database migrations
make seed         # Seed sample data
make test         # Run tests
make test-unit    # Run unit tests only
make test-int     # Run integration tests only
make docker-build # Build Docker image
make docker-up    # Start Docker containers
make docker-down  # Stop Docker containers
make clean        # Clean build artifacts
```

## 🔐 Authentication

The API uses JWT (HS256) for authentication.

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "user_id": "...",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

### Authenticated Requests

Include the JWT token in the Authorization header:

```bash
curl -X GET http://localhost:8080/api/v1/auth/user-info \
  -H "Authorization: Bearer <your-token>"
```

## 🌐 API Endpoints

### Auth
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/user-info` - Get user profile (requires auth)

### Products
- `GET /api/v1/products` - List products (with filters)
- `GET /api/v1/products/:id_or_slug` - Get product details
- `POST /api/v1/products` - Create product (admin/editor)
- `PUT /api/v1/products/:id` - Update product (admin/editor)
- `DELETE /api/v1/products/:id` - Delete product (admin only)

### Categories
- `GET /api/v1/categories` - List categories
- `GET /api/v1/categories/:id_or_slug` - Get category details
- `POST /api/v1/categories` - Create category (admin/editor)
- `PUT /api/v1/categories/:id` - Update category (admin/editor)
- `DELETE /api/v1/categories/:id` - Delete category (admin only)

### File Upload
- `POST /api/v1/uploads` - Upload image (authenticated)

### Health
- `GET /healthz` - Health check

## 🧪 Testing

Run all tests:
```bash
make test
```

Run only unit tests:
```bash
make test-unit
```

Run only integration tests:
```bash
make test-int
```

## 🗄️ Database

### Seed Data

The database is seeded with:
- 2 categories (Electronics, Furniture)
- 6 products
- 1 admin user (email: `admin@example.com`, password: `admin123`)

### Migrations

Migrations are automatically applied when the application starts in the **correct order** to avoid foreign key issues:
1. Users table (no dependencies)
2. Categories table (no dependencies)
3. Products table (depends on Categories)

**Manual migration:**
```bash
make migrate
```

**Reset database (clean slate):**
```bash
# Linux/Mac
chmod +x scripts/reset_db.sh
./scripts/reset_db.sh

# Windows
scripts\reset_db.bat

# Docker
docker-compose down -v
docker-compose up -d
```

See `MIGRATION_GUIDE.md` for detailed information.

## 🎯 Features

- ✅ Clean Architecture (domain, usecase, repository, transport)
- ✅ JWT Authentication & Authorization
- ✅ Role-based Access Control (admin, editor)
- ✅ Search & Filtering (name, category, price range)
- ✅ Pagination with total count
- ✅ Unified error handling
- ✅ Structured logging (Zap)
- ✅ CORS middleware
- ✅ Request validation
- ✅ Docker support
- ✅ Unit & Integration tests
- ✅ File upload support
- ✅ Health check endpoint

## 🐳 Docker

Build and run with Docker Compose:

```bash
docker-compose up --build
```

Stop containers:

```bash
docker-compose down
```

Remove volumes:

```bash
docker-compose down -v
```

## 📖 License

MIT License

