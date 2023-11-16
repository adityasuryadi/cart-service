package repository

import "cart-service/entity"

type ProductRepository interface {
	InsertProduct(product *entity.Product)
}
