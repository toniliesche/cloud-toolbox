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

package config

import (
	"cloud-toolbox/internal/infrastructure/errors"
	"fmt"
)

const (
	redisServerDefaultHost     = "localhost"
	redisServerDefaultPort     = 6379
	redisServerDefaultDatabase = 0
)

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	Database int64  `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (c *RedisConfig) Validate(path string) error {
	if c.Host == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.host", path))
	}

	if c.Port < 1 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError(fmt.Sprintf("%s.port", path))
	}

	if c.Database < 0 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError(fmt.Sprintf("%s.database", path))
	}

	return nil
}

func getDefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     redisServerDefaultHost,
		Port:     redisServerDefaultPort,
		Database: redisServerDefaultDatabase,
	}
}

func getRedisConfigFromEnvironment() (*RedisConfig, errors.ApplicationError) {
	cfg := getDefaultRedisConfig()

	host := GetEnvironmentString("REDIS_HOST", redisServerDefaultHost)
	if host != "" {
		cfg.Host = host
	}

	port, err := GetEnvironmentInt("REDIS_PORT", redisServerDefaultPort)
	if err != nil {
		return nil, err
	}

	if port != 0 {
		cfg.Port = port
	}

	database, err := GetEnvironmentInt("REDIS_DATABASE", redisServerDefaultDatabase)
	if err != nil {
		return nil, err
	}

	if database != 0 {
		cfg.Database = database
	}

	username := GetEnvironmentString("REDIS_USERNAME", "")
	if username != "" {
		cfg.Username = username
	}

	password := GetEnvironmentString("REDIS_PASSWORD", "")
	if password != "" {
		cfg.Password = password
	}

	return cfg, nil
}
