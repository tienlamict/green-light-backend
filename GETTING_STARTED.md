# Getting Started Guide

This guide will help you set up and run the Green Light Backend API locally and with Docker.

## Prerequisites

- **Go 1.24.2+** - [Download Go](https://golang.org/dl/)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **MySQL 8.0+** (for local development without Docker)
- **Make** (optional but recommended)

## Quick Start with Docker (Recommended)

This is the fastest way to get the application running.

### 1. Clone and Setup

```bash
git clone <repository-url>
cd green-light-backend
```

### 2. Start Everything with Docker

```bash
docker-compose up -d
```

This will:
- Start MySQL database
- Build and start the API server
- Apply database migrations automatically

Wait about 30 seconds for services to fully initialize.

### 3. Seed the Database

```bash
# Wait for the containers to be healthy
sleep 10

# Seed sample data
make docker-seed
```

Or manually:
```bash
docker-compose exec api go run /root/scripts/seed.go
```

### 4. Test the API

The API is now running at `http://localhost:8080`

**Test health check:**
```bash
curl http://localhost:8080/healthz
```

**Login as admin:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

**View Swagger docs:**
Open `http://localhost:8080/api/docs/index.html` in your browser

### 5. Stop Docker

```bash
docker-compose down
```

To remove volumes (database data):
```bash
docker-compose down -v
```

---

## Local Development (Without Docker)

### 1. Install Dependencies

```bash
go mod download
```

### 2. Start MySQL

**Option A: Using Docker for MySQL only**
```bash
docker run -d \
  --name greenlight_mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=greenlight_db \
  mysql:8.0.32
```

**Option B: Use existing MySQL installation**
Make sure MySQL is running and create the database:
```sql
CREATE DATABASE greenlight_db;
```

### 3. Configure Environment

The `.env` file is already configured for local development. Update if needed:

```bash
# .env is already present, but you can modify:
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=greenlight_db
```

### 4. Run Migrations

Migrations run automatically on startup, or you can run the app once:

```bash
make dev
```

Press `Ctrl+C` after it starts.

### 5. Seed Database

```bash
make seed
```

### 6. Run the Application

```bash
make dev
```

The server will start on `http://localhost:8080`

---

## API Usage Examples

### Authentication

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

Save the token from the response.

**Get User Info:**
```bash
curl http://localhost:8080/api/v1/auth/user-info \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Products

**List Products:**
```bash
curl "http://localhost:8080/api/v1/products?page=1&limit=10"
```

**Get Product by Slug:**
```bash
curl http://localhost:8080/api/v1/products/wireless-bluetooth-headphones
```

**Search Products:**
```bash
curl "http://localhost:8080/api/v1/products?q=phone&min_price=500&max_price=1000"
```

**Create Product (requires auth):**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Product",
    "sku": "NP-001",
    "description": "A great product",
    "price": 199.99,
    "stock": 50,
    "category_id": "CATEGORY_ID_HERE",
    "is_active": true
  }'
```

**Update Product (requires auth):**
```bash
curl -X PUT http://localhost:8080/api/v1/products/PRODUCT_ID \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "price": 179.99,
    "stock": 100
  }'
```

**Delete Product (admin only):**
```bash
curl -X DELETE http://localhost:8080/api/v1/products/PRODUCT_ID \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Categories

**List Categories:**
```bash
curl http://localhost:8080/api/v1/categories
```

**Get Category:**
```bash
curl http://localhost:8080/api/v1/categories/electronics
```

**Create Category (requires auth):**
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Category",
    "slug": "new-category",
    "description": "Category description",
    "is_active": true
  }'
```

### File Upload

**Upload Image (requires auth):**
```bash
curl -X POST http://localhost:8080/api/v1/uploads \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -F "file=@/path/to/image.jpg"
```

---

## Swagger Documentation

Interactive API documentation is available at:

```
http://localhost:8080/api/docs/index.html
```

You can test all endpoints directly from the Swagger UI. Click "Authorize" and enter your JWT token as:
```
Bearer YOUR_TOKEN_HERE
```

---

## Default Credentials

**Admin User:**
- Email: `admin@example.com`
- Password: `admin123`

**Seeded Data:**
- 2 Categories (Electronics, Furniture)
- 6 Products (3 Electronics, 3 Furniture)

---

## Common Make Commands

```bash
make help          # Show all available commands
make dev           # Run in development mode
make build         # Build binary
make seed          # Seed database
make test          # Run all tests
make test-unit     # Run unit tests only
make docker-up     # Start with Docker
make docker-down   # Stop Docker containers
make swagger       # Regenerate Swagger docs
make clean         # Clean build artifacts
```

---

## Troubleshooting

### Database Connection Issues

**Error: "Failed to connect to database"**

- Make sure MySQL is running
- Check credentials in `.env` file
- For Docker: `docker-compose ps` to verify containers are running
- For local MySQL: `mysql -u root -p` to test connection

### Port Already in Use

**Error: "bind: address already in use"**

Change the port in `.env`:
```
PORT=8081
```

### Docker Issues

**Containers won't start:**
```bash
# View logs
docker-compose logs

# Restart containers
docker-compose restart

# Rebuild from scratch
docker-compose down -v
docker-compose up --build
```

### Migration Issues

Migrations run automatically. To reset the database:

```bash
# With Docker
docker-compose down -v
docker-compose up -d
make docker-seed

# Without Docker
mysql -u root -p
DROP DATABASE greenlight_db;
CREATE DATABASE greenlight_db;
exit
make dev  # Migrations run on startup
make seed
```

---

## Project Structure Overview

```
green-light-backend/
├── cmd/api/                    # Application entry point
├── internal/
│   ├── domain/                 # Domain entities and interfaces
│   ├── repository/             # Data access layer (GORM)
│   ├── usecase/                # Business logic
│   └── transport/http/         # HTTP handlers and middleware
├── pkg/                        # Shared packages
│   ├── config/                 # Configuration
│   ├── database/               # Database connection
│   ├── logger/                 # Logger setup
│   └── utils/                  # Utilities (JWT, password, etc.)
├── scripts/                    # Helper scripts (seeder)
├── docs/                       # Swagger documentation
├── uploads/                    # Uploaded files
├── docker-compose.yml          # Docker Compose configuration
├── Dockerfile                  # Docker image definition
├── Makefile                    # Build and run commands
└── .env                        # Environment variables
```

---

## Next Steps

1. **Explore the API** using Swagger UI
2. **Run Tests**: `make test`
3. **Modify Entities** in `internal/domain/`
4. **Add New Endpoints** in `internal/transport/http/handler/`
5. **Update Documentation**: Regenerate with `make swagger`

---

## Production Deployment

For production:

1. Update `.env` with production values
2. Change `JWT_SECRET` to a strong random value
3. Use a managed MySQL instance
4. Set `ENV=production`
5. Enable HTTPS
6. Configure proper CORS origins
7. Set up monitoring and logging

---

## Support

For issues or questions, please check:
- API Documentation: `/api/docs/index.html`
- Project README: `README.md`
- Test files for usage examples

Happy coding! 🚀

