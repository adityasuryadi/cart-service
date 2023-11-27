package config

import (
	"cart-service/exception"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectionClose() {

}

type MqChannel struct {
	exchange   string
	queue      string
	routingKey string
}

func NewRabbitMqConn() (*amqp.Connection, error) {
	return amqp.Dial("amqp://guest:guest@rabbitmq_container:5672/")
}

func (c *MqChannel) ChannelDeclare(data interface{}) {
	conn, err := InitRabbitMQ()
	if err != nil {
		exception.FailOnError(err, "failed to connect Rabbit MQ")
	}

	defer conn.Close()

	ch, err := conn.Channel()
	exception.FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	err = ch.ExchangeDeclare(
		c.exchange, // name
		"direct",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments)
	)
	exception.FailOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q, err := ch.QueueDeclare(
		c.queue, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	exception.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(q.Name, c.routingKey, c.exchange, false, nil)
	if err != nil {
		exception.FailOnError(err, "Failed to declare a queue")
	}

	// body := bodyFrom(os.Args)
	// body := "Hello"
	body, err := json.Marshal(data)
	if err != nil {
		exception.FailOnError(err, "failed convert body")
	}
	err = ch.PublishWithContext(ctx,
		c.exchange,   // exchange
		c.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	exception.FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func ChannelCLose() {

}

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

type ListenFunction func() error

func (c *MqChannel) ListenQueue() (<-chan amqp.Delivery, error) {
	exchange := c.exchange
	queue := c.queue
	routingKey := c.routingKey

	conn, err := InitRabbitMQ()
	if err != nil {
		exception.FailOnError(err, "failed to connect Rabbit MQ")
	}

	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		failOnError(err, "Failed to declare a queue")
	}
	log.Print("producer: declaring binding")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	fmt.Printf("mssgs ", msgs)
	// go func() {
	// 	for d := range msgs {
	// 		// log.Printf(" [x] %s", d.Body)
	// 		var product *entity.Product
	// 		json.Unmarshal(d.Body, &product)

	// 		productRepo.InsertProduct(product)
	// 	}
	// }()
	return msgs, err
}
