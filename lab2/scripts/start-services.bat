@echo off
chcp 65001 > nul
echo.
echo ================================================================
echo       ЗАПУСК СЕРВИСОВ LAB2
echo ================================================================
echo.

echo Запуск баз данных...
docker-compose up -d

echo.
echo Ожидание запуска баз данных (5 секунд)...
timeout /t 5 /nobreak > nul

echo.
echo Запуск Product Service на порту 8081...
start "Product Service" cmd /k "cd product-service && go run ."

timeout /t 2 /nobreak > nul

echo Запуск Order Service на порту 8082...
start "Order Service" cmd /k "cd order-service && go run ."

timeout /t 2 /nobreak > nul

echo Запуск Inventory Service на порту 8083...
start "Inventory Service" cmd /k "cd inventory-service && go run ."

timeout /t 2 /nobreak > nul

echo Запуск Notification Service на порту 8084...
start "Notification Service" cmd /k "cd notification-service && go run ."

echo.
echo Все сервисы запущены!
echo.
echo Product Service:      http://localhost:8081
echo Order Service:        http://localhost:8082
echo Inventory Service:    http://localhost:8083
echo Notification Service: http://localhost:8084
echo.
pause

