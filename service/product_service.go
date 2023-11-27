package service

import "cart-service/entity"

type ProductService interface {
	CreateProduct(body []byte) error
	EditProduct(entity *entity.Product) error
	DeleteProduct(id string) error
}
