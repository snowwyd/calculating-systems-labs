@echo off
chcp 65001 > nul
echo.
echo ЗАПУСК СЕРВИСОВ LAB4
echo.

echo Запуск Product Service на порту 8081...
start "Product Service" cmd /k "cd product-service && go run ."

timeout /t 2 /nobreak > nul

echo Запуск Order Service на порту 8082...
start "Order Service" cmd /k "cd order-service && go run ."

timeout /t 2 /nobreak > nul

echo Запуск GraphQL Gateway на порту 8080...
start "GraphQL Gateway" cmd /k "cd graphql-gateway && go run ."

echo.
echo Все сервисы запущены!
echo.
echo Product Service:  http://localhost:8081
echo Order Service:    http://localhost:8082
echo GraphQL Gateway:  http://localhost:8080
echo GraphiQL UI:      http://localhost:8080
echo.
pause

