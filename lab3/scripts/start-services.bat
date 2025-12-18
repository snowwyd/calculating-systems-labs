@echo off
chcp 65001 > nul
echo.
echo ЗАПУСК СЕРВИСОВ LAB3
echo.

echo Запуск RabbitMQ...
docker-compose up -d

echo.
echo Ожидание запуска RabbitMQ (10 секунд)...
timeout /t 10 /nobreak > nul

echo.
echo Запуск Inventory Service на порту 8083...
start "Inventory Service" cmd /k "cd inventory-service && go run ."

timeout /t 2 /nobreak > nul

echo Запуск Order Service на порту 8082...
start "Order Service" cmd /k "cd order-service && go run ."

timeout /t 2 /nobreak > nul

echo Запуск Notification Service (слушатель RabbitMQ)...
start "Notification Service" cmd /k "cd notification-service && go run ."

echo.
echo Все сервисы запущены!
echo.
echo Inventory Service:    http://localhost:8083
echo Order Service:        http://localhost:8082
echo RabbitMQ Management:  http://localhost:15672 (admin/password)
echo.
pause

