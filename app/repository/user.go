package repository

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"

	"github.com/rezaAmiri123/service-user/app/model"
	"github.com/rezaAmiri123/service-user/pkg/utils"
)

const (
	userByIdCacheDuration = 3600
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	IsFollowing(ctx context.Context, a, b *model.User) (bool, error)
	Follow(ctx context.Context, a, b *model.User) error
	Unfollow(ctx context.Context, a, b *model.User) error
}

type ORMUserRepository struct {
	db        *gorm.DB
	cacheRepo UserCacheRepository
}

func NewORMUserRepository(db *gorm.DB, cacheRepo UserCacheRepository) *ORMUserRepository {
	return &ORMUserRepository{db: db, cacheRepo: cacheRepo}
}

// Create create a user
func (repo *ORMUserRepository) Create(ctx context.Context, user *model.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Create")
	defer span.Finish()

	return repo.db.Create(user).Error
}

//Update update the user
func (repo *ORMUserRepository) Update(ctx context.Context, user *model.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Update")
	defer span.Finish()

	repo.cacheRepo.DeleteByID(ctx, utils.UintToString(user.ID))

	return repo.db.Model(user).Update(user).Error
}

// GetByEmail finds a user from email
func (repo *ORMUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetByEmail")
	defer span.Finish()

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
func (repo *ORMUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetByID")
	defer span.Finish()

	cachedUser, _ := repo.cacheRepo.GetByID(ctx, utils.UintToString(id))
	if cachedUser != nil {
		return cachedUser, nil
	}

	var u model.User
	if err := repo.db.Find(&u, id).Error; err != nil {
		return nil, err
	}

	repo.cacheRepo.SetByID(ctx, utils.UintToString(u.ID), userByIdCacheDuration, &u)

	return &u, nil
}

// GetByUsername finds a user from username
func (repo *ORMUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetByUsername")
	defer span.Finish()

	var u model.User
	if err := repo.db.Where(model.User{Username: username}).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// IsFollowing returns whether user A follows user B or not
func (repo *ORMUserRepository) IsFollowing(ctx context.Context, a, b *model.User) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.IsFollowing")
	defer span.Finish()

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
func (repo *ORMUserRepository) Follow(ctx context.Context, a, b *model.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Follow")
	defer span.Finish()

	return repo.db.Model(a).Association("Follows").Append(b).Error
}

// Unfollow delete follow relationship to user B from user a
func (repo *ORMUserRepository) Unfollow(ctx context.Context, a, b *model.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Unfollow")
	defer span.Finish()

	return repo.db.Model(a).Association("Follows").Delete(b).Error
}
