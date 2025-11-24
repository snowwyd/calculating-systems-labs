@echo off
chcp 65001 > nul
echo.
echo ТЕСТИРОВАНИЕ DOCKER COMPOSE
echo.

echo Проверка статуса контейнеров...
docker ps --filter "name=lab5" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo.

echo === ТЕСТ 1: Product Service ===
echo GET http://localhost:8081/products
curl -s http://localhost:8081/products
echo.
echo.

echo === ТЕСТ 2: Order Service ===
echo GET http://localhost:8082/orders
curl -s http://localhost:8082/orders
echo.
echo.

echo === ТЕСТ 3: Создание заказа ===
echo POST http://localhost:8082/orders
curl -X POST http://localhost:8082/orders ^
  -H "Content-Type: application/json" ^
  -d "{\"customer_id\":\"CUST001\",\"product_id\":\"1\",\"quantity\":2,\"total_price\":100000}"
echo.
echo.

echo === ТЕСТ 4: Проверка созданного заказа ===
echo GET http://localhost:8082/orders
curl -s http://localhost:8082/orders
echo.
echo.

echo V Все тесты завершены!
echo.
pause

