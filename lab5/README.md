# Lab5

## Быстрый старт (Docker Compose)

```cmd
# 1. Собрать образы
scripts\build-images.bat

# 2. Запустить
scripts\start-docker.bat

# Готово! Сервисы доступны:
# http://localhost:8081/products
# http://localhost:8082/orders
```

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

## Часть 2: Kubernetes

### Требования для Kubernetes


#### Docker Desktop (рекомендуется для Windows)

1. Откройте Docker Desktop
2. Settings → Kubernetes
3. Включите "Enable Kubernetes"
4. Нажмите "Apply & Restart"
5. Дождитесь запуска (зелёная иконка Kubernetes внизу)

#### Проверка готовности

```cmd
scripts\check-k8s.bat
```

### Развертывание

#### 0. Проверка готовности

```cmd
scripts\check-k8s.bat
```

#### 1. Подготовка образов для Kubernetes

```cmd
scripts\build-images.bat
```

#### 2. Развертывание в Kubernetes

```cmd
scripts\deploy-k8s.bat
```

#### 3. Масштабирование

```bash
# Увеличить количество реплик
kubectl scale deployment product-service --replicas=3 -n lab5
kubectl scale deployment order-service --replicas=3 -n lab5

# Проверить
kubectl get pods -n lab5
```

#### 4. Удаление

```cmd
scripts\delete-k8s.bat
```