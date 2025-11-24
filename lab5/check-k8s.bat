@echo off
chcp 65001 > nul
echo.
echo ПРОВЕРКА KUBERNETES
echo.

echo Проверка kubectl...
kubectl version --client >nul 2>&1
if errorlevel 1 (
    echo X kubectl не установлен или недоступен
    echo.
    echo Установите kubectl: https://kubernetes.io/docs/tasks/tools/
    pause
    exit /b 1
)
echo V kubectl установлен

echo.
echo Проверка подключения к кластеру...
kubectl cluster-info >nul 2>&1
if errorlevel 1 (
    echo X Kubernetes кластер НЕ запущен или недоступен
    echo.
    echo === ВАРИАНТЫ РЕШЕНИЯ ===
    echo.
    echo 1. Docker Desktop:
    echo    - Откройте Docker Desktop
    echo    - Settings - Kubernetes
    echo    - Включите Enable Kubernetes
    echo    - Нажмите Apply and Restart
    echo.
    echo 2. Minikube:
    echo    - Установите: https://minikube.sigs.k8s.io/docs/start/
    echo    - Запустите: minikube start
    echo.
    echo 3. Используйте Docker Compose вместо Kubernetes
    echo    - start-docker.bat
    echo.
    pause
    exit /b 1
)

echo V Kubernetes кластер доступен
echo.
echo Информация о кластере:
kubectl cluster-info
echo.
echo Версия kubectl:
kubectl version --client
echo.
echo V Всё готово для развертывания!
echo.
pause
