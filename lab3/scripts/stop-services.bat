@echo off
chcp 65001 > nul
echo.
echo Остановка сервисов...
echo.

REM Остановка процессов на портах
for %%p in (8082 8083) do (
    for /f "tokens=5" %%a in ('netstat -aon ^| find ":%%p" ^| find "LISTENING"') do (
        taskkill /F /PID %%a 2>nul
        if not errorlevel 1 echo Порт %%p освобожден
    )
)

echo.
echo Остановка RabbitMQ...
docker-compose down

echo.
echo ✓ Все сервисы остановлены!
echo.
pause

