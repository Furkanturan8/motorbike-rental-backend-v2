package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

var defaultCache *RedisCache

func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	logger.Info("Redis bağlantısı başlatılıyor: %s", addr)

	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Password:        password,
		DB:              db,
		MaxRetries:      5,
		MaxRetryBackoff: 5 * time.Second,
		MinRetryBackoff: time.Second * 2,
		DialTimeout:     time.Second * 5,
	})

	// Bağlantıyı test et
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Redis ping hatası: %v", err)
		return nil, fmt.Errorf("redis ping hatası: %v", err)
	}

	return &RedisCache{client: client}, nil
}

func InitDefaultCache(addr, password string, db int) error {
	cache, err := NewRedisCache(addr, password, db)
	if err != nil {
		return err
	}
	defaultCache = cache
	return nil
}

// Veriyi JSON olarak cache'e yazar
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, json, expiration).Err()
}

// Cache'den veriyi okur ve verilen struct'a unmarshal eder
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Cache'den veriyi siler
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Birden fazla key'i siler
func (c *RedisCache) DeleteMany(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}
	return nil
}

// Key'in var olup olmadığını kontrol eder
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// Key'in süresini günceller
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// Global fonksiyonlar
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if defaultCache == nil {
		return errorx.ErrCacheNotInitialized
	}
	return defaultCache.Set(ctx, key, value, expiration)
}

func Get(ctx context.Context, key string, dest interface{}) error {
	if defaultCache == nil {
		return errorx.ErrCacheNotInitialized
	}
	return defaultCache.Get(ctx, key, dest)
}

func Delete(ctx context.Context, key string) error {
	if defaultCache == nil {
		return errorx.ErrCacheNotInitialized
	}
	return defaultCache.Delete(ctx, key)
}

func DeleteMany(ctx context.Context, pattern string) error {
	if defaultCache == nil {
		return errorx.ErrCacheNotInitialized
	}
	return defaultCache.DeleteMany(ctx, pattern)
}

func Exists(ctx context.Context, key string) (bool, error) {
	if defaultCache == nil {
		return false, errorx.ErrCacheNotInitialized
	}
	return defaultCache.Exists(ctx, key)
}

func Expire(ctx context.Context, key string, expiration time.Duration) error {
	if defaultCache == nil {
		return errorx.ErrCacheNotInitialized
	}
	return defaultCache.Expire(ctx, key, expiration)
}
