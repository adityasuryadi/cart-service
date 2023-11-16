package controller

import (
	"cart-service/model"
	"cart-service/service"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func NewCartController(service service.CartService) CartController {
	return CartController{
		Service: service,
	}
}

type CartController struct {
	Service service.CartService
}

func (controller *CartController) Route(app *fiber.App) {
	app.Post("cart", controller.AddToCart)
	app.Get("carts", controller.GetCarts)
	app.Delete("cart/:id", controller.RemoveCart)
	app.Put("cart/:id/decrement", controller.DecrementQty)
	app.Put("cart/:id/increment", controller.IncrementQty)
}

func (controller *CartController) GetCarts(ctx *fiber.Ctx) error {
	userId := ctx.Get("X-Kong-Jwt-Claim-User_id")
	carts, err := controller.Service.GetUserCart(userId)
	if userId == "" {
		response := model.GetResponse(401, errors.New("UnAuthorize"), "Unauthorize")
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if err != nil {
		response := model.GetResponse(500, err.Error(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}
	response := model.GetResponse(200, carts, "success")
	fmt.Println(response)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (controller *CartController) AddToCart(ctx *fiber.Ctx) error {
	var request model.InsertCartRequest
	ctx.BodyParser(&request)
	request.UserId = ctx.Get("X-Kong-Jwt-Claim-User_id")
	request.Email = ctx.Get("X-Kong-Jwt-Claim-Email")
	responseCode, cart, err := controller.Service.AddToCart(&request)
	// response := model.GetResponse(responseCode, err.Error(), "")

	if responseCode == 400 {
		response := model.GetResponse(responseCode, err.Error(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	// if responseCode == 500 {
	// 	response := model.GetResponse(responseCode, err.Error(), err.Error())
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	// }

	if responseCode == 200 {
		response := model.GetResponse(200, cart, "")
		return ctx.Status(fiber.StatusOK).JSON(response)
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(nil)
}

func (controller *CartController) RemoveCart(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	responseCode, err := controller.Service.DestroyCart(id)
	if responseCode == 500 {
		response := model.GetResponse(responseCode, err.Error(), "")
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}
	response := model.GetResponse(responseCode, "", "Success Remove")
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (controller *CartController) IncrementQty(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	responseCode, err := controller.Service.IncrementQty(id)
	if responseCode == 500 {
		response := model.GetResponse(responseCode, err.Error(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if responseCode == 404 {
		response := model.GetResponse(responseCode, err.Error(), err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(response)
	}

	if responseCode == 400 {
		response := model.GetResponse(responseCode, err.Error(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := model.GetResponse(responseCode, "", "success update")
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (controller *CartController) DecrementQty(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	responseCode, err := controller.Service.DecrementQty(id)
	if responseCode == 500 {
		response := model.GetResponse(responseCode, err.Error(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if responseCode == 404 {
		response := model.GetResponse(responseCode, err.Error(), err.Error())
		return ctx.Status(fiber.StatusNotFound).JSON(response)
	}

	if responseCode == 400 {
		response := model.GetResponse(responseCode, err.Error(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(response)
	}

	response := model.GetResponse(responseCode, "", "success update")
	return ctx.Status(fiber.StatusOK).JSON(response)
}
