package repository

import (
	"log"

	"cart-service/entity"

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

// FindProductById implements ProductRepository.
func (repository *ProductRepositoryImpl) FindProductById(id string) (*entity.Product, error) {
	var entity *entity.Product
	err := repository.db.First(&entity, "id = ? ", id).Debug().Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// InsertProduct implements ProductRepository.
func (repository *ProductRepositoryImpl) InsertProduct(product *entity.Product) {
	err := repository.db.Create(product).Debug().Error
	if err != nil {
		log.Print(err)
	}
}
