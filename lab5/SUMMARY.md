# Lab 5: Итоговый отчет

## ✅ Задание выполнено полностью

**Задача:** Создать контейнеры системы с помощью Docker. Развернуть архитектурное решение с помощью Docker Compose и Kubernetes.

## 📦 Что реализовано:

### 1. Docker Контейнеры

#### Product Service
- **Технология:** Multi-stage Docker build
- **Базовый образ:** golang:1.21-alpine → alpine:latest
- **Размер образа:** ~26 MB
- **Порт:** 8081
- **Функционал:** CRUD для продуктов (in-memory)

#### Order Service
- **Технология:** Multi-stage Docker build
- **Базовый образ:** golang:1.21-alpine → alpine:latest
- **Размер образа:** ~26 MB
- **Порт:** 8082
- **Функционал:** CRUD для заказов (in-memory)

### 2. Docker Compose

**Файл:** `docker-compose.yml`

**Особенности:**
- ✅ Автоматическая сборка образов
- ✅ Health checks для мониторинга
- ✅ Автоматический перезапуск (restart: unless-stopped)
- ✅ Изолированная сеть (lab5-network)
- ✅ Маппинг портов на хост

**Запуск:**
```bash
build-images.bat      # Сборка образов
start-docker.bat      # Запуск сервисов
test-docker.bat       # Тестирование
stop-docker.bat       # Остановка
```

**Результат тестирования:**
```
✓ Product Service: http://localhost:8081/products
✓ Order Service:   http://localhost:8082/orders
✓ Создание заказов работает
✓ Health checks активны
```

### 3. Kubernetes

**Файлы манифестов:**
- `kubernetes/namespace.yaml` - Namespace для изоляции
- `kubernetes/product-service-deployment.yaml` - Deployment + Service
- `kubernetes/order-service-deployment.yaml` - Deployment + Service
- `kubernetes/ingress.yaml` - Ingress для маршрутизации

**Особенности:**
- ✅ 2 реплики каждого сервиса (высокая доступность)
- ✅ Liveness и Readiness проbes
- ✅ Resource limits (CPU/Memory)
- ✅ ClusterIP Services для внутренней балансировки
- ✅ Ingress для единой точки входа
- ✅ Автоматическое восстановление подов

**Запуск:**
```bash
check-k8s.bat         # Проверка кластера
build-images.bat      # Сборка образов
deploy-k8s.bat        # Развертывание
test-k8s-simple.bat   # Проверка статуса
delete-k8s.bat        # Удаление
```

**Результат развертывания:**
```
✓ Namespace lab5 создан
✓ 4 пода запущены (2 Product + 2 Order)
✓ 2 сервиса работают
✓ Deployments готовы (2/2)
✓ ReplicaSets активны
✓ Ingress настроен
```

## 📊 Сравнение решений

| Характеристика | Docker Compose | Kubernetes |
|---|---|---|
| **Сложность настройки** | ⭐ Простая | ⭐⭐⭐ Средняя |
| **Скорость запуска** | ~5 сек | ~15 сек |
| **Масштабируемость** | Ручная | Автоматическая |
| **Высокая доступность** | Нет | Да (2+ реплики) |
| **Балансировка нагрузки** | Нет | Да (автоматическая) |
| **Self-healing** | Перезапуск | Да (автовосстановление) |
| **Мониторинг** | Health checks | Probes + Metrics |
| **Подходит для** | Dev/Test | Prod |

## 🎯 Архитектурное решение

### Docker Compose (Development/Testing)

```
┌────────────────────────────────┐
│     Docker Compose             │
├────────────────────────────────┤
│  Product Service :8081         │
│  Order Service   :8082         │
└────────────────────────────────┘
          ↓
  lab5-network (bridge)
```

**Плюсы:**
- Простота настройки и запуска
- Идеально для локальной разработки
- Быстрое тестирование изменений

**Минусы:**
- Один хост
- Ручное масштабирование
- Нет автоматической балансировки

### Kubernetes (Production)

```
┌─────────────────────────────────────────┐
│         Kubernetes Cluster              │
├─────────────────────────────────────────┤
│  Namespace: lab5                        │
│                                         │
│  ┌─ Product Service Deployment ───┐    │
│  │  Pod 1 (Replica 1)              │    │
│  │  Pod 2 (Replica 2)              │    │
│  └─────────────────────────────────┘    │
│           ↓ Service (ClusterIP)         │
│                                         │
│  ┌─ Order Service Deployment ──────┐   │
│  │  Pod 3 (Replica 1)              │    │
│  │  Pod 4 (Replica 2)              │    │
│  └─────────────────────────────────┘    │
│           ↓ Service (ClusterIP)         │
│                                         │
│  Ingress (lab5.local)                   │
└─────────────────────────────────────────┘
```

