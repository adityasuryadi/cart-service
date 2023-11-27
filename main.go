package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cart-service/config"
	"cart-service/controller"
	"cart-service/entity"
	"cart-service/exception"
	"cart-service/interface/rabbitmq"
	"cart-service/model"
	"cart-service/repository"
	"cart-service/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db := config.InitDB()
	cartRepository := repository.NewCartRepository(db)
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepository, productRepo)
	cartController := controller.NewCartController(cartService)

	_, cancel := context.WithCancel(context.Background())

	amqpConn, err := config.NewRabbitMqConn()
	if err != nil {
		exception.FailOnError(err, "failed connect to rabbit mq")
	}

	cartConsumer := rabbitmq.NewCartConsumer(amqpConn, productRepo)
	go func() {
		// listen create product
		exchange := "product.created"
		queue := "product.create"
		routingKey := "create"
		params := model.RabbitMQConusmerParams{
			WorkerPoolSize: 5,
			Exchange:       exchange,
			QueueName:      queue,
			BindingKey:     routingKey,
			ConsumerTag:    "",
		}
		err := cartConsumer.StartConsumer(params, func(b []byte) error {
			var product *entity.Product
			json.Unmarshal(b, &product)
			log.Print(product)
			err := productService.CreateProduct(b)
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			fmt.Printf("StartConsumer: %v", err)
			cancel()
		}

	}()

	go func() {
		// listen delete product
		exchange := "product.deleted"
		queue := "product.delete"
		routingKey := "delete"
		paramsDelete := model.RabbitMQConusmerParams{
			WorkerPoolSize: 5,
			Exchange:       exchange,
			QueueName:      queue,
			BindingKey:     routingKey,
			ConsumerTag:    "",
		}
		err = cartConsumer.StartConsumer(paramsDelete, func(b []byte) error {
			var id string
			json.Unmarshal(b, &id)
			fmt.Println("id", id)
			err := productService.DeleteProduct(id)
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			fmt.Printf("StartConsumer: %v", err)
			cancel()
		}

	}()

	go func() {
		// listen delete product
		exchange := "product.updated"
		queue := "product.update"
		routingKey := "update"

		paramsDelete := model.RabbitMQConusmerParams{
			WorkerPoolSize: 5,
			Exchange:       exchange,
			QueueName:      queue,
			BindingKey:     routingKey,
			ConsumerTag:    "",
		}
		err = cartConsumer.StartConsumer(paramsDelete, func(b []byte) error {
			var product *entity.Product
			json.Unmarshal(b, &product)
			err := productService.EditProduct(product)
			if err != nil {
				return err
			}
			return nil
		})

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

	cartController.Route(app)
	app.Listen(":5004")

}
