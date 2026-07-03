package rabbitmq

import "sync"

// Monitor tracks producer/consumer statistics.
type Monitor struct {
	messagesSent     int64
	messagesReceived int64
	errors           int64
	mu               sync.Mutex
}

// NewMonitor creates a new Monitor.
func NewMonitor() *Monitor {
	return &Monitor{}
}

// IncSent increments the sent counter.
func (m *Monitor) IncSent() {
	m.mu.Lock()
	m.messagesSent++
	m.mu.Unlock()
}

// IncReceived increments the received counter.
func (m *Monitor) IncReceived() {
	m.mu.Lock()
	m.messagesReceived++
	m.mu.Unlock()
}

// IncError increments the error counter.
func (m *Monitor) IncError() {
	m.mu.Lock()
	m.errors++
	m.mu.Unlock()
}

// Stats returns the current statistics.
func (m *Monitor) Stats() map[string]int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return map[string]int64{
		"sent":     m.messagesSent,
		"received": m.messagesReceived,
		"errors":   m.errors,
	}
}
