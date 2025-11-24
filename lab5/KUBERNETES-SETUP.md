# Настройка Kubernetes для Lab5

## Что такое Kubernetes?

Kubernetes (K8s) - это система оркестрации контейнеров. Для работы нужен **работающий кластер**.

## Варианты настройки

### ✅ Вариант 1: Docker Desktop (самый простой для Windows)

#### Шаги:

1. **Убедитесь, что Docker Desktop установлен и запущен**

2. **Откройте Docker Desktop**
   - Нажмите на иконку Docker в трее

3. **Перейдите в настройки**
   - Кликните на шестеренку (Settings)

4. **Включите Kubernetes**
   - Слева выберите "Kubernetes"
   - Поставьте галочку "Enable Kubernetes"
   - Нажмите "Apply & Restart"

5. **Дождитесь запуска**
   - Внизу Docker Desktop появятся две зелёные иконки:
     - 🟢 Docker (Engine is running)
     - 🟢 Kubernetes (is running)
   - Это может занять 2-5 минут

6. **Проверьте**
   ```cmd
   check-k8s.bat
   ```

#### Если Kubernetes не запускается в Docker Desktop:

- Перезапустите Docker Desktop
- Проверьте, что WSL2 включен
- В настройках Docker Desktop попробуйте сбросить Kubernetes:
  - Settings → Kubernetes → Reset Kubernetes Cluster

### ✅ Вариант 2: Minikube

Если Docker Desktop не работает или вы хотите отдельный кластер.

#### Установка:

```powershell
# PowerShell от администратора

# Через Chocolatey
choco install minikube

# Или скачайте с https://minikube.sigs.k8s.io/docs/start/
```

#### Запуск:

```powershell
# Запуск кластера
minikube start

# Проверка статуса
minikube status

# Использование Docker демона Minikube (для образов)
minikube docker-env | Invoke-Expression
```

#### После запуска Minikube:

```cmd
# Проверка
check-k8s.bat

# Собрать образы
build-images.bat

# Развернуть
deploy-k8s.bat
```

### ✅ Вариант 3: Docker Compose (без Kubernetes)

Если Kubernetes не нужен, используйте Docker Compose:

```cmd
# Собрать образы
build-images.bat

# Запустить через Docker Compose
start-docker.bat
```

Это проще и быстрее для локальной разработки!

## Проверка установки

```cmd
# Запустите скрипт проверки
check-k8s.bat
```

Должно быть:
```
✓ kubectl установлен
✓ Kubernetes кластер доступен
```

## Типичные ошибки

### ❌ "dial tcp [::1]:8080: connectex: No connection"

**Проблема:** Kubernetes кластер не запущен

**Решение:**
1. Если используете Docker Desktop - включите Kubernetes в настройках
2. Если используете Minikube - запустите `minikube start`
3. Или используйте Docker Compose вместо Kubernetes

### ❌ "kubectl: command not found"

**Проблема:** kubectl не установлен

**Решение:**
- Docker Desktop устанавливает kubectl автоматически
- Или установите вручную: https://kubernetes.io/docs/tasks/tools/

### ❌ Образы не находятся в Kubernetes

**Проблема:** Kubernetes не видит локальные Docker образы

**Решение для Minikube:**
```bash
# Используйте Docker демон Minikube
minikube docker-env | Invoke-Expression

# Пересоберите образы
build-images.bat
```

**Решение для Docker Desktop:**
- Образы должны быть видны автоматически
- Убедитесь, что в Deployment указано `imagePullPolicy: IfNotPresent`

## Полезные команды

```bash
# Проверка кластера
kubectl cluster-info
kubectl get nodes

# Проверка подов
kubectl get pods -n lab5

# Логи
kubectl logs -f deployment/product-service -n lab5

# Описание пода
kubectl describe pod <pod-name> -n lab5

# Удаление всего
kubectl delete namespace lab5
```

## Рекомендации

### Для обучения и разработки:
1. Начните с **Docker Compose** - это проще
2. Потом попробуйте **Docker Desktop Kubernetes** - встроенный
3. Для продакшена используйте облачные решения (GKE, EKS, AKS)

### Для продакшена:
- Google Kubernetes Engine (GKE)
- Amazon Elastic Kubernetes Service (EKS)
- Azure Kubernetes Service (AKS)

## Дополнительные ресурсы

- [Документация Kubernetes](https://kubernetes.io/docs/home/)
- [Docker Desktop Kubernetes](https://docs.docker.com/desktop/kubernetes/)
- [Minikube Docs](https://minikube.sigs.k8s.io/docs/)

