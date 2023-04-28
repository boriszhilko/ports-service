package redis

import (
	"context"
	"encoding/json"
	"github.com/boriszhilko/ports-service/internal/port/domain"
	"github.com/go-redis/redis/v8"
)

type Repository struct {
	client *redis.Client
}

func NewRedisRepository(url string) *Repository {
	return &Repository{
		client: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: "",
			DB:       0,
		})}
}

func (mr *Repository) CreateOrUpdate(ctx context.Context, ports []domain.Port) error {
	err := mr.writeBatch(ctx, ports)
	if err != nil {
		return err
	}
	return nil
}

func (mr *Repository) Get(ctx context.Context, id string) (domain.Port, error) {
	var port domain.Port

	result, err := mr.client.Get(ctx, id).Result()
	if err != nil {
		if err == redis.Nil {
			return port, domain.ErrPortNotFound
		} else {
			return port, err
		}
	}

	err = json.Unmarshal([]byte(result), &port)
	if err != nil {
		return port, err
	}
	return port, nil
}

func (mr *Repository) writeBatch(ctx context.Context, ports []domain.Port) error {
	pipe := mr.client.Pipeline()
	for _, record := range ports {
		key := record.ID

		value, err := json.Marshal(record)
		if err != nil {
			return err
		}

		pipe.Set(ctx, key, value, 0)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (mr *Repository) Stop() error {
	return mr.client.Close()
}
