package model

import (
	"github.com/gofiber/fiber/v2"
)

type WebResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetResponse(responseCode int, data interface{}, message string) WebResponse {
	var Code int
	var Status string
	switch responseCode {
	case 404:
		Code = fiber.StatusNotFound
		Status = "NOT_FOUND"
	case 200:
		Code = fiber.StatusOK
		Status = "OK"
	case 201:
		Code = fiber.StatusCreated
		Status = "CREATED"
	case 400:
		Code = fiber.StatusBadRequest
		Status = "BAD_REQUEST"
	case 401:
		Code = fiber.StatusUnauthorized
		Status = "UNAUTHORIZED"
	case 500:
		Code = fiber.StatusInternalServerError
		Status = "INTERNAL_SERVER_ERROR"
	}

	return WebResponse{
		Code:    Code,
		Status:  Status,
		Message: message,
		Data:    data,
	}
}
