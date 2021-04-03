package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/rezaAmiri123/service-user/app/model"
)

type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id uint) (*model.User, error)
}

type ORMUserRepository struct {
	db *gorm.DB
}

// Create create a user
func (repo *ORMUserRepository) Create(user *model.User) error {
	return repo.db.Create(user).Error
}

//Update update the user
func (repo *ORMUserRepository) Update(user *model.User) error {
	return repo.db.Model(model.User{}).Update(user).Error
}

// GetByEmail finds a user from email
func (repo *ORMUserRepository) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := repo.db.Where(model.User{
		Email: email,
	}).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByID finds a user from id
func (repo *ORMUserRepository) GetByID(id uint) (*model.User, error) {
	var u model.User
	if err := repo.db.Find(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
