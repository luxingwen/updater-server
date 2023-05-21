package redisop

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	standaloneClient *redis.Client
	clusterClient    *redis.ClusterClient
	isCluster        bool
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	addrs := strings.Split(addr, ",")
	client := &RedisClient{}

	if len(addrs) > 1 {
		client.isCluster = true
		client.clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: password,
		})
	} else {
		client.isCluster = false
		client.standaloneClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
	}

	return client
}

func (c *RedisClient) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	if c.isCluster {
		return c.clusterClient.Set(ctx, key, value, expiration).Err()
	}
	return c.standaloneClient.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	if c.isCluster {
		return c.clusterClient.Del(ctx, key).Err()
	}
	return c.standaloneClient.Del(ctx, key).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	if c.isCluster {
		return c.clusterClient.Get(ctx, key).Result()
	}
	return c.standaloneClient.Get(ctx, key).Result()
}

func (c *RedisClient) Close() error {
	if c.isCluster {
		return c.clusterClient.Close()
	}
	return c.standaloneClient.Close()
}
