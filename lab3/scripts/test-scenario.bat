@echo off
chcp 65001 > nul
echo.
echo === ТЕСТОВЫЙ СЦЕНАРИЙ LAB3 ===
echo.

echo Проверка доступности сервисов...
echo.

:check_inventory
echo Проверка Inventory Service (8083)...
curl -s -o nul -w "%%{http_code}" http://localhost:8083/inventory >nul 2>&1
if errorlevel 1 (
    echo   Inventory Service недоступен. Подождите...
    timeout /t 2 /nobreak > nul
    goto check_inventory
)
echo   Inventory Service доступен
echo.

:check_order
echo Проверка Order Service (8082)...
curl -s -o nul -w "%%{http_code}" http://localhost:8082/orders >nul 2>&1
if errorlevel 1 (
    echo   Order Service недоступен. Подождите...
    timeout /t 2 /nobreak > nul
    goto check_order
)
echo   Order Service доступен
echo.

echo === Начинаем тестирование ===
echo.

echo 1. Проверка Inventory Service...
curl -X GET http://localhost:8083/inventory
echo.
echo.

echo 2. Создание заказа (синхронная проверка + асинхронное уведомление)...
curl -X POST http://localhost:8082/orders ^
  -H "Content-Type: application/json" ^
  -d "{\"customer_id\":\"CUST001\",\"items\":[{\"product_code\":\"PROD001\",\"product_name\":\"Ноутбук\",\"quantity\":2,\"price\":50000}]}"
echo.
echo.

echo 3. Проверка остатков в Inventory...
curl -X GET http://localhost:8083/inventory
echo.
echo.

echo 4. Просмотр созданных заказов...
curl -X GET http://localhost:8082/orders
echo.
echo.

echo === Проверьте окно Notification Service для уведомлений! ===
echo.
pause

