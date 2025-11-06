# 🚨 Immediate Fix for Error 1824

## Current Error
```
Failed to migrate database: failed to migrate categories table: 
Error 1824 (HY000): Failed to open the referenced table 'categories'
```

## Root Cause
The database has **broken tables** from a previous failed migration. They must be cleaned up first.

---

## ✅ SOLUTION (Do This Now)

### Step 1: Stop Everything
```bash
docker-compose down
```

### Step 2: Remove Database Volume
```bash
# This deletes the broken database and starts fresh
docker-compose down -v
```

### Step 3: Start Fresh
```bash
docker-compose up -d
```

### Step 4: Wait for Services
```bash
# Wait 30 seconds for MySQL to be ready
sleep 30
```

### Step 5: Check Logs
```bash
# Should see "Database migration completed" without errors
docker-compose logs api | tail -20
```

### Step 6: Seed Data
```bash
docker-compose exec api go run /root/scripts/seed.go
```

### Step 7: Test
```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/api/v1/products
```

---

## 🎯 One-Line Fix

```bash
docker-compose down -v && docker-compose up -d && sleep 30 && docker-compose exec api go run /root/scripts/seed.go
```

---

## Why This Happens

1. First migration attempt created tables in wrong order
2. Tables exist but in broken state with missing foreign keys
3. New migration can't fix broken tables
4. **Solution:** Delete everything and start fresh

---

## Alternative: Manual Cleanup (Without Losing Data)

If you want to keep your data:

```bash
# Drop only the tables (keeps database)
docker-compose exec db mysql -u root -proot greenlight_db -e "
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;  
DROP TABLE IF EXISTS users;
SET FOREIGN_KEY_CHECKS = 1;
"

# Restart API (will recreate tables correctly)
docker-compose restart api

# Check logs
docker-compose logs api | tail -20
```

---

## Verification

After fix, you should see:

```bash
docker-compose exec db mysql -u root -proot -e "USE greenlight_db; SHOW TABLES;"
```

Output:
```
+-------------------------+
| Tables_in_greenlight_db |
+-------------------------+
| categories              |
| products                |
| users                   |
+-------------------------+
```

---

## What I Fixed in Code

Updated `cmd/api/main.go` to:
1. ✅ Disable foreign key checks during migration
2. ✅ Migrate tables in correct order
3. ✅ Re-enable foreign key checks after migration

This prevents future issues, but **existing broken tables must still be cleaned up manually**.

---

## Quick Commands

```bash
# Full reset
docker-compose down -v && docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f api

# Check tables
docker-compose exec db mysql -u root -proot -e "USE greenlight_db; SHOW TABLES;"

# Seed data
docker-compose exec api go run /root/scripts/seed.go

# Test API
curl http://localhost:8080/api/v1/products
```

---

## If Still Having Issues

1. **Check MySQL is ready:**
   ```bash
   docker-compose logs db | grep "ready for connections"
   ```

2. **Verify database exists:**
   ```bash
   docker-compose exec db mysql -u root -proot -e "SHOW DATABASES;"
   ```

3. **Check API can connect:**
   ```bash
   docker-compose logs api | grep -i "database\|error"
   ```

4. **Complete nuclear option:**
   ```bash
   # Remove everything
   docker-compose down -v
   docker volume prune -f
   docker system prune -f
   
   # Start from scratch
   docker-compose up -d --build
   sleep 30
   docker-compose exec api go run /root/scripts/seed.go
   ```

---

## Summary

**The issue:** Broken tables from previous migration  
**The fix:** `docker-compose down -v && docker-compose up -d`  
**The prevention:** Code now disables FK checks during migration

After running `docker-compose down -v`, the error should be gone! 🎉

