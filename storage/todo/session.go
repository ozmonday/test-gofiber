package todo

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Session interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, interface{}) error
}

type session struct {
	client   *redis.Client
	duration time.Duration
}

func NewSession(coon *redis.Client) Session {
	return &session{
		client:   coon,
		duration: 10 * time.Minute,
	}
}

func (s *session) Get(ctx context.Context, key string) (string, error) {
	key = fmt.Sprintf("TODO:%s", key)
	return s.client.Get(ctx, key).Result()
}

func (s *session) Set(ctx context.Context, key string, value interface{}) error {
	key = fmt.Sprintf("TODO:%s", key)
	_, err := s.client.Set(ctx, key, value, s.duration).Result()
	return err
}
