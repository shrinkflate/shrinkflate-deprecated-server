package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type shrinkflateQueue struct {
	rdb      *redis.Client
	host     string
	port     int
	password string
}

func (queue shrinkflateQueue) New() (shrinkflateQueue, error) {
	queue.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", queue.host, queue.port),
		Password: queue.password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := queue.rdb.Ping(ctx).Result()

	if err != nil {
		return queue, err
	}

	return queue, nil
}

func (queue shrinkflateQueue) Remember(key, value string, expiresInSeconds int64) (string, error) {
	rdb, ctx, cancel := queue.GetRdbWithContext()
	defer cancel()

	return rdb.Set(ctx, key, value, time.Duration(expiresInSeconds)).Result()
}

func (queue shrinkflateQueue) Get(key string) (string, error) {
	rdb, ctx, cancel := queue.GetRdbWithContext()
	defer cancel()

	return rdb.Get(ctx, key).Result()
}

func (queue shrinkflateQueue) Forget(key string) (int64, error) {
	rdb, ctx, cancel := queue.GetRdbWithContext()
	defer cancel()

	return rdb.Del(ctx, key).Result()
}

func (queue shrinkflateQueue) GetRdbWithContext() (*redis.Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	return queue.rdb, ctx, cancel
}
