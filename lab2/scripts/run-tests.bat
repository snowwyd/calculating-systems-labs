@echo off
chcp 65001 > nul
echo.
echo ИНТЕГРАЦИОННОЕ ТЕСТИРОВАНИЕ МИКРОСЕРВИСОВ LAB2 
echo.

echo Проверка доступности сервисов...
echo.

set ALL_OK=1

curl -s -o nul -w "%%{http_code}" http://localhost:8081/products >nul 2>&1
if errorlevel 1 (
    echo   Product Service ^(порт 8081^) - недоступен
    set ALL_OK=0
) else (
    echo   Product Service ^(порт 8081^) - доступен
)

curl -s -o nul -w "%%{http_code}" http://localhost:8082/orders >nul 2>&1
if errorlevel 1 (
    echo   Order Service ^(порт 8082^) - недоступен
    set ALL_OK=0
) else (
    echo   Order Service ^(порт 8082^) - доступен
)

curl -s -o nul -w "%%{http_code}" http://localhost:8083/inventory >nul 2>&1
if errorlevel 1 (
    echo   Inventory Service ^(порт 8083^) - недоступен
    set ALL_OK=0
) else (
    echo   Inventory Service ^(порт 8083^) - доступен
)

curl -s -o nul -w "%%{http_code}" http://localhost:8084/notifications >nul 2>&1
if errorlevel 1 (
    echo   Notification Service ^(порт 8084^) - недоступен
    set ALL_OK=0
) else (
    echo   Notification Service ^(порт 8084^) - доступен
)

echo.

if %ALL_OK%==0 (
    echo Не все сервисы запущены!
    echo Запустите сервисы командой: start-services.bat
    pause
    exit /b 1
)

echo Все сервисы доступны. Запуск тестов...
echo.

cd integration-test
go run .
cd ..

echo.
pause
