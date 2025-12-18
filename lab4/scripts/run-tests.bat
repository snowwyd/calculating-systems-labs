@echo off
chcp 65001 > nul
echo.
echo ИНТЕГРАЦИОННОЕ ТЕСТИРОВАНИЕ LAB4 - GraphQL Gateway
echo.

echo Проверка доступности сервисов...
echo.

set ALL_OK=1

curl -s -o nul http://localhost:8081/products >nul 2>&1
if errorlevel 1 (
    echo   Product Service ^(8081^) - недоступен
    set ALL_OK=0
) else (
    echo   Product Service ^(8081^) - доступен
)

curl -s -o nul http://localhost:8082/orders >nul 2>&1
if errorlevel 1 (
    echo   Order Service ^(8082^) - недоступен
    set ALL_OK=0
) else (
    echo   Order Service ^(8082^) - доступен
)

curl -s -o nul http://localhost:8080/graphql >nul 2>&1
if errorlevel 1 (
    echo   GraphQL Gateway ^(8080^) - недоступен
    set ALL_OK=0
) else (
    echo   GraphQL Gateway ^(8080^) - доступен
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

