package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// SetupDLX declares a Dead Letter Exchange and wires the main queue to it.
func SetupDLX(ch *amqp.Channel, mainQueue, dlxQueue, exchangeName string) error {
	// 1. Declare the DLX (fanout).
	err := ch.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}
	// 2. Declare the dead-letter queue.
	_, err = ch.QueueDeclare(dlxQueue, true, false, false, false, nil)
	if err != nil {
		return err
	}
	// 3. Bind the DLX to the dead-letter queue (routing key unused for fanout).
	err = ch.QueueBind(dlxQueue, "", exchangeName, false, nil)
	if err != nil {
		return err
	}
	// 4. Declare the main queue routing dead-lettered messages to the DLX.
	args := amqp.Table{
		"x-dead-letter-exchange":    exchangeName,
		"x-dead-letter-routing-key": dlxQueue,
	}
	_, err = ch.QueueDeclare(mainQueue, true, false, false, false, args)
	return err
}
