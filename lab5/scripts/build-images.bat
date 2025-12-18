@echo off
chcp 65001 > nul
echo.
echo СБОРКА DOCKER ОБРАЗОВ LAB5
echo.

echo Сборка Product Service...
docker build -t lab5/product-service:latest ./product-service
if errorlevel 1 (
    echo ОШИБКА: Не удалось собрать Product Service
    pause
    exit /b 1
)
echo Product Service собран
echo.

echo Сборка Order Service...
docker build -t lab5/order-service:latest ./order-service
if errorlevel 1 (
    echo ОШИБКА: Не удалось собрать Order Service
    pause
    exit /b 1
)
echo Order Service собран
echo.

echo Все образы успешно собраны!
echo.
docker images | findstr "lab5"
echo.
pause

