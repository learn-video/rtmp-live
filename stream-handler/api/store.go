package api

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrStreamNotFound = errors.New("stream not found")
)

type Stream struct {
	Name     string `json:"name"`
	Manifest string `json:"manifest"`
	Host     string `json:"host"`
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
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return r.Set(context.Background(), s.Name, data, 30*time.Second).Err()
}

func FetchStream(streamName string, r *redis.Client) (Stream, error) {
	res, err := r.Get(context.Background(), streamName).Result()
	if err == redis.Nil {
		return Stream{}, ErrStreamNotFound
	} else if err != nil {
		return Stream{}, err
	} else {
		var stream Stream
		err := json.Unmarshal([]byte(res), &stream)
		if err != nil {
			return Stream{}, err
		}
		return stream, nil
	}
}
