package main

import (
	"context"
	"fmt"
	"log"

	"cart-service/config"
	"cart-service/controller"
	"cart-service/exception"
	"cart-service/interface/rabbitmq"
	"cart-service/repository"
	"cart-service/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// func ListenQueue() <-chan amqp.Delivery {
// 	exchange := "product.created"
// 	queue := "product.create"
// 	routingKey := "create"

// 	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq_container:5672/")
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	err = ch.ExchangeDeclare(
// 		exchange, // name
// 		"direct", // type
// 		true,     // durable
// 		false,    // auto-deleted
// 		false,    // internal
// 		false,    // no-wait
// 		nil,      // arguments
// 	)
// 	failOnError(err, "Failed to declare an exchange")

// 	q, err := ch.QueueDeclare(
// 		queue, // name
// 		false, // durable
// 		false, // delete when unused
// 		false, // exclusive
// 		false, // no-wait
// 		nil,   // arguments
// 	)
// 	failOnError(err, "Failed to declare a queue")

// 	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
// 	if err != nil {
// 		failOnError(err, "Failed to declare a queue")
// 	}
// 	log.Print("producer: declaring binding")

// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		true,   // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	failOnError(err, "Failed to register a consumer")
// 	// go func() {
// 	// for d := range msgs {
// 	// 	var product *entity.Product
// 	// 	json.Unmarshal(d.Body, &product)
// 	// 	// productRepo.InsertProduct(product)
// 	// }
// 	// }()
// 	return msgs
// }

func main() {
	db := config.InitDB()
	cartRepository := repository.NewCartRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartService := service.NewCartService(cartRepository, productRepo)
	cartController := controller.NewCartController(cartService)

	_, cancel := context.WithCancel(context.Background())

	exchange := "product.created"
	queue := "product.create"
	routingKey := "create"

	amqpConn, err := config.NewRabbitMqConn(exchange, queue, routingKey)
	if err != nil {
		exception.FailOnError(err, "failed connect to rabbit mq")
	}

	cartConsumer := rabbitmq.NewCartConsumer(amqpConn, productRepo)
	// ch, err := cartConsumer.CreateChannel(exchange, queue, routingKey, "")
	// if err != nil {
	// 	exception.FailOnError(err, "create Channel")
	// }
	// defer ch.Close()

	// deliveries,err := ch.Consume(
	// 	queue,
	// 	"",
	// 	true,   // auto-ack
	// 	false,  // exclusive
	// 	false,  // no-local
	// 	false,  // no-wait
	// 	nil,
	// )
	// if err != nil {
	// 	exception.FailOnError(err, "consume")
	// }

	// go func() {

	// }()
	go func() {
		err := cartConsumer.StartConsumer(5, exchange, queue, routingKey, "")
		if err != nil {
			fmt.Printf("StartConsumer: %v", err)
			cancel()
		}
	}()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Range, Authorization",
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("cart service")
	})

	// exchange := "product.created"
	// queue := "product.create"
	// routingKey := "create"

	// ch := config.NewRabbitMqConn(exchange, queue, routingKey)
	// msgs, err := ch.ListenQueue()
	// if err != nil {
	// 	exception.FailOnError(err, "failed listen queue")
	// }
	// go func() {
	// for d := range msgs {
	// 	var product *entity.Product
	// 	json.Unmarshal(d.Body, &product)
	// 	productRepo.InsertProduct(product)
	// }
	// }()
	// msgs := ListenQueue()
	// for d := range msgs {
	// 	var product *entity.Product
	// 	json.Unmarshal(d.Body, &product)
	// 	log.Print(product)
	// 	// productRepo.InsertProduct(product)
	// }
	cartController.Route(app)
	app.Listen(":5004")

}
