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

### Конфигурация через environment variables

Сервис загружает переменные окружения из `.env` через библиотеку `godotenv`.

1. Скопируйте шаблон:

```powershell
Copy-Item .env.example .env
```

2. Измените значения в `.env` под ваше окружение.

Поддерживаемые переменные:

- `GRPC_PORT`
- `GRPC_ADDR`
- `METRICS_ADDR`
- `AUTH_PASSWORD`
- `JWT_SECRET`
- `JWT_TTL_HOURS`
- `KAFKA_BROKER`
- `KAFKA_TOPIC`
- `KAFKA_GROUP`
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`

Примечание: `.env` уже игнорируется в git, поэтому секреты не попадут в репозиторий.

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

## Мониторинг: Prometheus + Grafana

### 1) Запустите сервис с метриками

Метрики уже поднимаются в `cmd/api/main.go` на адресе `:2112/metrics`.

```powershell
go run ./cmd/api
```

Проверьте, что endpoint доступен:

```powershell
Invoke-WebRequest http://localhost:2112/metrics | Select-Object -ExpandProperty StatusCode
```

Ожидается код `200`.

### 2) Запустите Prometheus и Grafana

В репозитории добавлены:

- `docker-compose.monitoring.yml`
- `deploy/monitoring/prometheus.yml`

Запуск:

```powershell
docker compose -f docker-compose.monitoring.yml up -d
```

UI:

- Prometheus: `http://localhost:9091`
- Grafana: `http://localhost:3000` (login/password: `admin` / `admin`)

### 3) Подключите Prometheus в Grafana

1. Откройте Grafana -> Connections -> Data sources -> Add data source.
2. Выберите Prometheus.
3. URL: `http://prometheus:9090`
4. Save & test.

### 4) Постройте графики

Примеры запросов PromQL для панели:

- RPS по методам:
	- `sum by (method) (rate(grpc_requests_total[1m]))`
- P95 latency по методам:
	- `histogram_quantile(0.95, sum by (le, method) (rate(grpc_request_duration_seconds_bucket[5m])))`
- Средняя задержка по методам:
	- `sum by (method) (rate(grpc_request_duration_seconds_sum[5m])) / sum by (method) (rate(grpc_request_duration_seconds_count[5m]))`

### 5) Остановить мониторинг

```powershell
docker compose -f docker-compose.monitoring.yml down
```
