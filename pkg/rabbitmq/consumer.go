package rabbitmq

import (
	"log"

	"github.com/lekan-pvp/grade/go-rabbit/internal/models"
)

// Consumer receives and processes messages from RabbitMQ.
type Consumer struct {
	connection *Connection
	monitor    *Monitor
}

// NewConsumer creates a new Consumer.
func NewConsumer(conn *Connection, monitor *Monitor) *Consumer {
	return &Consumer{connection: conn, monitor: monitor}
}

// ConsumeJSON subscribes to a queue and dispatches each JSON message to handler.
// handler returns true to ACK (message removed) or false to NACK (may go to DLX).
func (c *Consumer) ConsumeJSON(queue string, handler func(*models.Message) bool) error {
	msgs, err := c.connection.channel.Consume(
		queue,
		"",    // consumer tag (anonymous)
		false, // autoAck: false -> manual ACK/NACK
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range msgs {
			msg, err := models.FromJSON(delivery.Body)
			if err != nil {
				c.monitor.IncError()
				log.Printf("failed to parse JSON message: %v", err)
				_ = delivery.Reject(false)
				continue
			}
			if handler(msg) {
				if err := delivery.Ack(false); err != nil {
					log.Printf("failed to send ACK: %v", err)
				} else {
					c.monitor.IncReceived()
				}
			} else {
				// NACK via Reject with requeue=false -> may be dead-lettered to DLX.
				if err := delivery.Reject(false); err != nil {
					log.Printf("failed to send NACK: %v", err)
				}
			}
		}
	}()
	return nil
}

// Close stops the consumer by closing the channel.
func (c *Consumer) Close() error {
	return c.connection.channel.Close()
}
