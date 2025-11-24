# Настройка Kubernetes для Lab5

## Что такое Kubernetes?

Kubernetes (K8s) - это система оркестрации контейнеров. Для работы нужен **работающий кластер**.

## Варианты настройки

### Вариант 1: Docker Desktop (самый простой для Windows)

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
     - Docker (Engine is running)
     - Kubernetes (is running)
   - Это может занять 2-5 минут

6. **Проверьте**
   ```cmd
   scripts\check-k8s.bat
   ```

#### Если Kubernetes не запускается в Docker Desktop:

- Перезапустите Docker Desktop
- Проверьте, что WSL2 включен
- В настройках Docker Desktop попробуйте сбросить Kubernetes:
  - Settings → Kubernetes → Reset Kubernetes Cluster

### Вариант 2: Minikube

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
scripts\check-k8s.bat

# Собрать образы
scripts\build-images.bat

# Развернуть
scripts\deploy-k8s.bat
```

### Вариант 3: Docker Compose (без Kubernetes)

Если Kubernetes не нужен, используйте Docker Compose:

```cmd
# Собрать образы
scripts\build-images.bat

# Запустить через Docker Compose
scripts\start-docker.bat
```

Это проще и быстрее для локальной разработки!

## Проверка установки

```cmd
# Запустите скрипт проверки
scripts\check-k8s.bat
```

Должно быть:
```
✓ kubectl установлен
✓ Kubernetes кластер доступен
```