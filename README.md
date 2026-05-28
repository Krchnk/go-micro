# User Service (gRPC, Protobuf, In-Memory)

Микросервис управления пользователями на gRPC с бинарной сериализацией Protocol Buffers.

## Структура

```text
.
├── api
│   └── proto
│       └── users
│           └── v1
│               └── users.proto
├── cmd
│   ├── api
│   │   └── main.go
│   ├── client
│   │   └── main.go
│   └── perf
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── gen
│   │   └── users
│   │       └── v1
│   │           ├── users.pb.go
│   │           └── users_grpc.pb.go
│   ├── grpcapi
│   │   └── server.go
│   └── users
│       ├── inmemory_repository.go
│       ├── repository.go
│       ├── service.go
│       └── user.go
├── go.mod
└── go.sum
```

## gRPC методы

Сервис `users.v1.UserService` реализует:

- `CreateUser`
- `UpdateUser`
- `DeleteUser`
- `ListUsers`

Описание контрактов находится в `api/proto/users/v1/users.proto`.

## Конфигурация

- `GRPC_PORT`: порт gRPC-сервера (по умолчанию `9090`)
- `GRPC_ADDR`: адрес для gRPC-клиента (по умолчанию `localhost:9090`)

## Запуск сервера

```powershell
$env:GRPC_PORT="9090"
go run ./cmd/api
```

## Запуск тестового клиента (CRUD сценарий)

```powershell
$env:GRPC_ADDR="localhost:9090"
go run ./cmd/client
```

Клиент выполняет последовательность:
1. `CreateUser`
2. `ListUsers`
3. `UpdateUser`
4. `DeleteUser`
5. `ListUsers`

## Исследование производительности (JSON vs Protobuf)

```powershell
go run ./cmd/perf
```

Пример результата (100000 итераций, 100 пользователей в payload):

- `Payload size (JSON): 5687 bytes`
- `Payload size (Proto): 3384 bytes`
- `Marshal JSON: 1.1884992s`
- `Marshal Proto: 693.6168ms`
- `Unmarshal JSON: 7.2767115s`
- `Unmarshal Proto: 1.6785275s`

Вывод: protobuf payload меньше и сериализация/десериализация выполняется быстрее, особенно при чтении (unmarshal), что снижает сетевую нагрузку и latency при росте нагрузки.
