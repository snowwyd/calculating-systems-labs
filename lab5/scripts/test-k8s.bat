@echo off
chcp 65001 > nul
echo.
echo ТЕСТИРОВАНИЕ KUBERNETES РАЗВЕРТЫВАНИЯ
echo.

echo Проверка статуса подов...
kubectl get pods -n lab5
echo.

echo Для доступа к сервисам используйте port-forward:
echo.
echo 1. Product Service:
echo    kubectl port-forward svc/product-service 8081:8081 -n lab5
echo    Затем: http://localhost:8081/products
echo.
echo 2. Order Service:
echo    kubectl port-forward svc/order-service 8082:8082 -n lab5
echo    Затем: http://localhost:8082/orders
echo.
echo Открыть автоматически? (откроется port-forward для Product Service)
pause

echo.
echo Запуск port-forward для Product Service...
echo Откройте http://localhost:8081/products в браузере
echo Нажмите Ctrl+C для остановки
echo.
kubectl port-forward svc/product-service 8081:8081 -n lab5

