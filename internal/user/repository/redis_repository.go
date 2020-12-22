package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"time"
)

// Auth redis repository
type userRedisRepo struct {
	redisClient *redis.Client
	basePrefix  string
	logger      logger.Logger
}

// Auth redis repository constructor
func NewUserRedisRepo(redisClient *redis.Client, logger logger.Logger) *userRedisRepo {
	return &userRedisRepo{redisClient: redisClient, basePrefix: "user:", logger: logger}
}

// Get user by id
func (r *userRedisRepo) GetByIDCtx(ctx context.Context, key string) *models.User {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRedisRepo.GetByIDCtx")
	defer span.Finish()

	userBytes, err := r.redisClient.Get(ctx, r.createKey(key)).Bytes()
	if err != nil {
		return nil
	}
	user := &models.User{}
	if err = json.Unmarshal(userBytes, user); err != nil {
		r.logger.Errorf("json.Unmarshal: %v", err)
		return nil
	}
	return user
}

// Cache user with duration in seconds
func (r *userRedisRepo) SetUserCtx(ctx context.Context, key string, seconds int, user *models.User) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRedisRepo.SetUserCtx")
	defer span.Finish()

	userBytes, err := json.Marshal(user)
	if err != nil {
		r.logger.Errorf("json.Unmarshal: %v", err)
		return
	}
	if err = r.redisClient.Set(ctx, r.createKey(key), userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		r.logger.Errorf("userRedisRepo.SetUserCtx.redisClient.Set: %v", err)
	}
}

// Delete user by key
func (r *userRedisRepo) DeleteUserCtx(ctx context.Context, key string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userRedisRepo.DeleteUserCtx")
	defer span.Finish()

	if err := r.redisClient.Del(ctx, r.createKey(key)).Err(); err != nil {
		r.logger.Errorf("userRedisRepo.DeleteUserCtx.redisClient.Del: %v", err)
	}
}

func (r *userRedisRepo) createKey(value string) string {
	return fmt.Sprintf("%s: %s", r.basePrefix, value)
}
