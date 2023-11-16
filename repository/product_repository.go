package repository

import "cart-service/entity"

type ProductRepository interface {
	FindProductById(id string) (*entity.Product, error)
	InsertProduct(product *entity.Product)
}
