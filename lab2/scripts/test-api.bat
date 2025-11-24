@echo off
chcp 65001 > nul
echo.
echo === Тестирование Product Service (8081) ===
echo.

echo Создание продукта...
curl -X POST http://localhost:8081/products ^
  -H "Content-Type: application/json" ^
  -d "{\"code\":\"P001\",\"name\":\"Ноутбук\",\"weight\":2.5,\"description\":\"Игровой ноутбук\"}"
echo.
echo.

echo Получение всех продуктов...
curl -X GET http://localhost:8081/products
echo.
echo.

echo === Тестирование Order Service (8082) ===
echo.

echo Создание заказа...
curl -X POST http://localhost:8082/orders ^
  -H "Content-Type: application/json" ^
  -d "{\"order_code\":\"O001\",\"order_date\":\"2024-01-01\",\"product_items\":[{\"product_code\":\"P001\",\"name\":\"Ноутбук\",\"quantity\":2,\"cost\":1500.00}]}"
echo.
echo.

echo Получение всех заказов...
curl -X GET http://localhost:8082/orders
echo.
echo.

echo === Тестирование Inventory Service (8083) ===
echo.

echo Добавление товара на склад...
curl -X POST http://localhost:8083/inventory ^
  -H "Content-Type: application/json" ^
  -d "{\"product_code\":\"P001\",\"name\":\"Ноутбук\",\"quantity\":50,\"current_price\":\"1500.00\"}"
echo.
echo.

echo Получение данных склада...
curl -X GET http://localhost:8083/inventory
echo.
echo.

echo === Тестирование Notification Service (8084) ===
echo.

echo Создание уведомления...
curl -X POST http://localhost:8084/notifications ^
  -H "Content-Type: application/json" ^
  -d "{\"message_type\":\"INFO\",\"description\":\"Заказ создан\",\"date\":\"%date% %time%\"}"
echo.
echo.

echo Получение всех уведомлений...
curl -X GET http://localhost:8084/notifications
echo.
echo.

echo === Тестирование завершено! ===
echo.
pause

