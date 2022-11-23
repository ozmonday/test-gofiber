package todo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testfiber/pkg/entities"
	"time"

	"github.com/go-redis/redis/v8"
)

type Session interface {
	Get(context.Context, string) (*entities.Todo, error)
	Set(context.Context, string, entities.Todo) error
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

func (s *session) Get(ctx context.Context, key string) (*entities.Todo, error) {
	value := new(entities.Todo)
	key = fmt.Sprintf("TODO:%s", key)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil || data == "" {
		e := errors.New("thera are someting wrong or data is empty")
		return nil, e
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return nil, err
	}
	return value, nil
}

func (s *session) Set(ctx context.Context, key string, value entities.Todo) error {
	key = fmt.Sprintf("TODO:%s", key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.client.Set(ctx, key, string(data), s.duration).Result()
	return err
}
