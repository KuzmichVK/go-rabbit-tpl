package main

import (
	"log"
	"os"

	"github.com/lekan-pvp/grade/go-rabbit/config"
	"github.com/lekan-pvp/grade/go-rabbit/internal/models"
	"github.com/lekan-pvp/grade/go-rabbit/pkg/rabbitmq"
	"gopkg.in/yaml.v2"
)

func main() {
	// 1. Read config.
	cfgData, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Fatal("Failed to read config: ", err)
	}
	var cfg config.Config
	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		log.Fatal("Failed to parse config: ", err)
	}

	// 2. Connect.
	conn, err := rabbitmq.NewConnection(cfg.RabbitMQ.URL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}
	defer conn.Close()

	// 3. Monitor.
	monitor := rabbitmq.NewMonitor()

	// 4. Producer.
	producer := rabbitmq.NewProducer(conn, monitor)

	ch := conn.Channel()

	// 5. Priority queue.
	if err := rabbitmq.DeclarePriorityQueue(ch, cfg.RabbitMQ.Queue, 10); err != nil {
		log.Printf("Failed to declare priority queue: %v", err)
	}

	// 6. DLX.
	// NOTE: re-declaring the same main queue with different args (DLX) than the
	// priority declare above triggers PRECONDITION_FAILED. See README "Known
	// issues". The tests never combine priority + DLX on one queue.
	if err := rabbitmq.SetupDLX(ch, cfg.RabbitMQ.Queue, cfg.RabbitMQ.DLXQueue, cfg.RabbitMQ.Exchange); err != nil {
		log.Printf("Failed to setup DLX: %v", err)
	}

	// 7. Sample message.
	msg := &models.Message{ID: 1, Content: "Hello with priority", Priority: 5}

	// 8. Publish.
	if err := producer.PublishJSON("", cfg.RabbitMQ.Queue, msg); err != nil {
		log.Printf("Failed to publish message: %v", err)
	} else {
		log.Println("Message sent successfully!")
	}

	// 9. Stats.
	stats := monitor.Stats()
	log.Printf("Stats: sent=%d, errors=%d", stats["sent"], stats["errors"])
}
