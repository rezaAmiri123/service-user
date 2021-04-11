package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"

	"github.com/rezaAmiri123/service-user/app/model"
	"github.com/rezaAmiri123/service-user/pkg/grpc_errors"
	"github.com/rezaAmiri123/service-user/pkg/logger"
)

// UserCacheRepository
type UserCacheRepository interface {
	GetByID(ctx context.Context, key string) (*model.User, error)
	SetByID(ctx context.Context, key string, seconds int, user *model.User) error
	DeleteByID(ctx context.Context, key string) error
}

type userRedisRepo struct {
	redisClient *redis.Client
	basePrefix  string
	logger      logger.Logger
}

func NewUserRedisRepo(redisClient *redis.Client, basePrefix string, logger logger.Logger) *userRedisRepo {
	return &userRedisRepo{redisClient: redisClient, basePrefix: basePrefix, logger: logger}
}

func (r *userRedisRepo) GetByID(ctx context.Context, key string) (*model.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRedisRepo.GetByID")
	defer span.Finish()

	userBytes, err := r.redisClient.Get(ctx, r.createKey(key)).Bytes()
	if err != nil {
		if err != redis.Nil {
			return nil, grpc_errors.ErrNotFound
		}
		return nil, err
	}
	user := &model.User{}
	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRedisRepo) SetByID(ctx context.Context, key string, seconds int, user *model.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRedisRepo.SetByID")
	defer span.Finish()

	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.redisClient.Set(ctx, r.createKey(key), userBytes, time.Second*time.Duration(seconds)).Err()
}

func (r *userRedisRepo) DeleteByID(ctx context.Context, key string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRedisRepo.DeleteByID")
	defer span.Finish()

	return r.redisClient.Del(ctx, r.createKey(key)).Err()
}

func (r *userRedisRepo) createKey(value string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, value)
}
