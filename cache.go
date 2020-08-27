package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type shrinkflateCache struct {
	rdb      *redis.Client
	host     string
	port     string
	password string
}

func (cache shrinkflateCache) New() (*shrinkflateCache, error) {
	cache.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cache.host, cache.port),
		Password: cache.password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cache.rdb.Ping(ctx).Result()

	if err != nil {
		return &cache, err
	}

	return &cache, nil
}

func (cache shrinkflateCache) Remember(key, value string, expiresInSeconds int64) (string, error) {
	rdb, ctx, cancel := cache.GetRdbWithContext()
	defer cancel()

	return rdb.Set(ctx, key, value, time.Duration(expiresInSeconds)).Result()
}

func (cache shrinkflateCache) Get(key string) (string, error) {
	rdb, ctx, cancel := cache.GetRdbWithContext()
	defer cancel()

	return rdb.Get(ctx, key).Result()
}

func (cache shrinkflateCache) Forget(key string) (int64, error) {
	rdb, ctx, cancel := cache.GetRdbWithContext()
	defer cancel()

	return rdb.Del(ctx, key).Result()
}

func (cache shrinkflateCache) GetRdbWithContext() (*redis.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	return cache.rdb, ctx, cancel
}
