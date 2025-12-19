#!/bin/bash
# Docker initialization script
# This runs SQL migration before starting the API

set -e

echo "🔄 Waiting for MySQL to be ready..."
max_attempts=30
attempt=0
while ! nc -z db 3306; do
  attempt=$((attempt + 1))
  if [ $attempt -gt $max_attempts ]; then
    echo "❌ MySQL not ready after $max_attempts attempts"
    exit 1
  fi
  echo "   Attempt $attempt/$max_attempts..."
  sleep 2
done
echo "✅ MySQL is ready!"

echo ""
echo "🗄️  Running SQL migrations..."
if mysql --skip-ssl -h db -u root -proot greenlight_db < /root/migrations/001_create_tables.sql 2>&1; then
  echo "✅ Migration completed successfully"
else
  echo "⚠️  Migration may have already run or failed"
fi

echo ""
echo "🔍 Verifying tables..."
TABLE_COUNT=$(mysql --skip-ssl -h db -u root -proot greenlight_db -se "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='greenlight_db' AND table_name IN ('users','categories','products')")
echo "   Found $TABLE_COUNT tables"

if [ "$TABLE_COUNT" -lt 3 ]; then
  echo "❌ Tables not created! Expected 3, found $TABLE_COUNT"
  echo "   Listing tables:"
  mysql --skip-ssl -h db -u root -proot greenlight_db -e "SHOW TABLES"
  exit 1
fi

echo "✅ All tables verified"
echo ""
echo "🚀 Starting API..."
exec /root/main

