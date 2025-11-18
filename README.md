# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюеров на Pull Request'ы.

## Описание

Микросервис, который автоматически назначает ревьюеров на Pull Request'ы из команды автора, позволяет выполнять переназначение ревьюверов и получать список PR'ов, назначенных конкретному пользователю, а также управлять командами и активностью пользователей.

### Локальный запуск

1. Клонируйте репозиторий
2. Установите зависимости:
   ```
   go mod download
   ```
3. Соберите приложение:
   ```
   make build
   # или linux
   go build -o bin/server ./cmd/server
   # или windows
   go build -o bin/server.exe ./cmd/server
   ```
4. Запустите сервер:
   ```bash
   make run
   # или linux
   ./bin/server
   # или windows
   ./bin/server.exe
   ```

Сервер будет доступен на `http://localhost:8080`

### Запуск через Docker Compose

```bash
docker-compose up
```

Сервис будет доступен на `http://localhost:8080`

## API Endpoints

Подробная спецификация API доступна в файле `openapi.yml`.

### Teams

#### `POST /team/add` - Создать команду с участниками

Создаёт новую команду и добавляет/обновляет пользователей в этой команде. Если команда уже существует, вернёт ошибку.

**Запрос bash | Linux:**
```bash
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": true
      },
      {
        "user_id": "u2",
        "username": "Bob",
        "is_active": true
      },
      {
        "user_id": "u3",
        "username": "Charlie",
        "is_active": true
      }
    ]
  }'
```

**Запрос PowerShell | Windows:**
```PowerShell
curl.exe -X POST http://localhost:8080/team/add `
  -H "Content-Type: application/json" `
  -d '{\"team_name\": \"backend\", \"members\": [{\"user_id\": \"u1\", \"username\": \"Alice\", \"is_active\": true}, {\"user_id\": \"u2\", \"username\": \"Bob\", \"is_active\": true}, {\"user_id\": \"u3\", \"username\": \"Charlie\", \"is_active\": true}]}'
```

**Успешный ответ (201):**
```json
{
  "team": {
    "team_name": "backend",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": true
      },
      {
        "user_id": "u2",
        "username": "Bob",
        "is_active": true
      },
      {
        "user_id": "u3",
        "username": "Charlie",
        "is_active": true
      }
    ]
  }
}
```

**Ошибка (400):** Команда уже существует
```json
{
  "error": {
    "code": "TEAM_EXISTS",
    "message": "team_name already exists"
  }
}
```

#### `GET /team/get` - Получить команду с участниками

Возвращает информацию о команде и всех её участниках (включая неактивных).

**Запрос bash | Linux:**
```bash
curl -X GET http://localhost:8080/team/get?team_name=backend
```

**Запрос PowerShell | Windows:**
```PowerShell
curl.exe http://localhost:8080/team/get?team_name=backend 
```

**Успешный ответ (200):**
```json
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "u1",
      "username": "Alice",
      "is_active": true
    },
    {
      "user_id": "u2",
      "username": "Bob",
      "is_active": true
    }
    {
      "user_id": "u3",
      "username": "Charlie",
      "is_active": true
    }
  ]
}
```

**Ошибка (404):** Команда не найдена
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "resource not found"
  }
}
```

### Users

#### `POST /users/setIsActive` - Установить флаг активности пользователя

Изменяет статус активности пользователя. Неактивные пользователи (`is_active: false`) не назначаются на новые ревью, но их текущие назначения остаются видимыми.

**Запрос bash | Linux:**
```bash
curl -X POST http://localhost:8080/users/setIsActive \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "u2",
    "is_active": false
  }'
```

**Запрос PowerShell | Windows:**
```PowerShell
curl.exe -X POST http://localhost:8080/users/setIsActive `
  -H "Content-Type: application/json" `
  -d '{\"user_id\": \"u2\",\"is_active\": false}'
```

