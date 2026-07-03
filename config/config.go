// Package config holds the application configuration.
package config

// Config is the application configuration.
type Config struct {
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

// RabbitMQConfig holds RabbitMQ connection settings.
type RabbitMQConfig struct {
	URL      string `yaml:"url"`
	Queue    string `yaml:"queue"`
	DLXQueue string `yaml:"dlx_queue"`
	Exchange string `yaml:"exchange"`
}
