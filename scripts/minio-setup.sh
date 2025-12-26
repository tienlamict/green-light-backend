#!/bin/sh
set -e

echo "=========================================="
echo "MinIO Setup Script"
echo "=========================================="

# Wait for MinIO to be ready
echo "Waiting for MinIO to be ready..."
MAX_RETRIES=30
RETRY=0
until mc alias set myminio http://minio:9000 ${MINIO_ROOT_USER} ${MINIO_ROOT_PASSWORD} 2>/dev/null; do
  RETRY=$((RETRY+1))
  if [ $RETRY -ge $MAX_RETRIES ]; then
    echo "ERROR: MinIO is not ready after $MAX_RETRIES retries"
    exit 1
  fi
  echo "MinIO is not ready yet. Waiting... ($RETRY/$MAX_RETRIES)"
  sleep 2
done

echo "✓ MinIO is ready!"

# Create bucket if it doesn't exist
echo "Creating bucket '${MINIO_BUCKET}' if it doesn't exist..."
mc mb myminio/${MINIO_BUCKET} 2>/dev/null && echo "✓ Bucket created" || echo "ℹ Bucket already exists"

# Configure CORS for MinIO server (global setting)
echo "Configuring CORS for MinIO server..."
# Set CORS allow origin globally for MinIO server
# This affects all buckets
if mc admin config set myminio/ api cors_allow_origin="*" 2>/dev/null; then
  echo "✓ CORS allow origin set to '*'"
  echo "Restarting MinIO service to apply CORS configuration..."
  mc admin service restart myminio 2>/dev/null || echo "⚠ Warning: Could not restart MinIO. Please restart manually."
  echo "Waiting 5 seconds for MinIO to restart..."
  sleep 5
  # Re-establish connection after restart
  mc alias set myminio http://minio:9000 ${MINIO_ROOT_USER} ${MINIO_ROOT_PASSWORD} 2>/dev/null || echo "⚠ Warning: Could not reconnect to MinIO"
else
  echo "⚠ Note: Could not set CORS via admin config"
  echo "   You may need to configure CORS manually via MinIO Console"
fi

# Set bucket anonymous policy for public read access
echo "Setting bucket anonymous policy..."
mc anonymous set download myminio/${MINIO_BUCKET} 2>/dev/null && echo "✓ Anonymous download policy set" || echo "ℹ Anonymous policy already set or not needed"

echo ""
echo "✓ MinIO setup completed!"
echo ""
echo "To verify CORS configuration:"
echo "1. Check MinIO Console: http://localhost:9001"
echo "2. Go to Buckets > ${MINIO_BUCKET} > Management > CORS Configuration"
echo "3. Or run: mc admin config get myminio/ api | grep cors"

echo "=========================================="
echo "✓ MinIO setup completed successfully!"
echo "=========================================="

