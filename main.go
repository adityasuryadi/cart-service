package main

import (
	"encoding/json"
	"log"

	"cart-service/config"
	"cart-service/controller"
	"cart-service/entity"
	"cart-service/repository"
	"cart-service/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	db := config.InitDB()
	cartRepository := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepository)
	cartController := controller.NewCartController(cartService)

	productRepo := repository.NewProductRepository(db)

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

	exchange := "product.created"
	queue := "product.create"
	routingKey := "create"

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq_container:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
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

	go func() {
		for d := range msgs {
			// log.Printf(" [x] %s", d.Body)
			var product *entity.Product
			json.Unmarshal(d.Body, &product)

			productRepo.InsertProduct(product)
		}
	}()

	cartController.Route(app)
	app.Listen(":5004")
}
