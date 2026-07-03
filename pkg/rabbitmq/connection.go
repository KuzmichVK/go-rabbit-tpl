// Package rabbitmq implements a small client on top of amqp091-go.
package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection holds an AMQP connection and a channel.
type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewConnection dials RabbitMQ and opens a channel.
func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close() // close the connection if the channel failed
		return nil, err
	}
	return &Connection{conn: conn, channel: ch}, nil
}

// Close closes the connection (the channel is closed with it).
func (c *Connection) Close() error {
	return c.conn.Close()
}

// Channel returns the AMQP channel of this connection.
func (c *Connection) Channel() *amqp.Channel {
	return c.channel
}
