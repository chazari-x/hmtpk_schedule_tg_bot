package redis

import (
	"errors"
	"fmt"
	"github.com/chazari-x/hmtpk_schedule/config"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type Redis struct {
	Cfg *config.Redis
	r   *redis.Client
}

func NewRedis(cfg *config.Redis) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), // Адрес и порт Redis сервера
		Password: cfg.Pass,                                 // Пароль, если установлен
		DB:       0,                                        // Номер базы данных
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &Redis{r: client, Cfg: cfg}, nil
}

func (c *Redis) AddData(key, value string, minutes float64) error {
	if key == "" {
		return errors.New("key is nil")
	}

	if value == "" {
		return errors.New("value is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	if minutes <= 0 {
		minutes = 5
	}

	return c.r.Set(ctx, key, value, time.Minute*time.Duration(minutes)).Err()
}

func (c *Redis) GetData(key string) (string, error) {
	if key == "" {
		return "", errors.New("key is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return c.r.Get(ctx, key).Result()
}
