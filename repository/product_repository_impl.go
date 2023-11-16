package repository

import (
	"cart-service/entity"
	"log"

	"gorm.io/gorm"
)

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{
		db: db,
	}
}

type ProductRepositoryImpl struct {
	db *gorm.DB
}

// InsertProduct implements ProductRepository.
func (repository *ProductRepositoryImpl) InsertProduct(product *entity.Product) {
	err := repository.db.Create(product).Debug().Error
	if err != nil {
		log.Print(err)
	}
}
