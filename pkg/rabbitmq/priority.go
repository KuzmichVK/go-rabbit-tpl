package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// InvalidPriorityError signals an out-of-range priority value.
type InvalidPriorityError struct {
	MaxPriority int
	Message     string
}

// Error implements the error interface for InvalidPriorityError.
func (e *InvalidPriorityError) Error() string {
	return e.Message
}

// IsPriorityError reports whether err is an InvalidPriorityError.
func IsPriorityError(err error) bool {
	_, ok := err.(*InvalidPriorityError)
	return ok
}

// DeclarePriorityQueue declares a priority queue with the given max priority.
func DeclarePriorityQueue(ch *amqp.Channel, queueName string, maxPriority int) error {
	if maxPriority < 1 || maxPriority > 255 {
		return &InvalidPriorityError{
			MaxPriority: maxPriority,
			Message:     "maxPriority must be in range 1-255",
		}
	}
	args := amqp.Table{
		// FIX: amqp091 table encoder does not accept a bare int; use int32.
		"x-max-priority": int32(maxPriority),
	}
	_, err := ch.QueueDeclare(queueName, true, false, false, false, args)
	if err != nil {
		log.Printf("failed to declare priority queue %q: %v", queueName, err)
		return err
	}
	log.Printf("priority queue %q declared with maxPriority=%d", queueName, maxPriority)
	return nil
}
