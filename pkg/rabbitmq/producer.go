package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/lekan-pvp/grade/go-rabbit/internal/models"
)

// Producer publishes messages to RabbitMQ.
type Producer struct {
	connection *Connection
	monitor    *Monitor
}

// NewProducer creates a new Producer.
func NewProducer(conn *Connection, monitor *Monitor) *Producer {
	return &Producer{connection: conn, monitor: monitor}
}

// PublishJSON publishes a JSON message to the given queue through the exchange.
func (p *Producer) PublishJSON(exchange, queue string, msg *models.Message) error {
	body, err := msg.ToJSON()
	if err != nil {
		p.monitor.IncError()
		return err
	}
	// NOTE: Publish is deprecated in amqp091-go; PublishWithContext is preferred.
	err = p.connection.channel.Publish(
		exchange,
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Priority:    uint8(msg.Priority),
		},
	)
	if err == nil {
		p.monitor.IncSent()
	} else {
		p.monitor.IncError()
	}
	return err
}
