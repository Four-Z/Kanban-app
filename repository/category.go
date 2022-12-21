package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	categories := []entity.Category{}
	err := r.db.Raw("SELECT * FROM categories WHERE user_id = ?", id).Scan(&categories).Error
	// err := r.db.First(&entity.Category{}, "user_id = ?", id).Scan(&categories).Error

	if err != nil {
		return []entity.Category{}, err
	}

	return categories, nil
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	storedCategory := entity.Category{}
	err = r.db.Create(&category).Scan(&storedCategory).Error

	if err != nil {
		return 0, err
	}

	return storedCategory.ID, nil
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	err := r.db.Create(&categories).Error
	return err
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	category := entity.Category{}
	err := r.db.First(&entity.Category{}, id).Scan(&category).Error

	if err != nil {
		return entity.Category{}, err
	}

	return category, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	return r.db.Model(&entity.Category{}).Where("id = ?", category.ID).Updates(category).Error
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	// return r.db.Where("id = ?", id).Delete(&entity.Category{}).Error
	return r.db.Delete(&entity.Category{}, id).Error
}
