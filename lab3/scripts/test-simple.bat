@echo off
chcp 65001 > nul
echo.
echo === ПРОСТОЙ ТЕСТ LAB3 ===
echo.

echo Шаг 1: Проверка Inventory Service...
curl http://localhost:8083/inventory
if errorlevel 1 (
    echo.
    echo ОШИБКА: Inventory Service недоступен!
    echo Запустите: start-services.bat
    pause
    exit /b 1
)
echo.
echo Inventory Service работает
echo.
timeout /t 2 /nobreak > nul

echo Шаг 2: Создание простого заказа (1 товар)...
curl -X POST http://localhost:8082/orders -H "Content-Type: application/json" -d "{\"customer_id\":\"TEST\",\"items\":[{\"product_code\":\"PROD002\",\"product_name\":\"Мышь\",\"quantity\":1,\"price\":1000}]}"
if errorlevel 1 (
    echo.
    echo ОШИБКА: Не удалось создать заказ!
    pause
    exit /b 1
)
echo.
echo Заказ создан
echo.
timeout /t 2 /nobreak > nul

echo Шаг 3: Проверка созданного заказа...
curl http://localhost:8082/orders
echo.
echo.

echo Шаг 4: Проверка остатков на складе...
curl http://localhost:8083/inventory
echo.
echo.

echo ТЕСТ ЗАВЕРШЕН УСПЕШНО!
echo.
echo Проверьте окно Notification Service - там должно быть уведомление!
echo.
pause

