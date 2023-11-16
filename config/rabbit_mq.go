package config

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func InitRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq_container:5672/")
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		return nil, err
	}
	return conn, nil
}
