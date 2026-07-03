// Package models contains data structures exchanged via RabbitMQ.
package models

import "encoding/json"

// Message is a message transferred through RabbitMQ.
type Message struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Priority int    `json:"priority,omitempty"`
}

// ToJSON serializes Message into JSON.
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON deserializes JSON into a Message.
func FromJSON(data []byte) (*Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	return &m, err
}
