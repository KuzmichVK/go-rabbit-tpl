package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	// 4. Consumer.
	consumer := rabbitmq.NewConsumer(conn, monitor)

	// 5. Message handler.
	handler := func(msg *models.Message) bool {
		fmt.Printf("Received: ID=%d, Content='%s', Priority=%d\n", msg.ID, msg.Content, msg.Priority)
		if msg.ID%2 == 0 {
			fmt.Println("  -> Processing failed, sending to DLX")
			return false // NACK -> DLX
		}
		fmt.Println("  -> Processed successfully")
		return true // ACK
	}

	// 6. Start consuming.
	if err := consumer.ConsumeJSON(cfg.RabbitMQ.Queue, handler); err != nil {
		log.Fatal("Failed to consume messages: ", err)
	}
	fmt.Println("Consumer running... (press Ctrl+C to exit)")

	// 7. Periodic stats.
	tick := time.NewTicker(10 * time.Second)
	go func() {
		for range tick.C {
			stats := monitor.Stats()
			log.Printf("Stats: received=%d, errors=%d", stats["received"], stats["errors"])
		}
	}()

	// 8. Block forever.
	select {}
}
