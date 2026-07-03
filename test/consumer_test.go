package test

import (
	"testing"
	"time"

	"github.com/lekan-pvp/grade/go-rabbit/internal/models"
	"github.com/lekan-pvp/grade/go-rabbit/pkg/rabbitmq"
)

func TestConsumer_ConsumeJSON(t *testing.T) {
	// Создаём соединение
	conn, err := rabbitmq.NewConnection("amqp://guest:guest@localhost:5672")
	if err != nil {
		t.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	monitor := rabbitmq.NewMonitor()
	consumer := rabbitmq.NewConsumer(conn, monitor)

	// Объявляем очередь
	err = rabbitmq.DeclarePriorityQueue(conn.Channel(), "test_consume", 5)
	if err != nil {
		t.Fatal("Failed to declare queue:", err)
	}

	// Отправляем тестовое сообщение
	producer := rabbitmq.NewProducer(conn, monitor)
	testMsg := &models.Message{
		ID:       1,
		Content:  "test_consume",
		Priority: 1,
	}
	err = producer.PublishJSON("", "test_consume", testMsg)
	if err != nil {
		t.Fatal("Failed to publish test message:", err)
	}

	// Обработчик для теста
	var received bool
	handler := func(msg *models.Message) bool {
		received = true
		return true // ACK
	}

	// Запускаем потребителя на 5 секунд
	err = consumer.ConsumeJSON("test_consume", handler)
	if err != nil {
		t.Fatal("Consume failed:", err)
	}

	time.Sleep(5 * time.Second)

	if !received {
		t.Error("Expected to receive message, but didn't")
	}

	stats := monitor.Stats()
	if stats["received"] != 1 {
		t.Errorf("Ecpected 1 received message, got %d", stats["received"])
	}
}
