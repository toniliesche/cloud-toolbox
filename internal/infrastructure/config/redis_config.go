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
}

func (c *RedisConfig) Validate(path string) error {
	if c.Host == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.host", path))
	}

	if c.Port == 0 {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.port", path))
	}

	if c.Port < 0 {
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

	return cfg, nil
}
