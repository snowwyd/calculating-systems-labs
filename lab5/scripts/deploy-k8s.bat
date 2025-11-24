@echo off
chcp 65001 > nul
echo.
echo РАЗВЕРТЫВАНИЕ В KUBERNETES
echo.

echo Проверка доступности Kubernetes...
kubectl cluster-info >nul 2>&1
if errorlevel 1 (
    echo.
    echo ✗ ОШИБКА: Kubernetes кластер недоступен!
    echo.
    echo Запустите check-k8s.bat для диагностики
    echo.
    pause
    exit /b 1
)
echo ✓ Kubernetes кластер доступен
echo.

echo Создание namespace...
kubectl apply -f kubernetes/namespace.yaml

echo.
echo Развертывание Product Service...
kubectl apply -f kubernetes/product-service-deployment.yaml

echo.
echo Развертывание Order Service...
kubectl apply -f kubernetes/order-service-deployment.yaml

echo.
echo Создание Ingress...
kubectl apply -f kubernetes/ingress.yaml

echo.
echo ✓ Развертывание завершено!
echo.
echo Проверка статуса:
kubectl get all -n lab5

echo.
echo Получение подов:
kubectl get pods -n lab5

echo.
echo Логи сервиса:
echo   kubectl logs -f deployment/product-service -n lab5
echo   kubectl logs -f deployment/order-service -n lab5
echo.
pause

