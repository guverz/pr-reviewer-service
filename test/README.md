# Тестирование проекта

Этот каталог содержит тесты и скрипты для проверки всех условий из ТЗ.

## Условия для проверки

1. ✅ При создании PR автоматически назначаются до 2 активных ревьюеров из команды автора, исключая самого автора
2. ✅ Переназначение заменяет одного ревьюера на случайного активного участника из команды заменяемого ревьюера
3. ✅ После MERGED менять список ревьюверов нельзя
4. ✅ Если доступных кандидатов меньше двух, назначается доступное количество (0/1)
5. ✅ Пользователь с `isActive = false` не должен назначаться на ревью
6. ✅ Операция merge должна быть идемпотентной
7. ✅ Сервис должен подниматься через `docker-compose up` на порту 8080

## Способы тестирования

### 1. Интеграционные тесты (Go)

Запустите сервер в одном терминале:
```bash
make run
# или
docker-compose up
```

В другом терминале запустите тесты:
```bash
go test ./test/integration_test.go -v
```

### 2. Bash скрипт для ручного тестирования

Запустите сервер:
```bash
make run
# или
docker-compose up
```

### 3. Ручное тестирование через curl

Примеры запросов для проверки каждого условия:

#### Создание команды
```bash
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend",
    "members": [
      {"user_id": "u1", "username": "Alice", "is_active": true},
      {"user_id": "u2", "username": "Bob", "is_active": true},
      {"user_id": "u3", "username": "Charlie", "is_active": true},
      {"user_id": "u4", "username": "David", "is_active": false}
    ]
  }'
```

#### Создание PR (проверка автоматического назначения)
```bash
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-001",
    "pull_request_name": "Test PR",
    "author_id": "u1"
  }'
```

#### Переназначение ревьюера
```bash
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-001",
    "old_user_id": "u2"
  }'
```

#### Merge PR
```bash
curl -X POST http://localhost:8080/pullRequest/merge \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-001"
  }'
```

#### Попытка переназначения после merge (должна быть ошибка 409)
```bash
curl -X POST http://localhost:8080/pullRequest/reassign \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-001",
    "old_user_id": "u3"
  }'
```

## Проверка через Docker Compose

1. Убедитесь, что сервис запущен:
```bash
docker-compose up -d
```

2. Проверьте health check:
```bash
curl http://localhost:8080/healthz
```

3. Запустите тесты или скрипт (сервер должен быть доступен на localhost:8080)





