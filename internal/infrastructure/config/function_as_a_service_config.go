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
	"gopkg.in/yaml.v3"
	"os"
	"slices"
)

type FunctionAsAServiceConfig struct {
	ApplicationConfig
	HttpServer        *HttpServerConfig `yaml:"http"`
	Redis             *RedisConfig      `yaml:"redis"`
	Scylla            *ScyllaConfig     `yaml:"scylla"`
	StorageBackend    string            `json:"storage_backend"`
	StorageTtl        int64             `json:"storage_ttl"`
	Command           string            `yaml:"command"`
	ParallelExecution int64             `yaml:"parallel_execution"`
	ExecutionTimeout  int64             `yaml:"execution_timeout"`
}

func (c *FunctionAsAServiceConfig) Validate() errors.ApplicationError {
	if c.SystemConfig == nil {
		return errors.NewMissingConfigSectionError("system")
	}

	if err := c.SystemConfig.Validate("system"); err != nil {
		return errors.NewValidateConfigSectionError("system", err)
	}

	if c.HttpServer == nil {
		return errors.NewMissingConfigSectionError("http")
	}

	if err := c.HttpServer.Validate("html"); err != nil {
		return errors.NewValidateConfigSectionError("http", err)
	}

	if c.Command == "" {
		return errors.NewMissingConfigValueError("command")
	}

	if c.ParallelExecution < 1 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError("parallel_execution")
	}

	if c.ExecutionTimeout < 10 {
		return errors.NewConfigValueNeedsToBeGreaterThanOrEqualValueError("execution_timeout", 10)
	}

	if c.ExecutionTimeout > 900 {
		return errors.NewConfigValueNeedsToBeLessThanOrEqualValueError("execution_timeout", 900)
	}

	if c.StorageTtl < 0 {
		return errors.NewConfigValueNeedsToBeGreaterThanOrEqualValueError("storage_ttl", 0)
	}

	validBackends := []string{"inmemory", "redis", "scylla"}
	if !slices.Contains(validBackends, c.StorageBackend) {
		return errors.NewInvalidConfigValueError("storage_backend", validBackends, c.StorageBackend)
	}

	if c.StorageBackend == "redis" {
		if c.Redis == nil {
			return errors.NewMissingConfigSectionError("redis")
		}

		if err := c.Redis.Validate("redis"); err != nil {
			return errors.NewValidateConfigSectionError("redis", err)
		}
	} else {
		c.Redis = nil
	}

	if c.StorageBackend == "scylla" {
		if c.Scylla == nil {
			return errors.NewMissingConfigSectionError("scylla")
		}

		if err := c.Scylla.Validate("scylla"); err != nil {
			return errors.NewValidateConfigSectionError("scylla", err)
		}
	} else {
		c.Scylla = nil
	}

	return nil
}

func ProvideFunctionAsAServiceConfig() (*FunctionAsAServiceConfig, errors.ApplicationError) {
	configFile := GetEnvironmentString("FUNCTION_AS_A_SERVICE_CONFIG_FILE", "")

	var cfg *FunctionAsAServiceConfig
	var err errors.ApplicationError

	if configFile != "" {
		cfg, err = getFunctionAsAServiceConfigFromFile(configFile)
	} else {
		cfg, err = getFunctionAsAServiceConfigFromEnvironment()
	}

	if err != nil {
		return nil, errors.NewConfigCreationError(err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, errors.NewConfigValidationError(err)
	}

	return cfg, nil
}

func getDefaultFunctionAsAServiceConfig() *FunctionAsAServiceConfig {
	return &FunctionAsAServiceConfig{
		HttpServer:     getDefaultHttpServerConfig(),
		Redis:          getDefaultRedisConfig(),
		Scylla:         getDefaultScyllaConfig(),
		StorageBackend: "inmemory",
		StorageTtl:     60,
	}
}

func getFunctionAsAServiceConfigFromFile(file string) (*FunctionAsAServiceConfig, errors.ApplicationError) {
	fileContents, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.NewConfigFileReadingError(err)
	}

	cfg := getDefaultFunctionAsAServiceConfig()
	if err = yaml.Unmarshal(fileContents, cfg); err != nil {
		return nil, errors.NewConfigFileParsingError(err)
	}

	return cfg, nil
}

func getFunctionAsAServiceConfigFromEnvironment() (*FunctionAsAServiceConfig, errors.ApplicationError) {
	systemConfig, err := getSystemConfigFromEnvironment()
	if err != nil {
		return nil, err
	}

	command := GetEnvironmentString("FUNCTION_AS_A_SERVICE_COMMAND", "")
	if command == "" {
		return nil, errors.NewMissingEnvironmentVariableError("FUNCTION_AS_A_SERVICE_COMMAND")
	}

	parallelExecution, err := GetEnvironmentInt("FUNCTION_AS_A_SERVICE_PARALLEL_EXECUTION_LIMIT", 10)
	if err != nil {
		return nil, err
	}

	executionTimeout, err := GetEnvironmentInt("FUNCTION_AS_A_SERVICE_EXECUTION_TIMEOUT", 30)
	if err != nil {
		return nil, err
	}

	httpConfig, err := getHttpServerConfigFromEnvironment()
	if err != nil {
		return nil, err
	}

	storageBackend := GetEnvironmentString("FUNCTION_AS_A_SERVICE_STORAGE_BACKEND", "")

	storageTtl, err := GetEnvironmentInt("FUNCTION_AS_A_SERVICE_STORAGE_TTL", 60)
	if err != nil {
		return nil, err
	}

	var redisConfig *RedisConfig
	if storageBackend == "redis" {
		redisConfig, err = getRedisConfigFromEnvironment()
		if err != nil {
			return nil, err
		}
	}

	var scyllaConfig *ScyllaConfig
	if storageBackend == "scylla" {
		scyllaConfig, err = getScyllaConfigFromEnvironment()
		if err != nil {
			return nil, err
		}
	}

	return &FunctionAsAServiceConfig{
		ApplicationConfig: ApplicationConfig{
			SystemConfig: systemConfig,
		},
		HttpServer:        httpConfig,
		Redis:             redisConfig,
		Scylla:            scyllaConfig,
		Command:           command,
		ParallelExecution: parallelExecution,
		ExecutionTimeout:  executionTimeout,
		StorageBackend:    storageBackend,
		StorageTtl:        storageTtl,
	}, nil
}
