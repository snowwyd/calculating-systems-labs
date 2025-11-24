@echo off
chcp 65001 > nul
echo.
echo УДАЛЕНИЕ РЕСУРСОВ ИЗ KUBERNETES
echo.

echo Удаление Ingress...
kubectl delete -f kubernetes/ingress.yaml

echo.
echo Удаление сервисов...
kubectl delete -f kubernetes/order-service-deployment.yaml
kubectl delete -f kubernetes/product-service-deployment.yaml

echo.
echo Удаление namespace...
kubectl delete -f kubernetes/namespace.yaml

echo.
echo ✓ Все ресурсы удалены!
echo.
pause

