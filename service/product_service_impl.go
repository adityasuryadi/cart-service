package service

import (
	"cart-service/entity"
	"cart-service/repository"
	"encoding/json"
)

type ProductServiceImpl struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return &ProductServiceImpl{
		repository: repository,
	}
}

// DeleteProduct implements ProductService.
func (service *ProductServiceImpl) DeleteProduct(id string) error {
	service.repository.Delete(id)
	return nil
}

// EditProduct implements ProductService.
func (service *ProductServiceImpl) EditProduct(entity *entity.Product) error {
	err := service.repository.Update(entity)
	if err != nil {
		return err
	}
	return nil
}

// CreateProduct implements ProductService.
func (service *ProductServiceImpl) CreateProduct(body []byte) error {
	var product *entity.Product
	json.Unmarshal(body, &product)
	err := service.repository.InsertProduct(product)
	if err != nil {
		return err
	}
	return nil
}
