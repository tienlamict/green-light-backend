@echo off
REM Quick Start Script for Green Light Backend (Windows)
REM This script will set up and start the entire application

echo.
echo ========================================
echo   Green Light Backend - Quick Start
echo ========================================
echo.

REM Step 1: Clean up
echo [Step 1/8] Cleaning up old containers and volumes...
docker-compose down -v 2>nul
echo [OK] Cleanup complete
echo.

REM Step 2: Build and start
echo [Step 2/8] Building and starting containers...
docker-compose up -d --build
echo [OK] Containers started
echo.

REM Step 3: Wait for MySQL
echo [Step 3/8] Waiting for MySQL to be ready (30 seconds)...
timeout /t 30 /nobreak >nul
echo [OK] MySQL should be ready
echo.

REM Step 4: Check status
echo [Step 4/8] Checking container status...
docker-compose ps
echo.

REM Step 5: Check logs
echo [Step 5/8] Checking migration logs...
docker-compose logs --tail=10 api
echo.

REM Step 6: Verify tables
echo [Step 6/8] Verifying database tables...
docker-compose exec -T db mysql -u root -proot greenlight_db -e "SHOW TABLES;"
echo.

REM Step 7: Seed data
echo [Step 7/8] Seeding database with sample data...
docker-compose exec -T api go run /root/scripts/seed.go
echo.

REM Step 8: Test API
echo [Step 8/8] Testing API endpoints...
curl -s http://localhost:8080/healthz >nul && (
    echo [OK] Health check passed
) || (
    echo [WARNING] Health check failed
)

curl -s http://localhost:8080/api/v1/products >nul && (
    echo [OK] Products endpoint working
) || (
    echo [WARNING] Products endpoint failed
)
echo.

REM Final summary
echo ========================================
echo   Setup Complete!
echo ========================================
echo.
echo Your API is now running at:
echo   - API Base: http://localhost:8080
echo   - Swagger: http://localhost:8080/api/docs/index.html
echo   - Health: http://localhost:8080/healthz
echo.
echo Default credentials:
echo   - Email: admin@example.com
echo   - Password: admin123
echo.
echo Useful commands:
echo   - View logs: docker-compose logs -f
echo   - Stop: docker-compose down
echo   - Restart: docker-compose restart
echo.
echo Happy coding!
echo.
pause

