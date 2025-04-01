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

package config_test

import (
	"cloud-toolbox/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateRedisConfigFailsOnEmptyHost(t *testing.T) {
	t.Parallel()
	cfg := getValidRedisConfig()
	cfg.Host = ""

	err := cfg.Validate("redis")
	if assert.Error(t, err, "expected error") {
		assert.Equal(t, "config value `redis.host` must exist", err.Error())
	}
}

func TestValidateRedisConfigFailsOnInvalidPort(t *testing.T) {
	t.Parallel()
	cfg := getValidRedisConfig()
	cfg.Port = -1

	err := cfg.Validate("redis")
	if assert.Error(t, err, "expected error") {
		assert.Equal(t, "config value `redis.port` must be greater than `0`", err.Error())
	}
}

func TestValidateRedisConfigFailsOnInvalidDatabase(t *testing.T) {
	t.Parallel()
	cfg := getValidRedisConfig()
	cfg.Database = -1

	err := cfg.Validate("redis")
	if assert.Error(t, err, "expected error") {
		assert.Equal(t, "config value `redis.database` must be greater than or equal to `0`", err.Error())
	}
}

func TestValidateRedisConfigSucceedsOnValidConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidRedisConfig()

	err := cfg.Validate("redis")
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func getValidRedisConfig() *config.RedisConfig {
	return &config.RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Database: 0,
	}
}