**Успешный ответ (200):**
```json
{
  "user": {
    "user_id": "u2",
    "username": "Bob",
    "team_name": "backend",
    "is_active": false
  }
}
```

**Ошибка (404):** Пользователь не найден
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "resource not found"
  }
}
```

#### `GET /users/getReview` - Получить PR'ы, где пользователь назначен ревьюером

Возвращает список всех Pull Request'ов (открытых и закрытых), где указанный пользователь назначен ревьюером.

**Запрос bash | Linux:**
```bash
curl -X GET http://localhost:8080/users/getReview?user_id=u2
```

**Запрос PowerShell | Windows:**
```PowerShell
curl.exe -X GET http://localhost:8080/users/getReview?user_id=u2
```

**Успешный ответ (200):**
```json
{
  "user_id": "u2",
  "pull_requests": [
    {
      "pull_request_id": "pr-1001",
      "pull_request_name": "Add search",
      "author_id": "u1",
      "status": "OPEN"
    },
    {
      "pull_request_id": "pr-1002",
      "pull_request_name": "Fix bug",
      "author_id": "u3",
      "status": "MERGED"
    }
  ]
}
```

**Ошибка (404):** Пользователь не найден

### Pull Requests

#### `POST /pullRequest/create` - Создать PR и автоматически назначить ревьюеров

Создаёт новый Pull Request и автоматически назначает до 2 активных ревьюверов из команды автора (исключая самого автора). Если доступных кандидатов меньше двух, назначается доступное количество (0 или 1).

**Запрос bash | Linux:**
```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search feature",
    "author_id": "u1"
  }'
```

**Запрос PowerShell | Windows:**
```PowerShell
curl.exe -X POST http://localhost:8080/pullRequest/create `
  -H "Content-Type: application/json" `
  -d '{\"pull_request_id\": \"pr-001\", \"pull_request_name\": \"Add search feature\", \"author_id\": \"u1\"}'
```

**Успешный ответ (201):**
```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search feature",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": ["u2", "u3"],
    "createdAt": "2025-10-24T12:34:56Z",
    "mergedAt": null
  }
}
```

**Ошибки:**
- **404:** Автор или команда не найдены
- **409:** PR с таким ID уже существует
```json
{
  "error": {
    "code": "PR_EXISTS",
    "message": "PR id already exists"
  }
}
```

#### `POST /pullRequest/merge` - Пометить PR как MERGED

Помечает Pull Request как объединённый (MERGED). Операция идемпотентна — повторный вызов не приводит к ошибке и возвращает актуальное состояние PR. После merge изменение списка ревьюверов запрещено.

**Запрос bash | Linux:**
```bash
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001"
  }'
```

**Запрос PowerShell | Windows:**
```PowerShell
curl.exe -X POST http://localhost:8080/pullRequest/merge `
  -H "Content-Type: application/json" `
  -d '{\"pull_request_id\": \"pr-001\"}'
```

**Успешный ответ (200):**
```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search feature",
    "author_id": "u1",
    "status": "MERGED",
    "assigned_reviewers": ["u2", "u3"],
    "createdAt": "2025-10-24T12:34:56Z",
    "mergedAt": "2025-10-24T15:20:10Z"
  }
}
```

**Ошибка (404):** PR не найден

#### `POST /pullRequest/reassign` - Переназначить ревьюера

Заменяет одного ревьювера на случайного активного участника из команды заменяемого ревьювера. Нельзя переназначать ревьюверов в уже объединённых (MERGED) PR.

**Запрос bash | Linux:**
```bash
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1001",
    "old_user_id": "u2"
  }'
```

**Запрос PowerShell | Windows:**
```PowerShell
curl.exe -X POST http://localhost:8080/pullRequest/reassign `
  -H "Content-Type: application/json" `
  -d '{\"pull_request_id\": \"pr-001\", \"old_user_id\": \"u2\"}'
```

