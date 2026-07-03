# CHANGELOG — go-rabbit-tpl

Проект: клиент RabbitMQ на Go (проверочная работа модуля 1, курс «Работа с брокерами
сообщений», Яндекс Грейд). Репозиторий-шаблон форкнут/склонирован, все файлы
реализованы по авторскому решению из задания с исправлением ошибок.

---

## Что сделано

Реализованы все файлы (шаблон приходил с пустыми заглушками):

| Файл | Содержимое |
|---|---|
| `go.mod` | module + зависимости (`amqp091-go`, `yaml.v2`) |
| `internal/models/message.go` | `Message{ID, Content, Priority}`, `ToJSON`, `FromJSON` |
| `config/config.go` | `Config` / `RabbitMQConfig` (YAML-теги) |
| `config/config.yml` | значения подключения |
| `pkg/rabbitmq/connection.go` | `Connection`, `NewConnection`, `Close`, `Channel` |
| `pkg/rabbitmq/monitor.go` | `Monitor` + `IncSent/IncReceived/IncError`, `Stats` |
| `pkg/rabbitmq/producer.go` | `Producer`, `NewProducer`, `PublishJSON` |
| `pkg/rabbitmq/consumer.go` | `Consumer`, `NewConsumer`, `ConsumeJSON`, `Close` |
| `pkg/rabbitmq/priority.go` | `DeclarePriorityQueue`, `InvalidPriorityError`, `IsPriorityError` |
| `pkg/rabbitmq/dlx.go` | `SetupDLX` |
| `cmd/producer/main.go` | демо-продюсер |
| `cmd/consumer/main.go` | демо-потребитель |
| `docker-compose.yml` | RabbitMQ 3-management |
| `README.md` | документация проекта |

Файлы тестов `test/produser_test.go` и `test/consumer_test.go` — **не трогались**
(они контракт: задают сигнатуры, поля `Message`, ключи `Stats`, module-путь импортов).

---

## Исправления багов авторского решения

1. **`go.mod` — module path (критично для тестов).**
   Было (шаблон): `module github.com/PracticumGrade/go-rabbit-tpl`.
   Стало: `module github.com/lekan-pvp/grade/go-rabbit`.
   Причина: оба теста импортируют пакеты по пути
   `github.com/lekan-pvp/grade/go-rabbit/{internal/models,pkg/rabbitmq}`. При
   несовпадении module-пути `go test ./...` не собирается. Тесты — авторитет
   контракта, поэтому module приведён под них, и все внутренние импорты кода — тоже.

2. **`pkg/rabbitmq/priority.go` — тип значения `x-max-priority`.**
   Было: `"x-max-priority": maxPriority` (голый `int`).
   Стало: `"x-max-priority": int32(maxPriority)`.
   Причина: кодировщик AMQP-таблицы в `amqp091-go` не поддерживает `int` →
   объявление очереди завершилось бы ошибкой. Помечено в коде `// FIX:`.

3. **`docker-compose.yml` — опечатка в имени переменной.**
   Было: `RABBIYMQ_DEFAULT_PASS: guest`.
   Стало: `RABBITMQ_DEFAULT_PASS: guest`.

4. **`pkg/rabbitmq/dlx.go` — опечатка в имени параметра.**
   Было: `exchanngeName` (компилировалось, но грязно).
   Стало: `exchangeName`.

5. **Конфиг — двойной источник.**
   В авторском варианте фигурировали и `config.yml` (в корне), и `config/config.yaml`
   в структуре, а код читает `config/config.yml`. Оставлен **единственный**
   `config/config.yml`; лишний пустой корневой `config.yml` из шаблона удалён.

**Оставлено как в авторском (работает, тесты проходят):**
`channel.Publish(...)` устарел в `amqp091-go` (рекомендуется `PublishWithContext`) —
сохранён для соответствия авторскому решению; отмечено комментарием в коде.

---

## Известное ограничение (не баг сборки)

`cmd/producer/main.go` вызывает `DeclarePriorityQueue(main_queue)` **и**
`SetupDLX(main_queue)` — оба объявляют одну очередь с разными аргументами →
`PRECONDITION_FAILED` при повторном `QueueDeclare`. На тесты не влияет (они не
комбинируют приоритеты и DLX на одной очереди). Обход описан в `README.md`
(«Известные ограничения»).

---

## Команды (раскладка → сборка → проверка → git)

Выполнялись из корня склонированного репозитория. **Перед git-операциями —
поставить синхронизацию Яндекс.Диска на паузу** (репозиторий лежит на пути
Яндекс.Диска; одновременная синхронизация может повредить `.git/index`).

```bash
# раскладка файлов — heredoc-командами (см. историю чата), затем:
rm -f config.yml                 # убрать лишний пустой корневой конфиг

docker compose up -d             # RabbitMQ для тестов
go mod tidy                      # зависимости + go.sum
gofmt -w .                       # формат (табы)
go vet ./...                     # статический анализ
go build ./...                   # сборка
go test ./...                    # тесты (нужен запущенный RabbitMQ)

# git (при выключенной синхронизации Яндекс.Диска):
git add -A
git commit -m "feat: реализация клиента GoRabbit (producer/consumer, priority, DLX, monitor, config)"
git push
```
