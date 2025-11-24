# Lab 5: Docker и Kubernetes

Развертывание микросервисной архитектуры с помощью Docker Compose и Kubernetes.

## Быстрый старт (Docker Compose - рекомендуется)

```cmd
# 1. Собрать образы
scripts\build-images.bat

# 2. Запустить
scripts\start-docker.bat

# Готово! Сервисы доступны:
# http://localhost:8081/products
# http://localhost:8082/orders
```

**Для Kubernetes:** см. раздел "Часть 2" и файл `KUBERNETES-SETUP.md`

## Архитектура

```
Docker / Kubernetes         

Product Service (8081) - 2 реплики   
Order Service (8082) - 2 реплики   

```

## Часть 1: Docker Compose

### Структура

- **Dockerfile** для каждого сервиса (multi-stage build)
- **docker-compose.yml** для оркестрации
- Автоматический перезапуск
- Health checks
- Изолированная сеть

### Быстрый старт

#### 1. Сборка образов

```cmd
scripts\build-images.bat
```

Это создаст Docker образы:
- `lab5/product-service:latest`
- `lab5/order-service:latest`

#### 2. Запуск через Docker Compose

```cmd
scripts\start-docker.bat
```

Сервисы будут доступны на:
- Product Service: http://localhost:8081
- Order Service: http://localhost:8082

#### 3. Остановка

```cmd
scripts\stop-docker.bat
```

### Ручное управление Docker Compose

```bash
# Запуск
docker-compose up -d

# Логи
docker-compose logs -f

# Статус
docker-compose ps

# Остановка
docker-compose down

# Пересборка и запуск
docker-compose up -d --build
```

## Часть 2: Kubernetes

### Требования для Kubernetes

**ВАЖНО:** Kubernetes требует запущенного кластера!

#### Вариант 1: Docker Desktop (рекомендуется для Windows)

1. Откройте Docker Desktop
2. Settings → Kubernetes
3. Включите "Enable Kubernetes"
4. Нажмите "Apply & Restart"
5. Дождитесь запуска (зелёная иконка Kubernetes внизу)

#### Вариант 2: Minikube

```bash
# Установка (в PowerShell от администратора)
choco install minikube

# Запуск кластера
minikube start

# Использование Docker демона Minikube для образов
minikube docker-env | Invoke-Expression
```

#### Проверка готовности

```cmd
scripts\check-k8s.bat
```

Этот скрипт проверит:
- Установлен ли kubectl
- Запущен ли Kubernetes кластер
- Доступна ли API Kubernetes

### Структура манифестов

```
kubernetes/
├── namespace.yaml                    # Namespace lab5
├── product-service-deployment.yaml   # Deployment + Service
├── order-service-deployment.yaml     # Deployment + Service
└── ingress.yaml                      # Ingress для маршрутизации
```

### Развертывание

#### 0. Проверка готовности (ОБЯЗАТЕЛЬНО!)

```cmd
scripts\check-k8s.bat
```

Если кластер не запущен, следуйте инструкциям выше.

#### 1. Подготовка образов для Kubernetes

Если используете Minikube:

```bash
# Используем Docker демон Minikube
eval $(minikube docker-env)

# Собираем образы
docker build -t lab5/product-service:latest ./product-service
docker build -t lab5/order-service:latest ./order-service
```

Или используйте скрипт:

```cmd
scripts\build-images.bat
```

#### 2. Развертывание в Kubernetes

```cmd
scripts\deploy-k8s.bat
```

Это создаст:
- Namespace `lab5`
- 2 реплики Product Service
- 2 реплики Order Service
- Services для каждого
- Ingress для маршрутизации

#### 3. Проверка статуса

```bash
# Все ресурсы
kubectl get all -n lab5

# Поды
kubectl get pods -n lab5

# Сервисы
kubectl get svc -n lab5

# Ingress
kubectl get ingress -n lab5
```

#### 4. Логи

```bash
# Product Service
kubectl logs -f deployment/product-service -n lab5

# Order Service
kubectl logs -f deployment/order-service -n lab5
```

#### 5. Доступ к сервисам

**Через port-forward:**

```bash
# Product Service
kubectl port-forward svc/product-service 8081:8081 -n lab5

# Order Service
kubectl port-forward svc/order-service 8082:8082 -n lab5
```

Теперь доступны на:
- http://localhost:8081/products
- http://localhost:8082/orders

**Через Ingress (если настроен):**

Добавьте в `C:\Windows\System32\drivers\etc\hosts`:
```
127.0.0.1  lab5.local
```

Доступ:
- http://lab5.local/products
- http://lab5.local/orders

#### 6. Масштабирование

```bash
# Увеличить количество реплик
kubectl scale deployment product-service --replicas=3 -n lab5
kubectl scale deployment order-service --replicas=3 -n lab5

# Проверить
kubectl get pods -n lab5
```

#### 7. Удаление

```cmd
scripts\delete-k8s.bat
```

Или вручную:
```bash
kubectl delete namespace lab5
```

## Тестирование

### Docker Compose

```bash
# После scripts\start-docker.bat

curl http://localhost:8081/products
curl http://localhost:8082/orders

# Создать заказ
curl -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id":"CUST001","product_id":"1","quantity":2,"total_price":100000}'
```

### Kubernetes

```bash
# После port-forward

curl http://localhost:8081/products
curl http://localhost:8082/orders
```
