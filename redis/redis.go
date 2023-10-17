package redis

import (
	"github.com/chazari-x/hmtpk_schedule/config"
	"github.com/chazari-x/hmtpk_schedule/redis/redis"
)

func Redis(cfg *config.Redis) (*redis.Redis, error) {
	newRedis, err := redis.NewRedis(cfg)
	if err != nil {
		return nil, err
	}

	return newRedis, nil
}
