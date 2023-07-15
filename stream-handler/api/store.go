package api

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(c Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
		DB:       0,
	})

	return rdb
}

func ReportStream(s *Stream, r *redis.Client) {
	r.Set(context.Background(), s.Name, s.Manifest, 30*time.Second)
}
