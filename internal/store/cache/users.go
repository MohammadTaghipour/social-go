package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MohammadTaghipour/social/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {

	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, cacheKey, json, time.Minute).Err()
}

func (s *UserStore) Get(ctx context.Context, user_id int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", user_id)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}
	return &user, nil
}
