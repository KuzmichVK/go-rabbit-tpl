# GoRabbit — клиент RabbitMQ на Go

Учебный проект (проверочная работа модуля 1 курса «Работа с брокерами сообщений»,
Яндекс Грейд). Небольшая библиотека поверх `amqp091-go` для публикации и получения
сообщений в формате JSON с поддержкой приоритетных очередей, механизма DLX и
встроенного мониторинга.

## Возможности

- **JSON** — автоматическая сериализация структур Go в JSON и обратно.
- **Приоритетные очереди** — `x-max-priority`: более важные сообщения обрабатываются
  раньше.
- **DLX (Dead Letter Exchange)** — «очередь мёртвых писем» для сообщений, получивших
  NACK/Reject или просроченных.
- **Мониторинг** — счётчики отправленных, полученных сообщений и ошибок.
- **Конфигурация из YAML** — параметры подключения в `config/config.yml`.
- **Готовность к запуску** — `docker-compose.yml` + примеры в `cmd/`.

## Структура проекта

```
go-rabbit-tpl/
├── cmd/
│   ├── producer/main.go      # демо-продюсер
│   └── consumer/main.go      # демо-потребитель
├── pkg/rabbitmq/
│   ├── connection.go         # соединение + канал
│   ├── monitor.go            # счётчики статистики
│   ├── producer.go           # PublishJSON
│   ├── consumer.go           # ConsumeJSON (ACK/NACK)
│   ├── priority.go           # DeclarePriorityQueue
│   └── dlx.go                # SetupDLX
├── internal/models/
│   └── message.go            # Message + ToJSON/FromJSON
├── config/
│   ├── config.go             # структуры Config
│   └── config.yml            # значения конфигурации
├── test/
│   ├── produser_test.go      # тесты (шаблон, не менять)
│   └── consumer_test.go
├── docker-compose.yml
├── go.mod
└── README.md
```

## Требования

- **Go** (модуль требует `go 1.24.5`).
- **Docker** (Docker Desktop / colima) — для запуска RabbitMQ.

## Быстрый старт

Запускать из корня репозитория.

```bash
# 1. Поднять RabbitMQ (веб-панель: http://localhost:15672, guest/guest)
docker compose up -d

# 2. Подтянуть зависимости (создаст go.sum)
go mod tidy

# 3. Форматирование и статические проверки
gofmt -w .
go vet ./...

# 4. Сборка
go build ./...

# 5. Тесты (RabbitMQ из шага 1 должен быть запущен — тесты реально коннектятся)
go test ./...
```

> ⚠️ Тесты **не мокают** брокер: `TestProducer_PublishJSON` и
> `TestConsumer_ConsumeJSON` подключаются к `amqp://guest:guest@localhost:5672`.
> Без запущенного RabbitMQ (`docker compose up -d`) они упадут на `NewConnection`.

## Конфигурация

`config/config.yml` — единственный источник настроек (код читает именно его):

```yaml
rabbitmq:
  url: "amqp://guest:guest@localhost:5672"
  queue: "main_queue"
  dlx_queue: "dlx_queue"
  exchange: "dlx_exchange"
```

## Запуск демо

```bash
# в одном терминале — потребитель
go run ./cmd/consumer

# в другом — продюсер (отправит одно сообщение с приоритетом)
go run ./cmd/producer
```

Потребитель имитирует ошибку для сообщений с чётным `ID` (возвращает `false` → NACK
→ сообщение уходит в DLX).

## ⚠️ Известные ограничения

**Демо-продюсер и связка «приоритеты + DLX» на одной очереди.**
`cmd/producer/main.go` вызывает подряд `DeclarePriorityQueue(main_queue, 10)` и
`SetupDLX(main_queue, …)`. Обе функции объявляют **одну и ту же** очередь
`main_queue`, но с **разными** аргументами (`x-max-priority` против
`x-dead-letter-*`). Повторный `QueueDeclare` очереди с несовпадающими аргументами
даёт на брокере ошибку `PRECONDITION_FAILED` (inequivalent args) — второй вызов
залогируется как ошибка, а канал может стать невалидным.

Это ограничение самого демо, а **не** тестов: `produser_test.go` и
`consumer_test.go` объявляют приоритетную очередь **без** DLX, поэтому проходят
штатно.

Если нужно запустить демо с обоими механизмами сразу — варианты:

- объявлять `main_queue` **один раз** со всеми аргументами
  (`x-max-priority` + `x-dead-letter-exchange` + `x-dead-letter-routing-key`);
- либо использовать **разные** очереди для демонстрации приоритетов и DLX.

## Внесённые исправления (относительно авторского решения)

При оформлении авторского решения устранены ошибки, мешавшие сборке/работе (см.
`CHANGELOG.md` с деталями):

- ⚠️ **`go.mod`, module path.** Шаблонный `github.com/PracticumGrade/go-rabbit-tpl`
  заменён на `github.com/lekan-pvp/grade/go-rabbit` — именно этот путь импортируют
  **тесты** (`test/*.go`). Без совпадения `go test ./...` не компилируется.
- ⚠️ **`priority.go`.** `x-max-priority: maxPriority` (голый `int`) → `int32(maxPriority)`:
  кодировщик AMQP-таблицы в `amqp091-go` не принимает `int`, что приводило бы к
  ошибке при объявлении очереди.
- ⚠️ **`docker-compose.yml`.** Опечатка `RABBIYMQ_DEFAULT_PASS` → `RABBITMQ_DEFAULT_PASS`.
- ⚠️ **`dlx.go`.** Опечатка в имени параметра `exchanngeName` → `exchangeName`.
- ⚠️ **Конфиг.** Убран лишний пустой `config.yml` в корне; единый рабочий конфиг —
  `config/config.yml` (путь, который читает код в `cmd/`).

Оставлено как в авторском решении (работает, на прохождение тестов не влияет):

- `channel.Publish(...)` помечен как устаревший в `amqp091-go` (рекомендуется
  `PublishWithContext`), но сохранён для соответствия авторскому варианту.
