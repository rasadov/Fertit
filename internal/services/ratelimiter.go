package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RateLimiter interface {
	IsBlocked(ctx context.Context, key string) (bool, error)
	RecordFailedAttempt(ctx context.Context, key string) error
	RecordSuccessfulAttempt(ctx context.Context, key string)
	GetAttemptCount(ctx context.Context, key string) (int, error)
	GetTimeUntilReset(ctx context.Context, ip string) (time.Duration, error)
}

type RedisRateLimiter struct {
	client   *redis.Client
	maxTries int
	window   time.Duration
}

func NewRateLimiter(client *redis.Client, maxTries int, window time.Duration) RateLimiter {
	return &RedisRateLimiter{
		client:   client,
		maxTries: maxTries,
		window:   window,
	}
}

func (r *RedisRateLimiter) IsBlocked(ctx context.Context, ip string) (bool, error) {
	key := fmt.Sprintf("login_attempts:%s", ip)

	attempts, err := r.client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return attempts >= r.maxTries, nil
}

func (r *RedisRateLimiter) RecordFailedAttempt(ctx context.Context, ip string) error {
	key := fmt.Sprintf("login_attempts:%s", ip)

	// Increment with expiration
	pipe := r.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, r.window)
	_, err := pipe.Exec(ctx)

	return err
}

func (r *RedisRateLimiter) RecordSuccessfulAttempt(ctx context.Context, ip string) {
	key := fmt.Sprintf("login_attempts:%s", ip)
	r.client.Del(ctx, key)
}

func (r *RedisRateLimiter) GetAttemptCount(ctx context.Context, ip string) (int, error) {
	key := fmt.Sprintf("login_attempts:%s", ip)

	attempts, err := r.client.Get(ctx, key).Int()
	if err != nil {
		return 0, err
	}
	return attempts, nil
}

func (r *RedisRateLimiter) GetTimeUntilReset(ctx context.Context, ip string) (time.Duration, error) {
	key := fmt.Sprintf("login_attempts:%s", ip)

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if ttl < 0 {
		return 0, nil
	}
	return ttl, nil
}
