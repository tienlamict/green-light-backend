# Database Migration Guide

This guide explains how to handle database migrations for the Green Light Backend API.

## Problem Solved

**Error:** `Error 1824 (HY000): Failed to open the referenced table`

This error occurs when tables with foreign key constraints are created before their referenced tables exist. We've solved this by:
1. Creating tables in the correct order (parent tables first)
2. Providing explicit SQL migration scripts
3. Fixing GORM AutoMigrate to migrate tables sequentially

---

## Migration Methods

### Method 1: Automatic (Default - Recommended)

The application automatically runs migrations on startup in the correct order:

```bash
# Just start the app (locally)
make dev

# Or with Docker
docker-compose up -d
```

**Migration Order:**
1. Users table (no foreign keys)
2. Categories table (no foreign keys)
3. Products table (foreign key to Categories)

---

### Method 2: Manual SQL Migration

Use the SQL script directly for more control:

**Local MySQL:**
```bash
# Run migration
mysql -h localhost -u root -p greenlight_db < migrations/001_create_tables.sql

# Or use the Go migration tool
go run migrations/migrate.go

# Or use Make
make migrate
```

**With Docker:**
```bash
# Connect to MySQL container
docker-compose exec db mysql -u root -proot greenlight_db < migrations/001_create_tables.sql
```

---

### Method 3: Reset Database (Clean Slate)

**⚠️ WARNING: This deletes all data!**

**Local:**
```bash
# Make script executable (first time only)
chmod +x scripts/reset_db.sh

# Run reset
./scripts/reset_db.sh

# Then seed
make seed
```

**With Docker:**
```bash
# Drop and recreate
docker-compose down -v
docker-compose up -d
sleep 30

# Seed data
docker-compose exec api go run /root/scripts/seed.go
```

---

## Migration Files

### migrations/001_create_tables.sql
Complete SQL schema for all tables in correct order:
- Creates Users table
- Creates Categories table  
- Creates Products table with foreign key to Categories

### migrations/migrate.go
Go program that executes SQL migrations:
```bash
go run migrations/migrate.go
```

### migrations/drop_tables.sql
Drops all tables in reverse order (for reset):
```bash
mysql -u root -p greenlight_db < migrations/drop_tables.sql
```

---

## Table Creation Order (Important!)

Always create in this order:

```
1. users          (no dependencies)
2. categories     (no dependencies)
3. products       (FK: category_id -> categories.category_id)
```

**Why?** MySQL requires parent tables to exist before creating child tables with foreign keys.

---

## Common Scenarios

### First Time Setup

```bash
# Option A: Let app auto-migrate
docker-compose up -d
sleep 30
make docker-seed

# Option B: Manual migration
make migrate
make seed
make dev
```

### Database Already Exists

```bash
# If you get Error 1824, tables exist in wrong order
# Option 1: Drop and recreate
docker-compose down -v
docker-compose up -d

# Option 2: Manual drop
mysql -u root -p greenlight_db < migrations/drop_tables.sql
make migrate
```

### Change Database Schema

1. Stop the application
2. Update domain models in `internal/domain/`
3. Update SQL in `migrations/001_create_tables.sql`
4. Reset database:
   ```bash
   ./scripts/reset_db.sh
   make seed
   ```

### Production Migration

For production, use proper migration tools:

**Recommended Tools:**
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)
- [atlas](https://atlasgo.io/)

**Example with golang-migrate:**
```bash
# Install
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir migrations -seq create_tables

# Run migration
migrate -path migrations -database "mysql://root:root@tcp(localhost:3306)/greenlight_db" up
```

---

## Troubleshooting

### Error 1824: Failed to open referenced table

**Cause:** Products table trying to reference Categories before it exists.

**Solution:**
```bash
# Drop tables in reverse order
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
SET FOREIGN_KEY_CHECKS = 1;

# Then recreate
make migrate
```

### Error 1050: Table already exists

**Cause:** Table already exists.

**Solution:**
```bash
# Either drop specific table
mysql -u root -p -e "DROP TABLE greenlight_db.products;"

# Or reset all
./scripts/reset_db.sh
```

### Foreign Key Constraint Fails

**Cause:** Trying to insert Product with non-existent category_id.

**Solution:**
```bash
# Make sure categories exist first
make seed  # This creates categories before products
```

### Cannot Connect to Database

**Docker:**
```bash
# Check if MySQL is ready
docker-compose ps
docker-compose logs db

# Wait for health check
docker-compose up -d
sleep 30  # Wait for MySQL to be ready
```

**Local:**
```bash
# Test connection
mysql -h localhost -u root -p

# Check if database exists
SHOW DATABASES;
USE greenlight_db;
SHOW TABLES;
```

---

## Verification

### Check Tables Created

```bash
# Local
mysql -u root -p -e "USE greenlight_db; SHOW TABLES;"

# Docker
docker-compose exec db mysql -u root -proot -e "USE greenlight_db; SHOW TABLES;"
```

**Expected Output:**
```
+-------------------------+
| Tables_in_greenlight_db |
+-------------------------+
| categories              |
| products                |
| users                   |
+-------------------------+
```

### Check Foreign Keys

```bash
mysql -u root -p greenlight_db -e "
  SELECT 
    TABLE_NAME, 
    CONSTRAINT_NAME, 
    REFERENCED_TABLE_NAME 
  FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
  WHERE TABLE_SCHEMA = 'greenlight_db' 
  AND REFERENCED_TABLE_NAME IS NOT NULL;
"
```

**Expected Output:**
```
+------------+---------------------+-----------------------+
| TABLE_NAME | CONSTRAINT_NAME     | REFERENCED_TABLE_NAME |
+------------+---------------------+-----------------------+
| products   | fk_products_category| categories            |
+------------+---------------------+-----------------------+
```

### Check Data

```bash
# Count records
mysql -u root -p greenlight_db -e "
  SELECT 'users' as table_name, COUNT(*) as count FROM users
  UNION ALL
  SELECT 'categories', COUNT(*) FROM categories
  UNION ALL
  SELECT 'products', COUNT(*) FROM products;
"
```

**After Seeding:**
```
+------------+-------+
| table_name | count |
+------------+-------+
| users      |     1 |
| categories |     2 |
| products   |     6 |
+------------+-------+
```

---

## Quick Commands

```bash
# Create tables
make migrate

# Drop and recreate
./scripts/reset_db.sh

# Check if tables exist
docker-compose exec db mysql -u root -proot -e "USE greenlight_db; SHOW TABLES;"

# View table structure
docker-compose exec db mysql -u root -proot -e "USE greenlight_db; DESCRIBE products;"

# Seed data
make seed
```

---

## Files Reference

```
migrations/
├── 001_create_tables.sql    # Complete schema
├── drop_tables.sql           # Drop all tables
└── migrate.go                # Go migration runner

scripts/
├── reset_db.sh              # Reset script (Linux/Mac)
└── seed.go                  # Data seeder

cmd/api/main.go              # Auto-migration code
```

---

## Best Practices

1. **Always backup** before running migrations in production
2. **Test migrations** in development first
3. **Use transactions** for complex migrations
4. **Version control** migration files
5. **Document** schema changes
6. **Never edit** old migration files, create new ones
7. **Use migration tools** in production (golang-migrate, goose)

---

## Need Help?

Check these files:
- `README.md` - General documentation
- `GETTING_STARTED.md` - Setup guide
- `QUICK_REFERENCE.md` - Quick commands

Or view logs:
```bash
# Docker logs
docker-compose logs -f api

# Check migrations in code
cat cmd/api/main.go | grep -A 15 "migrateDatabase"
```

