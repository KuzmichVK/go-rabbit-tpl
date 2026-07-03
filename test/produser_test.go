package test

import (
	"testing"

	"github.com/lekan-pvp/grade/go-rabbit/internal/models"
	"github.com/lekan-pvp/grade/go-rabbit/pkg/rabbitmq"
)

func TestProducer_PublishJSON(t *testing.T) {
	// Создаём тестовое соединение
	conn, err := rabbitmq.NewConnection("amqp://guest:guest@localhost:5672")
	if err != nil {
		t.Fatal("Failed to connection:", err)
	}
	defer conn.Close()

	monitor := rabbitmq.NewMonitor()
	producer := rabbitmq.NewProducer(conn, monitor)

	// Объявляем тестовую очередь
	err = rabbitmq.DeclarePriorityQueue(conn.Channel(), "test_queue", 5)
	if err != nil {
		t.Fatal("Failed to declare queue:", err)
	}

	// Подготавливаем сообщение
	msg := &models.Message{
		ID:       990,
		Content:  "Test message",
		Priority: 3,
	}

	// Отправляем сообщение
	err = producer.PublishJSON("", "test_queue", msg)
	if err != nil {
		t.Error("Publish failed:", err)
	}

	stats := monitor.Stats()
	if stats["sent"] != 1 {
		t.Errorf("Expected 1 sent message, got %d", stats["sent"])
	}
}
