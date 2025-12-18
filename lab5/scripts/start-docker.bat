@echo off
chcp 65001 > nul
echo.
echo ЗАПУСК СЕРВИСОВ LAB5 ЧЕРЕЗ DOCKER COMPOSE
echo.

echo Остановка старых контейнеров...
docker-compose down

echo.
echo Запуск контейнеров...
docker-compose up -d

echo.
echo Ожидание готовности сервисов (10 секунд)...
timeout /t 10 /nobreak > nul

echo.
echo Проверка статуса контейнеров...
docker-compose ps

echo.
echo Сервисы запущены!
echo.
echo Product Service: http://localhost:8081
echo Order Service:   http://localhost:8082
echo.
echo Просмотр логов:
echo   docker-compose logs -f product-service
echo   docker-compose logs -f order-service
echo.
pause

