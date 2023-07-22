package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrStreamNotFound = errors.New("stream not found")
)

type Stream struct {
	Name     string
	Manifest string
	Host     string
}

func (s *Stream) Path() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Manifest)
}

func NewRedis(c Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisAddr,
		Password: c.RedisPassword,
		DB:       0,
	})

	return rdb
}

func ReportStream(s *Stream, r *redis.Client) error {
	return r.Set(context.Background(), s.Name, s.Path(), 30*time.Second).Err()
}

func FetchStream(streamName string, r *redis.Client) (string, error) {
	res, err := r.Get(context.Background(), streamName).Result()
	if err == redis.Nil {
		return "", ErrStreamNotFound
	} else if err != nil {
		return "", err
	} else {
		return res, nil
	}
}
