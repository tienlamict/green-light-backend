# 🚀 Complete Build Guide

## ✅ Step-by-Step: Build and Run the Project

### Step 1: Clean Start (Fresh Build)

```bash
# Stop any running containers and remove volumes
docker-compose down -v

# Remove old images (optional but recommended)
docker-compose build --no-cache
```

### Step 2: Start Services

```bash
# Start everything in detached mode
docker-compose up -d
```

You should see:
```
Creating network "green-light-backend_greenlight_network" with driver "bridge"
Creating greenlight_db ... done
Creating greenlight_api ... done
```

### Step 3: Wait for Database

```bash
# Wait 30 seconds for MySQL to initialize
sleep 30

# Or watch the logs until you see "ready for connections"
docker-compose logs -f db
```

Wait until you see:
```
greenlight_db | [Server] /usr/sbin/mysqld: ready for connections.
```

Press `Ctrl+C` to exit logs.

### Step 4: Check Migration

```bash
# Check API logs
docker-compose logs api | tail -20
```

You should see:
```
{"level":"info","timestamp":"...","msg":"Connected to database successfully"}
{"level":"info","timestamp":"...","msg":"Database migration completed"}
{"level":"info","timestamp":"...","msg":"Server starting on :8080"}
```

### Step 5: Verify Tables Created

```bash
docker-compose exec db mysql -u root -proot greenlight_db -e "SHOW TABLES;"
```

Expected output:
```
+-------------------------+
| Tables_in_greenlight_db |
+-------------------------+
| categories              |
| products                |
| users                   |
+-------------------------+
```

### Step 6: Seed Sample Data

```bash
docker-compose exec api go run /root/scripts/seed.go
```

Expected output:
```
Starting database seeding...
Seeding users...
  Created admin user: admin@example.com
Seeding categories...
  Created category: Electronics
  Created category: Furniture
Seeding products...
  Created product: Wireless Bluetooth Headphones
  Created product: 4K Smart TV 55 inch
  Created product: Smartphone 128GB
  Created product: Modern Office Desk
  Created product: Ergonomic Office Chair
  Created product: Bookshelf Cabinet
Database seeding completed successfully!
```

### Step 7: Test the API

**Health Check:**
```bash
curl http://localhost:8080/healthz
```

Response:
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "ok"
  }
}
```

**Get Products:**
```bash
curl http://localhost:8080/api/v1/products
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

### Step 8: Access Swagger Docs

Open in browser:
```
http://localhost:8080/api/docs/index.html
```

---

## 🎯 One-Command Full Setup

```bash
docker-compose down -v && docker-compose up -d && sleep 30 && docker-compose exec api go run /root/scripts/seed.go && curl http://localhost:8080/healthz
```

---

## 🔧 Using Makefile

If you have `make` installed:

```bash
# Complete reset and rebuild
make docker-reset

# Seed data
make docker-seed

# View logs
make docker-logs

# Stop everything
make docker-down
```

---

## 🐛 Troubleshooting

### Container Won't Start

```bash
# Check status
docker-compose ps

# Check logs
docker-compose logs api
docker-compose logs db
```

### Port Already in Use

```bash
# Find what's using port 8080
# Linux/Mac:
lsof -i :8080

# Windows:
netstat -ano | findstr :8080

# Change port in docker-compose.yml:
ports:
  - "8081:8080"  # Use 8081 instead
```

### Database Connection Failed

```bash
# Check MySQL is running
docker-compose ps db

# Should show "Up (healthy)"

# If unhealthy, wait longer or restart
docker-compose restart db
sleep 30
```

### Migration Failed (Error 1824)

```bash
# Clean everything and start fresh
docker-compose down -v
docker-compose up -d
sleep 30
docker-compose exec api go run /root/scripts/seed.go
```

### Can't Seed Data

```bash
# Make sure API container is running
docker-compose ps

# Try running seed script with logs
docker-compose exec api sh -c "cd /root && go run scripts/seed.go"
```

---

## 📊 Verification Checklist

- [ ] Containers are running: `docker-compose ps`
- [ ] MySQL is healthy: Shows `Up (healthy)`
- [ ] API is running: Shows `Up`
- [ ] Tables created: 3 tables (users, categories, products)
- [ ] Data seeded: 1 user, 2 categories, 6 products
- [ ] Health check works: `curl http://localhost:8080/healthz`
- [ ] Can login: Returns JWT token
- [ ] Can get products: Returns product list
- [ ] Swagger works: http://localhost:8080/api/docs/index.html

---

## 🎉 Success!

If all steps completed successfully, your API is now running at:

- **API Base URL:** http://localhost:8080
- **Swagger Docs:** http://localhost:8080/api/docs/index.html
- **Health Check:** http://localhost:8080/healthz

**Default Credentials:**
- Email: `admin@example.com`
- Password: `admin123`

---

## 📝 Next Steps

1. **Test endpoints** using Swagger UI
2. **Create products** via API
3. **Upload images** using `/api/v1/uploads`
4. **Read documentation** in README.md
5. **Check API examples** in API_EXAMPLES.md

---

## 🛑 Stop Everything

```bash
# Stop containers
docker-compose down

# Stop and remove volumes (deletes database)
docker-compose down -v
```

---

## 💡 Useful Commands

```bash
# View real-time logs
docker-compose logs -f

# View only API logs
docker-compose logs -f api

# View only DB logs
docker-compose logs -f db

# Shell into API container
docker-compose exec api sh

# Shell into MySQL
docker-compose exec db mysql -u root -proot greenlight_db

# Check database contents
docker-compose exec db mysql -u root -proot greenlight_db -e "SELECT COUNT(*) FROM products;"

# Restart a service
docker-compose restart api

# Rebuild without cache
docker-compose build --no-cache api
```

---

## 🚀 Ready to Build!

Run this complete sequence:

```bash
# 1. Clean slate
docker-compose down -v

# 2. Build and start
docker-compose up -d

# 3. Wait for services
echo "Waiting 30 seconds for MySQL..."
sleep 30

# 4. Check status
docker-compose ps

# 5. Seed data
docker-compose exec api go run /root/scripts/seed.go

# 6. Test
curl http://localhost:8080/healthz
curl http://localhost:8080/api/v1/products | jq

# 7. Done!
echo "✅ API is running at http://localhost:8080"
echo "📚 Swagger docs at http://localhost:8080/api/docs/index.html"
```

Happy coding! 🎉

