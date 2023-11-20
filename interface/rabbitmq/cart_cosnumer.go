package rabbitmq

import (
	"cart-service/entity"
	"cart-service/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	publishMandatory = false
	publishImmediate = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

type CartConsumer struct {
	amqpConn          *amqp.Connection
	productRepository repository.ProductRepository
}

func NewCartConsumer(amqpConn *amqp.Connection, productRepo repository.ProductRepository) *CartConsumer {
	return &CartConsumer{
		amqpConn:          amqpConn,
		productRepository: productRepo,
	}
}

func (c *CartConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, errors.New("error amqpConn.channel")
	}

	fmt.Printf("Declaring exchange: %s", exchangeName)
	err = ch.ExchangeDeclare(
		exchangeName,
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments)
	)

	if err != nil {
		return nil, errors.New("Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)

	if err != nil {
		return nil, errors.New("Error ch.QueueDeclare")
	}

	fmt.Printf("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchangeName,
		bindingKey)

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.New("Error ch.QueueBind")
	}

	fmt.Printf("Queue bound to exchange, starting to consume from queue, consumerTag: %v", consumerTag)

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.New("Error  ch.Qos")
	}

	return ch, nil
}

func (c *CartConsumer) worker(ctx context.Context, messages <-chan amqp.Delivery) {
	for delivery := range messages {
		fmt.Printf("processDeliveries deliveryTag% v", delivery.DeliveryTag)
		var product *entity.Product
		json.Unmarshal(delivery.Body, &product)
		c.productRepository.InsertProduct(product)
	}
	fmt.Println("Deliveries channel closed")
}

func (c *CartConsumer) StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := c.CreateChannel(exchange, queueName, bindingKey, consumerTag)
	if err != nil {
		return errors.New(err.Error())
	}
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return errors.New("Consume")
	}

	for i := 0; i < workerPoolSize; i++ {
		go c.worker(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	fmt.Printf("ch.NotifyClose: %v", chanErr)
	return chanErr
}
