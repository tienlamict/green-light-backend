# Troubleshooting Guide

Common issues and solutions for the Green Light Backend API.

---

## Database Issues

### Error 1824: Failed to open the referenced table

**Problem:**
```
Error 1824 (HY000): Failed to open the referenced table 'categories'
```

**Cause:** Tables with foreign keys were created before their parent tables.

**Solution:**
```bash
# Option 1: Drop and recreate (Docker)
docker-compose down -v
docker-compose up -d
sleep 30

# Option 2: Reset manually
./scripts/reset_db.sh    # Linux/Mac
scripts\reset_db.bat     # Windows

# Option 3: Drop tables manually
mysql -u root -p greenlight_db < migrations/drop_tables.sql
make migrate
```

**Prevention:** The app now migrates tables in correct order automatically. This shouldn't happen anymore.

See: `SOLUTION_SUMMARY.md` and `MIGRATION_GUIDE.md`

---

### Cannot connect to database

**Docker:**
```bash
# Check container status
docker-compose ps

# View MySQL logs
docker-compose logs db

# Restart MySQL
docker-compose restart db

# Check health
docker-compose exec db mysql -u root -proot -e "SELECT 1;"
```

**Local:**
```bash
# Test connection
mysql -h localhost -u root -p

# Check MySQL is running
# Mac:
brew services list | grep mysql

# Linux:
systemctl status mysql

# Windows:
net start MySQL
```

**Fix:**
```bash
# Update .env with correct credentials
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=greenlight_db
```

---

### Table doesn't exist

**Check:**
```bash
docker-compose exec db mysql -u root -proot greenlight_db -e "SHOW TABLES;"
```

**If empty:**
```bash
# Tables weren't created - run migration
make migrate

# Or restart app (auto-migration)
docker-compose restart api
```

---

### Foreign key constraint fails

**Problem:**
```
Error: Cannot add or update a child row: a foreign key constraint fails
```

**Cause:** Trying to insert a Product with non-existent category_id.

**Solution:**
```bash
# Make sure categories exist
curl http://localhost:8080/api/v1/categories

# Or seed database
make seed
```

---

## Docker Issues

### Port 8080 already in use

**Check what's using the port:**
```bash
# Linux/Mac
lsof -i :8080

# Windows
netstat -ano | findstr :8080
```

**Solutions:**

**Option 1: Change port**
```bash
# Edit .env
PORT=8081

# Edit docker-compose.yml
ports:
  - "8081:8080"
```

**Option 2: Kill process**
```bash
# Linux/Mac
lsof -ti:8080 | xargs kill -9

# Windows (find PID first, then)
taskkill /PID <PID> /F
```

---

### Docker containers won't start

**Check status:**
```bash
docker-compose ps
docker-compose logs
```

**Common fixes:**

**1. Network issues:**
```bash
docker network prune
docker-compose down
docker-compose up -d
```

**2. Volume issues:**
```bash
docker-compose down -v
docker volume prune
docker-compose up -d
```

**3. Build issues:**
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

---

### MySQL container unhealthy

**Check health:**
```bash
docker-compose ps
# Look for "health: starting" or "unhealthy"
```

**Solution:**
```bash
# Wait longer (MySQL takes 20-30 seconds)
sleep 30
docker-compose ps

# If still unhealthy, restart
docker-compose restart db
```

**Check logs:**
```bash
docker-compose logs db | tail -50
```

---

## API Issues

### 404 Not Found on all endpoints

**Possible causes:**

1. **Server not started:**
   ```bash
   docker-compose ps
   # Should show "Up"
   ```

2. **Wrong URL:**
   ```bash
   # Correct:
   curl http://localhost:8080/api/v1/products
   
   # Wrong:
   curl http://localhost:8080/products  # Missing /api/v1
   ```

3. **Port mismatch:**
   ```bash
   # Check .env
   cat .env | grep PORT
   
   # Should match URL
   ```

---

### 401 Unauthorized

**Problem:** Protected endpoints require authentication.

**Solution:**
```bash
# 1. Login first
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}' \
  | jq -r '.data.token')

# 2. Use token
curl http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN"
```

---

### 403 Forbidden

**Problem:** User doesn't have required role.

**Roles:**
- `admin` - Can do everything
- `editor` - Can create/update but not delete

**Example:**
```bash
# Only admin can delete
curl -X DELETE http://localhost:8080/api/v1/products/ID \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

### Cannot login (invalid credentials)

**Check:**
```bash
# 1. Database seeded?
docker-compose exec db mysql -u root -proot greenlight_db \
  -e "SELECT email, role FROM users;"
```

**If empty:**
```bash
make seed
# or
docker-compose exec api go run /root/scripts/seed.go
```

**Default credentials:**
- Email: `admin@example.com`
- Password: `admin123`

---

### Validation errors (400 Bad Request)

**Problem:**
```json
{
  "success": false,
  "error": "Key: 'CreateProductRequest.Name' Error:Field validation for 'Name' failed"
}
```

**Solution:** Check required fields:

**Product:**
- name (required, min 2, max 255)
- sku (required)
- price (required, > 0)
- stock (required, >= 0)
- category_id (required, must exist)

**Category:**
- name (required, min 2, max 255)

**Login:**
- email (required, valid email)
- password (required, min 6)

---

## Build Issues

### Go module errors

**Problem:**
```
go: module not found
```

**Solution:**
```bash
# Download dependencies
go mod download

