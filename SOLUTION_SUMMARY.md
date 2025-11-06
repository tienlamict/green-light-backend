# Solution: Fixed MySQL Error 1824

## Problem
```
Error 1824 (HY000): Failed to open the referenced table
```

This error occurred because tables with foreign key constraints were being created before their parent tables existed.

## Root Cause

The `Products` table has a foreign key constraint to the `Categories` table:
```sql
CONSTRAINT `fk_products_category` 
  FOREIGN KEY (`category_id`) 
  REFERENCES `categories`(`category_id`)
```

When GORM's `AutoMigrate` was called with:
```go
db.AutoMigrate(&domain.User{}, &domain.Category{}, &domain.Product{})
```

It sometimes tried to create `Products` before `Categories`, causing the foreign key constraint to fail.

## Solution Implemented

### 1. ✅ Fixed AutoMigrate Order (cmd/api/main.go)

Changed from:
```go
func migrateDatabase(db *gorm.DB) error {
    return db.AutoMigrate(
        &domain.User{},
        &domain.Category{},
        &domain.Product{},
    )
}
```

To:
```go
func migrateDatabase(db *gorm.DB) error {
    // Migrate tables in correct order to avoid foreign key issues
    // 1. Users (no dependencies)
    if err := db.AutoMigrate(&domain.User{}); err != nil {
        return fmt.Errorf("failed to migrate users table: %w", err)
    }

    // 2. Categories (no dependencies)
    if err := db.AutoMigrate(&domain.Category{}); err != nil {
        return fmt.Errorf("failed to migrate categories table: %w", err)
    }

    // 3. Products (depends on Categories)
    if err := db.AutoMigrate(&domain.Product{}); err != nil {
        return fmt.Errorf("failed to migrate products table: %w", err)
    }

    return nil
}
```

**Order matters!** Tables must be created from parents to children.

### 2. ✅ Created SQL Migration Script (migrations/001_create_tables.sql)

Complete SQL schema that creates tables in the correct order:
1. Users (no foreign keys)
2. Categories (no foreign keys)  
3. Products (foreign key to Categories)

### 3. ✅ Added Migration Runner (migrations/migrate.go)

Go program to run SQL migrations:
```bash
go run migrations/migrate.go
# or
make migrate
```

### 4. ✅ Added Reset Scripts

For development:
- `scripts/reset_db.sh` (Linux/Mac)
- `scripts/reset_db.bat` (Windows)
- `migrations/drop_tables.sql`

## How to Use

### Option 1: Automatic (Recommended)

Just start the app - migrations run automatically:

```bash
# With Docker
docker-compose up -d

# Or locally
make dev
```

The app now migrates tables sequentially in the correct order.

### Option 2: Manual SQL Migration

```bash
# Run SQL script directly
mysql -u root -p greenlight_db < migrations/001_create_tables.sql

# Or use Go migration tool
make migrate
```

### Option 3: Reset Database (Clean Slate)

**Linux/Mac:**
```bash
chmod +x scripts/reset_db.sh
./scripts/reset_db.sh
make seed
```

**Windows:**
```cmd
scripts\reset_db.bat
make seed
```

**Docker:**
```bash
docker-compose down -v
docker-compose up -d
sleep 30
make docker-seed
```

## Table Creation Order

**Critical Order:**
```
1. users       → No dependencies
2. categories  → No dependencies
3. products    → Depends on categories
```

## Verification

Check if tables were created:
```bash
# Local
mysql -u root -p greenlight_db -e "SHOW TABLES;"

# Docker
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

Check foreign keys:
```bash
docker-compose exec db mysql -u root -proot greenlight_db -e "
  SELECT 
    CONSTRAINT_NAME, 
    TABLE_NAME, 
    REFERENCED_TABLE_NAME 
  FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
  WHERE TABLE_SCHEMA = 'greenlight_db' 
  AND REFERENCED_TABLE_NAME IS NOT NULL;"
```

Expected:
```
+----------------------+------------+-----------------------+
| CONSTRAINT_NAME      | TABLE_NAME | REFERENCED_TABLE_NAME |
+----------------------+------------+-----------------------+
| fk_products_category | products   | categories            |
+----------------------+------------+-----------------------+
```

## Files Created/Modified

**New Files:**
- `migrations/001_create_tables.sql` - Complete schema
- `migrations/migrate.go` - Migration runner
- `migrations/drop_tables.sql` - Drop script
- `scripts/reset_db.sh` - Reset script (Unix)
- `scripts/reset_db.bat` - Reset script (Windows)
- `MIGRATION_GUIDE.md` - Detailed guide
- `SOLUTION_SUMMARY.md` - This file

**Modified Files:**
- `cmd/api/main.go` - Fixed migration order
- `Makefile` - Updated migrate target

## Testing the Fix

1. **Clean slate:**
   ```bash
   docker-compose down -v
   ```

2. **Start fresh:**
   ```bash
   docker-compose up -d
   ```

3. **Wait for MySQL:**
   ```bash
   sleep 30
   ```

4. **Check tables:**
   ```bash
   docker-compose exec db mysql -u root -proot greenlight_db -e "SHOW TABLES;"
   ```

5. **Seed data:**
   ```bash
   docker-compose exec api go run /root/scripts/seed.go
   ```

6. **Test API:**
   ```bash
   curl http://localhost:8080/healthz
   curl http://localhost:8080/api/v1/products
   ```

## Why This Works

**Before:**
- GORM AutoMigrate might create tables in any order
- Products table tried to create FK to Categories before Categories existed
- MySQL Error 1824

**After:**
- Sequential migration ensures parent tables exist first
- Categories created before Products
- Foreign key constraint succeeds

## Additional Resources

- **Full Guide:** See `MIGRATION_GUIDE.md`
- **Quick Commands:** See `QUICK_REFERENCE.md`
- **Setup Guide:** See `GETTING_STARTED.md`

## Quick Commands

```bash
# Start everything (auto-migration)
docker-compose up -d && sleep 30 && make docker-seed

# Manual migration
make migrate

# Reset database
./scripts/reset_db.sh  # Unix
scripts\reset_db.bat   # Windows

# Check tables
docker-compose exec db mysql -u root -proot -e "USE greenlight_db; SHOW TABLES;"
```

## Summary

✅ **Problem Solved:** Tables now create in correct order  
✅ **Migration Works:** Both automatic and manual methods  
✅ **Scripts Provided:** For reset and recovery  
✅ **Documented:** Complete migration guide included  
✅ **Tested:** Works with Docker and local MySQL  

**The Error 1824 issue is now fixed!** 🎉