**Успешный ответ (200):**
```json
{
  "pr": {
    "pull_request_id": "pr-1001",
    "pull_request_name": "Add search feature",
    "author_id": "u1",
    "status": "OPEN",
    "assigned_reviewers": ["u3", "u5"],
    "createdAt": "2025-10-24T12:34:56Z",
    "mergedAt": null
  },
  "replaced_by": "u5"
}
```

**Ошибки:**
- **404:** PR или пользователь не найдены
- **409:** Нарушение доменных правил:
  - `PR_MERGED` - нельзя менять ревьюверов после merge
  - `NOT_ASSIGNED` - указанный пользователь не был назначен ревьювером на этот PR
  - `NO_CANDIDATE` - нет доступных активных кандидатов в команде заменяемого ревьювера

Примеры ошибок:
```json
{
  "error": {
    "code": "PR_MERGED",
    "message": "cannot reassign on merged PR"
  }
}
```

```json
{
  "error": {
    "code": "NOT_ASSIGNED",
    "message": "reviewer is not assigned to this PR"
  }
}
```

```json
{
  "error": {
    "code": "NO_CANDIDATE",
    "message": "no active replacement candidate in team"
  }
}
```

### Health Check

#### `GET /healthz` - Проверка работоспособности сервиса

Простой эндпоинт для проверки доступности сервиса.

**Запрос bash | Linux:**
```bash
curl -X GET http://localhost:8080/healthz
```

**Запрос PowerShell | Windows:**
```PowerShell
curl.exe -X GET http://localhost:8080/healthz
```

**Ответ (200):**
```
OK
```

## Структура проекта

```
.
├── cmd/
│   └── server/          # Точка входа приложения
├── internal/
│   ├── api/            # HTTP handlers и роутинг
│   ├── app/            # Инициализация приложения
│   ├── config/         # Конфигурация
│   ├── domain/         # Доменные модели и ошибки
│   ├── httpserver/     # HTTP сервер
│   ├── repository/     # Интерфейсы и реализации репозиториев
│   │   └── inmemory/   # In-memory реализация 
│   └── service/        # Бизнес-логика
├── pkg/
│   ├── httpx/          # HTTP утилиты
│   └── logger/         # Логирование
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── README.md
```

## Конфигурация

Приложение использует переменные окружения:

- `HTTP_ADDR` - адрес для HTTP сервера (по умолчанию: `:8080`)
- `HTTP_READ_HEADER_TIMEOUT` - таймаут чтения заголовков (по умолчанию: `5s`)
- `HTTP_SHUTDOWN_TIMEOUT` - таймаут graceful shutdown (по умолчанию: `5s`)

## Реализованные функции

- Создание команд и управление пользователями
- Автоматическое назначение до 2 ревьюеров при создании PR
- Переназначение ревьюеров из команды заменяемого ревьюера
- Идемпотентная операция merge PR
- Получение списка PR'ов для пользователя
- Управление активностью пользователей
- In-memory хранилище для быстрого тестирования

## Принятые решения

### Хранение данных

В текущей реализации используется in-memory хранилище для упрощения разработки и тестирования. 

### Транзакции

Для in-memory реализации транзакции выполняются синхронно без реальной изоляции. При переходе на PostgreSQL необходимо использовать реальные транзакции БД.

## Дополнительные задания

Следующие задания из ТЗ реализованы:

- Интеграционное/E2E тестирование (реализовано в `test/integration_test.go`)


## Тестирование

Для проверки всех условий из ТЗ доступно несколько способов:

### 1. Интеграционные тесты

Запустите сервер в одном терминале:
```bash
make run
# или
docker-compose up
```

В другом терминале запустите тесты:
```bash
make test-integration
# или
go test ./test/integration_test.go -v
```

### 2. Ручное тестирование через curl

Примеры запросов и подробные инструкции см. в `test/README.md`.

## Разработка

### Сборка

```bash
make build
```

### Запуск

```bash
make run
```

### Тесты

```bash
make test-integration  # Интеграционные тесты (требует запущенный сервер)
```

### Очистка

```bash
make clean
```