# Tidy modules
go mod tidy

# Verify
go mod verify
```

---

### Cannot find package

**Problem:**
```
package green-light-backend/internal/domain is not in GOROOT
```

**Solution:**
```bash
# Make sure you're in project root
cd /path/to/green-light-backend

# Ensure go.mod exists
ls go.mod

# Download dependencies
go mod download
```

---

### Swagger docs not working

**Problem:** `/api/docs/index.html` returns 404

**Solution:**
```bash
# Generate swagger docs
make swagger

# Or install swag and generate
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs

# Rebuild
make build
make dev
```

---

## Performance Issues

### Slow queries

**Enable query logging:**
```go
// In pkg/database/database.go, change:
Logger: logger.Default.LogMode(logger.Info)
```

**Check slow queries:**
```bash
docker-compose exec db mysql -u root -proot -e "
  SELECT * FROM information_schema.PROCESSLIST 
  WHERE TIME > 5 AND COMMAND != 'Sleep';
"
```

**Add indexes:**
```sql
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_name ON products(name);
```

---

### High memory usage

**Check:**
```bash
docker stats
```

**Solutions:**

1. **Reduce connection pool:**
   ```go
   // In pkg/database/database.go
   sqlDB.SetMaxIdleConns(5)
   sqlDB.SetMaxOpenConns(25)
   ```

2. **Add pagination:**
   ```bash
   # Always use limit
   curl "http://localhost:8080/api/v1/products?limit=10"
   ```

---

## Testing Issues

### Tests fail

**Unit tests:**
```bash
# Run with verbose
go test -v ./internal/usecase/

# Run single test
go test -v -run TestAuthUseCase_Login ./internal/usecase/
```

**Integration tests:**
```bash
# Need test database
mysql -u root -p -e "CREATE DATABASE greenlight_test;"

# Run with tag
go test -v -tags=integration ./internal/repository/
```

---

## File Upload Issues

### Cannot upload files

**Check:**
1. **Authenticated?**
   ```bash
   # Need Bearer token
   curl -X POST http://localhost:8080/api/v1/uploads \
     -H "Authorization: Bearer $TOKEN" \
     -F "file=@image.jpg"
   ```

2. **File too large?**
   ```bash
   # Check .env
   MAX_UPLOAD_SIZE=10485760  # 10MB
   ```

3. **Invalid file type?**
   - Allowed: jpg, jpeg, png, gif, webp

4. **Upload directory writable?**
   ```bash
   # Check permissions
   ls -la uploads/
   
   # Create if missing
   mkdir -p uploads
   chmod 755 uploads
   ```

---

## Common Fixes

### Reset Everything

**Complete reset:**
```bash
# Stop everything
docker-compose down -v

# Clean Docker
docker system prune -f

# Start fresh
docker-compose up -d --build
sleep 30

# Seed
docker-compose exec api go run /root/scripts/seed.go

# Test
curl http://localhost:8080/healthz
```

---

### Check Logs

**Application logs:**
```bash
# Docker
docker-compose logs -f api

# Local
# Check console output
```

**Database logs:**
```bash
docker-compose logs -f db
```

**All logs:**
```bash
docker-compose logs -f
```

---

### Verify Setup

**Checklist:**
```bash
# 1. Containers running?
docker-compose ps

# 2. Database accessible?
docker-compose exec db mysql -u root -proot -e "SELECT 1;"

# 3. Tables exist?
docker-compose exec db mysql -u root -proot greenlight_db -e "SHOW TABLES;"

# 4. Data seeded?
docker-compose exec db mysql -u root -proot greenlight_db \
  -e "SELECT COUNT(*) FROM products;"

# 5. API responding?
curl http://localhost:8080/healthz

# 6. Can login?
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

---

## Get Help

### Enable Debug Logging

```bash
# Edit .env
LOG_LEVEL=debug

# Restart
docker-compose restart api
```

### View Error Details

```bash
# Check API logs for full error
docker-compose logs api | grep -i error

# Check MySQL errors
docker-compose logs db | grep -i error
```

### Useful Commands

```bash
# Shell into API container
docker-compose exec api sh

# Shell into MySQL container
docker-compose exec db mysql -u root -proot greenlight_db

# Check environment variables
docker-compose exec api env | grep DB_

# Test network
docker-compose exec api nc -zv db 3306
```

---

## Still Need Help?

Check these files:
- `SOLUTION_SUMMARY.md` - Migration issue fix
- `MIGRATION_GUIDE.md` - Database migration help
- `GETTING_STARTED.md` - Setup instructions
- `API_EXAMPLES.md` - API usage examples
- `QUICK_REFERENCE.md` - Quick commands

Or review logs for specific error messages.

