@echo off
chcp 65001 > nul
echo.
echo === ОТЛАДКА СЕРВИСОВ LAB3 ===
echo.

echo 1. Проверка Inventory Service...
echo.
curl -v http://localhost:8083/inventory
echo.
echo.

echo 2. Проверка Order Service...
echo.
curl -v http://localhost:8082/orders
echo.
echo.

echo 3. Простой тест создания заказа...
echo.
echo Отправка запроса...
curl -v -X POST http://localhost:8082/orders ^
  -H "Content-Type: application/json" ^
  -d "{\"customer_id\":\"CUST001\",\"items\":[{\"product_code\":\"PROD001\",\"product_name\":\"Ноутбук\",\"quantity\":1,\"price\":50000}]}"
echo.
echo.

pause

