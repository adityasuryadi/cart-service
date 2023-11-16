package service

import (
	"errors"
	"fmt"

	"cart-service/entity"
	"cart-service/model"
	"cart-service/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewCartService(repository repository.CartRepository, productRepository repository.ProductRepository) CartService {
	return &CartServiceImpl{
		repository:  repository,
		ProductRepo: productRepository,
	}
}

type CartServiceImpl struct {
	repository  repository.CartRepository
	ProductRepo repository.ProductRepository
}

// DecrementQty implements CartService.
func (service *CartServiceImpl) DecrementQty(id string) (responseCode int, err error) {
	cart, err := service.repository.GetCartById(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, errors.New("product not found")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return 500, err
	}

	if cart.Qty <= 1 {
		return 400, errors.New("qty min 1")
	}

	currentQty := cart.Qty - 1
	cart.TotalPrice = cart.ProductPrice * float64(currentQty)
	cart.Qty = currentQty
	service.repository.UpdateCart(cart)
	return 200, nil
}

// IncrementQty implements CartService.
func (service *CartServiceImpl) IncrementQty(id string) (responseCode int, err error) {
	cart, err := service.repository.GetCartById(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 404, errors.New("product not found")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return 500, err
	}
	fmt.Println(cart)
	if cart.Qty >= 5 {
		return 400, errors.New("qty max 5")
	}

	currentQty := cart.Qty + 1
	cart.TotalPrice = cart.ProductPrice * float64(currentQty)
	cart.Qty = currentQty
	service.repository.UpdateCart(cart)
	return 200, nil
}

// DestroyCart implements CartService.
func (service *CartServiceImpl) DestroyCart(id string) (responseCode int, err error) {
	err = service.repository.DeletCart(id)
	if err != nil {
		return 500, err
	}
	return 200, nil
}

// GetUserCart implements CartService.
func (service *CartServiceImpl) GetUserCart(userId string) (carts []model.CartResponse, err error) {
	enitities, err := service.repository.GetCartByUserId(userId)
	if err != nil {
		return nil, err
	}
	for _, v := range enitities {
		carts = append(carts, model.CartResponse{
			Id:           v.Id.String(),
			ProductName:  v.ProductName,
			ProductPrice: v.ProductPrice,
			ProductId:    v.ProductId.String(),
			Qty:          v.Qty,
			TotalPrice:   v.TotalPrice,
		})
	}
	if len(carts) == 0 {
		carts = []model.CartResponse{}
	}
	return carts, nil
}

// AddToCart implements CartService.
func (service *CartServiceImpl) AddToCart(request *model.InsertCartRequest) (responseCode int, responseCart *model.CartResponse, err error) {
	// get product api
	// url := "http://product_service:5002/product/" + request.ProductId
	// client := resty.New()
	// client.SetTimeout(1 * time.Minute)
	// resp, err := client.R().Get(url)

	// if resp.IsError() {
	// 	return 500, nil, err
	// }

	// response := make(map[string]interface{})
	// json.Unmarshal(resp.Body(), &response)
	// data := response["data"].(map[string]interface{})

	// // check validation stock
	// if data["qty"].(float64) < 1 {
	// 	return 400, nil, errors.New("out of stock")
	// }

	product, err := service.ProductRepo.FindProductById(request.ProductId)
	if err != nil {
		return 500, nil, errors.New(err.Error())
	}

	if product == nil {
		return 400, nil, errors.New("product not found")
	}

	if product.Qty < 1 {
		return 400, nil, errors.New("out of stock")
	}

	// check if product exist on cart,increment if product exist
	cart, err := service.repository.GetCartByProductId(request.ProductId)

	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return 500, nil, err
	}

	if cart != nil && err == nil {
		currentQty := cart.Qty + 1
		cart.TotalPrice = cart.ProductPrice * float64(currentQty)
		cart.Qty = currentQty
		service.repository.UpdateCart(cart)

		responseCart = &model.CartResponse{
			Id:           cart.Id.String(),
			ProductName:  cart.ProductName,
			ProductPrice: cart.ProductPrice,
			Qty:          cart.Qty,
			TotalPrice:   cart.TotalPrice,
		}

		return 200, responseCart, nil
	}

	cart = &entity.Cart{
		ProductId:    uuid.MustParse(request.ProductId),
		Email:        request.Email,
		ProductName:  product.Name,
		Qty:          1,
		ProductPrice: product.Price,
		UserId:       uuid.MustParse(request.UserId),
		TotalPrice:   product.Price * 1,
	}
	entityCart, err := service.repository.Insert(cart)
	if err != nil {
		return 500, nil, err
	}

	responseCart = &model.CartResponse{
		Id:           entityCart.Id.String(),
		ProductName:  entityCart.ProductName,
		ProductPrice: entityCart.ProductPrice,
		ProductId:    entityCart.ProductId.String(),
		Qty:          entityCart.Qty,
		TotalPrice:   entityCart.TotalPrice,
	}
	return 200, responseCart, nil
}
