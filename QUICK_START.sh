#!/bin/bash

# Quick Start Script for Green Light Backend
# This script will set up and start the entire application

set -e

echo "🚀 Green Light Backend - Quick Start"
echo "===================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Clean up
echo "📦 Step 1: Cleaning up old containers and volumes..."
docker-compose down -v 2>/dev/null || true
echo -e "${GREEN}✓ Cleanup complete${NC}"
echo ""

# Step 2: Build and start
echo "🏗️  Step 2: Building and starting containers..."
docker-compose up -d --build
echo -e "${GREEN}✓ Containers started${NC}"
echo ""

# Step 3: Wait for MySQL
echo "⏳ Step 3: Waiting for MySQL to be ready (30 seconds)..."
sleep 30
echo -e "${GREEN}✓ MySQL should be ready${NC}"
echo ""

# Step 4: Check status
echo "🔍 Step 4: Checking container status..."
docker-compose ps
echo ""

# Step 5: Check logs for migration
echo "📋 Step 5: Checking migration logs..."
docker-compose logs api | tail -10
echo ""

# Step 6: Verify tables
echo "🗄️  Step 6: Verifying database tables..."
docker-compose exec -T db mysql -u root -proot greenlight_db -e "SHOW TABLES;" || echo -e "${YELLOW}⚠ Could not verify tables${NC}"
echo ""

# Step 7: Seed data
echo "🌱 Step 7: Seeding database with sample data..."
docker-compose exec -T api go run /root/scripts/seed.go || echo -e "${YELLOW}⚠ Seeding may have failed - check logs${NC}"
echo ""

# Step 8: Test API
echo "🧪 Step 8: Testing API endpoints..."
echo -n "  Health check: "
curl -s http://localhost:8080/healthz > /dev/null && echo -e "${GREEN}✓ OK${NC}" || echo -e "${YELLOW}✗ Failed${NC}"

echo -n "  Products endpoint: "
curl -s http://localhost:8080/api/v1/products > /dev/null && echo -e "${GREEN}✓ OK${NC}" || echo -e "${YELLOW}✗ Failed${NC}"
echo ""

# Final summary
echo "=================================="
echo "🎉 Setup Complete!"
echo "=================================="
echo ""
echo "Your API is now running at:"
echo "  📍 API Base: http://localhost:8080"
echo "  💚 Health: http://localhost:8080/healthz"
echo ""
echo "Default credentials:"
echo "  📧 Email: admin@example.com"
echo "  🔑 Password: admin123"
echo ""
echo "Useful commands:"
echo "  📊 View logs: docker-compose logs -f"
echo "  🛑 Stop: docker-compose down"
echo "  🔄 Restart: docker-compose restart"
echo ""
echo "Happy coding! 🚀"

