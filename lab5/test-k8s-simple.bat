@echo off
chcp 65001 > nul
echo.
echo БЫСТРЫЙ ТЕСТ KUBERNETES
echo.

echo Статус подов:
kubectl get pods -n lab5
echo.

echo Статус сервисов:
kubectl get svc -n lab5
echo.

echo Логи Product Service:
kubectl logs deployment/product-service -n lab5 --tail=5
echo.

echo Логи Order Service:
kubectl logs deployment/order-service -n lab5 --tail=5
echo.

echo === ДЛЯ ДОСТУПА К API ===
echo.
echo Запустите в ОТДЕЛЬНЫХ окнах PowerShell:
echo.
echo Product Service:
echo   kubectl port-forward svc/product-service 8081:8081 -n lab5
echo   curl http://localhost:8081/products
echo.
echo Order Service:
echo   kubectl port-forward svc/order-service 8082:8082 -n lab5
echo   curl http://localhost:8082/orders
echo.
pause

