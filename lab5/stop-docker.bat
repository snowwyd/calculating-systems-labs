@echo off
chcp 65001 > nul
echo.
echo ОСТАНОВКА DOCKER COMPOSE
echo.

docker-compose down

echo.
echo ✓ Все контейнеры остановлены и удалены!
echo.
pause