**Плюсы:**
- Автоматическая балансировка нагрузки
- Self-healing (автовосстановление)
- Горизонтальное масштабирование
- Rolling updates без простоя
- Production-ready

**Минусы:**
- Требует кластер Kubernetes
- Более сложная настройка
- Требует знаний K8s

## 🚀 Инструкции по использованию

### Вариант 1: Docker Compose (Рекомендуется для начала)

```bash
# 1. Собрать образы
.\build-images.bat

# 2. Запустить
.\start-docker.bat

# 3. Тестировать
.\test-docker.bat

# 4. Остановить
.\stop-docker.bat
```

**Доступ:**
- Product Service: http://localhost:8081/products
- Order Service: http://localhost:8082/orders

### Вариант 2: Kubernetes

```bash
# 0. Проверить кластер
.\check-k8s.bat

# 1. Собрать образы
.\build-images.bat

# 2. Развернуть
.\deploy-k8s.bat

# 3. Проверить статус
.\test-k8s-simple.bat

# 4. Доступ через port-forward
kubectl port-forward svc/product-service 8081:8081 -n lab5
kubectl port-forward svc/order-service 8082:8082 -n lab5

# 5. Удалить
.\delete-k8s.bat
```

## 📁 Структура проекта

```
lab5/
├── docker-compose.yml              # Docker Compose конфигурация
├── build-images.bat                # Сборка Docker образов
├── start-docker.bat                # Запуск Docker Compose
├── stop-docker.bat                 # Остановка Docker Compose
├── test-docker.bat                 # Тестирование Docker Compose
├── check-k8s.bat                   # Проверка Kubernetes
├── deploy-k8s.bat                  # Развертывание в K8s
├── delete-k8s.bat                  # Удаление из K8s
├── test-k8s.bat                    # Тестирование K8s (port-forward)
├── test-k8s-simple.bat             # Проверка статуса K8s
├── README.md                       # Основная документация
├── KUBERNETES-SETUP.md             # Настройка Kubernetes
├── SUMMARY.md                      # Этот файл
├── .dockerignore                   # Игнорируемые файлы для Docker
├── .gitignore                      # Игнорируемые файлы для Git
│
├── product-service/
│   ├── Dockerfile                  # Multi-stage build
│   ├── main.go                     # REST API сервис
│   ├── go.mod
│   └── go.sum
│
├── order-service/
│   ├── Dockerfile                  # Multi-stage build
│   ├── main.go                     # REST API сервис
│   ├── go.mod
│   └── go.sum
│
└── kubernetes/
    ├── namespace.yaml              # Namespace lab5
    ├── product-service-deployment.yaml
    ├── order-service-deployment.yaml
    └── ingress.yaml
```

## 🎓 Выводы

### Достигнутые цели:

1. ✅ **Контейнеризация** - оба сервиса упакованы в Docker контейнеры
2. ✅ **Docker Compose** - оркестрация для локальной разработки
3. ✅ **Kubernetes** - production-ready развертывание
4. ✅ **Автоматизация** - скрипты для всех операций
5. ✅ **Документация** - полное описание решения

### Полученные навыки:

- 📦 Multi-stage Docker builds
- 🐳 Docker Compose оркестрация
- ☸️ Kubernetes Deployments, Services, Ingress
- 🔄 Health checks и Probes
- 📊 Балансировка нагрузки
- 🛠️ Автоматизация развертывания

### Best Practices применены:

- ✅ Multi-stage builds для минимизации размера образов
- ✅ Health checks для мониторинга
- ✅ Resource limits для стабильности
- ✅ Namespace для изоляции ресурсов
- ✅ Multiple replicas для высокой доступности
- ✅ Liveness/Readiness probes для надежности
- ✅ .dockerignore для оптимизации сборки

## 📚 Дополнительные материалы

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Multi-stage builds](https://docs.docker.com/build/building/multi-stage/)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)

---

**Проект выполнен:** 20.11.2025  
**Статус:** ✅ Полностью работает  
**Тестирование:** ✅ Пройдено

