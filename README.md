# User Service (Go, In-Memory)

Микросервис для управления пользователями с хранением данных в памяти (map).

Проект разложен по пакетам: в main только запуск и wiring зависимостей.

## Структура

```text
.
├── cmd
│   └── api
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── httpapi
│   │   └── handler.go
│   └── users
│       ├── inmemory_repository.go
│       ├── repository.go
│       ├── service.go
│       └── user.go
├── go.mod
└── go.sum
```

## Конфигурация через environment variables

- HTTP_PORT: порт HTTP сервера (по умолчанию 8080)

Пример запуска:

```bash
HTTP_PORT=8081 go run ./cmd/api
```

Для PowerShell:

```powershell
$env:HTTP_PORT="8081"
go run ./cmd/api
```

Сервис стартует на http://localhost:{HTTP_PORT}.

## API

### 1) Регистрация пользователя

POST /users

Пример запроса:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Ivan","email":"ivan@example.com"}'
```

### 2) Обновление пользователя

PUT /users/{id}

Пример запроса:

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Ivan Petrov","email":"ivan.petrov@example.com"}'
```

### 3) Удаление пользователя

DELETE /users/{id}

Пример запроса:

```bash
curl -X DELETE http://localhost:8080/users/1
```

### 4) Получение списка пользователей

GET /users

Пример запроса:

```bash
curl http://localhost:8080/users
```

## Модель пользователя

```json
{
  "id": 1,
  "name": "Ivan",
  "email": "ivan@example.com"
}
```

## Коды ответов

- 201 Created - пользователь создан
- 200 OK - успешное получение/обновление
- 204 No Content - пользователь удален
- 400 Bad Request - ошибка в пути или JSON
- 404 Not Found - пользователь не найден
- 405 Method Not Allowed - метод не поддерживается
