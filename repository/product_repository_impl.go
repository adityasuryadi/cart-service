package repository

import (
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

// Delete implements ProductRepository.
func (repository *ProductRepositoryImpl) Delete(id string) error {
	var product *entity.Product
	err := repository.db.Where("id = ?", id).Delete(&product).Error
	if err != nil {
		return err
	}
	return nil
}

// Update implements ProductRepository.
func (repository *ProductRepositoryImpl) Update(product *entity.Product) error {
	err := repository.db.Save(&product).Error
	if err != nil {
		return err
	}
	return nil
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
func (repository *ProductRepositoryImpl) InsertProduct(product *entity.Product) error {
	err := repository.db.Create(product).Debug().Error
	if err != nil {
		return err
	}
	return nil
}
