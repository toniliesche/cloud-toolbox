// MIT License
// Copyright (c) 2025 Toni Liesche
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

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
