package repository

import (
	"cart-service/entity"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewCartRepository(db *gorm.DB) CartRepository {
	return &CartRepositoryImpl{
		db: db,
	}
}

type CartRepositoryImpl struct {
	db *gorm.DB
}

// DeletecartsByProductId implements CartRepository.
func (repository *CartRepositoryImpl) DeletecartsByProductId(id string) error {
	var cart *entity.Cart
	err := repository.db.Where("product_id", id).Delete(&cart).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateCartsByProduct implements CartRepository.
func (repository *CartRepositoryImpl) UpdateCartsByProductId(cart *entity.Cart) error {
	err := repository.db.Debug().Where("product_id", cart.ProductId).Updates(&entity.Cart{ProductName: cart.ProductName, ProductPrice: cart.ProductPrice}).Error
	if err != nil {
		return err
	}
	return nil
}

// GetCartById implements CartRepository.
func (repository *CartRepositoryImpl) GetCartById(id string) (cart *entity.Cart, err error) {
	err = repository.db.First(&cart, "id = ?", id).Debug().Error
	if err != nil {
		return nil, err
	}
	fmt.Println("cart", cart)
	return cart, nil
}

// DeletCart implements CartRepository.
func (repository *CartRepositoryImpl) DeletCart(id string) error {
	var entity entity.Cart
	err := repository.db.Where("id = ?", uuid.MustParse(id)).Delete(&entity).Error
	if err != nil {
		return err
	}
	return nil
}

// updateCart implements CartRepository.
func (repository *CartRepositoryImpl) UpdateCart(entity *entity.Cart) (*entity.Cart, error) {
	// var cart entity.Cart
	result := repository.db.Save(entity)

	fmt.Println("update", result.RowsAffected)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected > 0 {
		return entity, nil
	}

	return entity, nil
}

// GetCartByProductId implements CartRepository.
func (repository *CartRepositoryImpl) GetCartByProductId(productId string) (cart *entity.Cart, err error) {
	var entity *entity.Cart
	err = repository.db.Where("product_id", uuid.MustParse(productId)).Find(&entity).Error
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// GetCartByUserId implements CartRepository.
func (repository *CartRepositoryImpl) GetCartByUserId(userId string) (entities []entity.Cart, err error) {
	err = repository.db.Where("user_id = ?", uuid.MustParse(userId)).Debug().Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}

// Insert implements CartRepository.
func (repository *CartRepositoryImpl) Insert(entity *entity.Cart) (*entity.Cart, error) {
	err := repository.db.Create(entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
