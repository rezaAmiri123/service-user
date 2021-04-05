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
	GetByUsername(username string) (*model.User, error)
	IsFollowing(a, b *model.User) (bool, error)
	Follow(a,b *model.User)error
	Unfollow(a,b *model.User)error
}

type ORMUserRepository struct {
	db *gorm.DB
}

func NewORMUserRepository(db *gorm.DB) *ORMUserRepository {
	return &ORMUserRepository{db: db}
}

// Create create a user
func (repo *ORMUserRepository) Create(user *model.User) error {
	return repo.db.Create(user).Error
}

//Update update the user
func (repo *ORMUserRepository) Update(user *model.User) error {
	return repo.db.Model(user).Update(user).Error
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

// GetByUsername finds a user from username
func (repo *ORMUserRepository) GetByUsername(username string) (*model.User, error) {
	var u model.User
	if err := repo.db.Where(model.User{Username: username}).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// IsFollowing returns whether user A follows user B or not
func (repo *ORMUserRepository) IsFollowing(a, b *model.User) (bool, error) {
	if a == nil || b == nil {
		return false, nil
	}
	var count int
	err := repo.db.Table("follows").
		Where("from_user_id = ? AND to_user_id = ?", a.ID, b.ID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Follow create follow relationship to user B from user A
func (repo *ORMUserRepository)Follow(a,b *model.User)error{
	return repo.db.Model(a).Association("Follows").Append(b).Error
}

// Unfollow delete follow relationship to user B from user a
func (repo *ORMUserRepository)Unfollow(a,b *model.User)error{
	return repo.db.Model(a).Association("Follows").Delete(b).Error
}
