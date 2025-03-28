package connectors

import (
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewRedis(container *di.Container) (*redis.Client, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError("Redis")
	}

	if container.RedisConfig == nil {
		return nil, errors.NewResolveDependencyError("Redis", "RedisConfig")
	}

	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", container.RedisConfig.Host, container.RedisConfig.Port),
		DB:   int(container.RedisConfig.Database),
	}), nil
}
