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

func TestValidateFunctionAsAServiceConfigFailsOnMissingSystemConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.SystemConfig = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing system section error") {
		return
	}

	if !assert.Equal(t, "config section `system` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnInvalidSystemConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.SystemConfig.Log = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid system config error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `system`", "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnMissingHttpConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.HttpServer = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing http section error") {
		return
	}

	if !assert.Equal(t, "config section `http` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnInvalidHttpConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.HttpServer.Port = 0

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid http config error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `http`", "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnMissingCommand(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.Command = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing command error") {
		return
	}

	if !assert.Equal(t, "config value `command` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnMissingFunctionName(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.FunctionName = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing function name error") {
		return
	}

	if !assert.Equal(t, "config value `function_name` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnZeroParallelExecution(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.ParallelExecution = 0

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch zero parallel execution error") {
		return
	}

	if !assert.Equal(t, "config value `parallel_execution` must be greater than `0`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnLowExecutionTimeout(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.ExecutionTimeout = 5

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch low execution timeout error") {
		return
	}

	if !assert.Equal(t, "config value `execution_timeout` must be greater than or equal to `10`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnHighExecutionTimeout(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.ExecutionTimeout = 1000

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch high execution timeout error") {
		return
	}

	if !assert.Equal(t, "config value `execution_timeout` must be less than or equal to `900`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnNegativeStorageTtl(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.StorageTtl = -1

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch negative storage ttl error") {
		return
	}

	if !assert.Equal(t, "config value `storage_ttl` must be greater than or equal to `0`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnInvalidStorageBackend(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.StorageBackend = "invalid"

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid storage backend error") {
		return
	}

	if !assert.Equal(t, "config value `storage_backend` must be one of [inmemory redis scylla], but is `invalid`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnMissingRedisConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.StorageBackend = "redis"
	cfg.Redis = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing redis config error") {
		return
	}

	if !assert.Equal(t, "config section `redis` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnInvalidRedisConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.StorageBackend = "redis"
	cfg.Redis.Host = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid redis config error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `redis`", "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnMissingScyllaConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.StorageBackend = "scylla"
	cfg.Scylla = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing scylla config error") {
		return
	}

	if !assert.Equal(t, "config section `scylla` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigFailsOnInvalidScyllaConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.StorageBackend = "scylla"
	cfg.Scylla = getValidScyllaConfig()
	cfg.Scylla.Host = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid scylla config error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `scylla`", "unexpected error message") {
		return
	}
}

func TestValidateFunctionAsAServiceConfigSucceedsOnValidConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()

	err := cfg.Validate()
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func getValidFunctionAsAServiceConfig() *config.FunctionAsAServiceConfig {
	return &config.FunctionAsAServiceConfig{
		ApplicationConfig: config.ApplicationConfig{
			SystemConfig: getValidSystemConfig(),
		},
		HttpServer:        getValidHttpServerConfig(),
		Command:           "echo 'Hello, World!'",
		FunctionName:      "hello_world",
		ParallelExecution: 1,
		ExecutionTimeout:  30,
		StorageBackend:    "redis",
		StorageTtl:        10,
		Redis:             getValidRedisConfig(),
	}
}
