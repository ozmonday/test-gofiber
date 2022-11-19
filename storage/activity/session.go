package activity

import (
	"context"
	"errors"
	"fmt"
	"testfiber/storage/entities"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
)

type Session interface {
	Set(context.Context, string, entities.Activity) error
	Get(context.Context, string, *entities.Activity) error
}

type session struct {
	client   *redis.Client
	duration time.Duration
}

func NewSession(conn *redis.Client) Session {
	return &session{
		client:   conn,
		duration: 10 * time.Minute,
	}
}

func (s *session) Set(ctx context.Context, key string, value entities.Activity) error {
	key = fmt.Sprintf("ACT:%s", key)
	data, err := json.Marshal(value)

	if err != nil {
		return err
	}

	_, err = s.client.Set(ctx, key, string(data), s.duration).Result()
	return err

}

func (s *session) Get(ctx context.Context, key string, value *entities.Activity) error {
	key = fmt.Sprintf("ACT:%s", key)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil || data == string(redis.Nil) {
		e := errors.New("thera are someting wrong or data is empty")
		return e
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return err
	}
	return nil
}
