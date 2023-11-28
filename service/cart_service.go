package service

import "cart-service/model"

type CartService interface {
	AddToCart(request *model.InsertCartRequest) (responseCode int, cartResponse *model.CartResponse, err error)
	GetUserCart(userId string) (carts []model.CartResponse, err error)
	DestroyCart(productId string) (responseCode int, err error)
	DecrementQty(productId string) (responseCode int, err error)
	IncrementQty(productId string) (responseCode int, err error)
	EditCartByProductId(request *model.UpdateCartRequest) error
	DeleteCartByProductId(id string) error
}
