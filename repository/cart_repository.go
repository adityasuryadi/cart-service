package repository

import (
	"cart-service/entity"
)

type CartRepository interface {
	Insert(entity *entity.Cart) (*entity.Cart, error)
	GetCartByUserId(userId string) (entities []entity.Cart, err error)
	GetCartByProductId(productId string) (cart *entity.Cart, err error)
	GetCartById(id string) (cart *entity.Cart, err error)
	UpdateCart(cart *entity.Cart) (*entity.Cart, error)
	DeletCart(productId string) error
}
