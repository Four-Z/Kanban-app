package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	user := entity.User{}
	err := r.db.First(&entity.User{}, id).Scan(&user).Error

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	user := &entity.User{}
	// err := r.db.First(&entity.User{}, "email = ?", email).Scan(&user).Error
	result := r.db.Where("email = ?", email).Find(&entity.User{}).Scan(&user)

	if result.RowsAffected > 0 {
		return *user, errors.New("email already exists")
	} else {
		return entity.User{}, nil
	}

}

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := r.db.Model(&entity.User{}).Where("id = ?", user.ID).Updates(user).Error

	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	err := r.db.Where("id = ?", id).Delete(&entity.User{}).Error
	if err != nil {
		return err
	}

	return nil
}
