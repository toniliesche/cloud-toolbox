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
)

type FunctionTriggerConfig struct {
	ApplicationConfig
	HttpServer *HttpServerConfig `yaml:"http"`
	RabbitMQ   *RabbitMQConfig   `yaml:"rabbitmq"`
}

func (c *FunctionTriggerConfig) Validate() errors.ApplicationError {
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

	if c.RabbitMQ == nil {
		return errors.NewMissingConfigSectionError("rabbitmq")
	}

	if err := c.RabbitMQ.Validate("rabbitmq", RabbitMQModeSubscriber); err != nil {
		return errors.NewValidateConfigSectionError("rabbitmq", err)
	}

	return nil
}

func ProvideFunctionTriggerConfig() (*FunctionTriggerConfig, errors.ApplicationError) {
	configFile := GetEnvironmentString("FUNCTION_TRIGGER_CONFIG_FILE", "")

	var cfg *FunctionTriggerConfig
	var err errors.ApplicationError

	if configFile != "" {
		cfg, err = getFunctionTriggerConfigFromFile(configFile)
	} else {
		cfg, err = getFunctionTriggerConfigFromEnvironment()
	}

	if err != nil {
		return nil, errors.NewConfigCreationError(err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, errors.NewConfigValidationError(err)
	}

	return cfg, nil
}

func getDefaultFunctionTriggerConfig() *FunctionTriggerConfig {
	return &FunctionTriggerConfig{
		HttpServer: getDefaultHttpServerConfig(),
		RabbitMQ:   getDefaultRabbitMQConfig(),
	}
}

func getFunctionTriggerConfigFromFile(file string) (*FunctionTriggerConfig, errors.ApplicationError) {
	fileContents, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.NewConfigFileReadingError(err)
	}

	cfg := getDefaultFunctionTriggerConfig()
	if err = yaml.Unmarshal(fileContents, cfg); err != nil {
		return nil, errors.NewConfigFileParsingError(err)
	}

	return cfg, nil
}

func getFunctionTriggerConfigFromEnvironment() (*FunctionTriggerConfig, errors.ApplicationError) {
	systemConfig, err := getSystemConfigFromEnvironment()
	if err != nil {
		return nil, err
	}

	httpConfig, err := getHttpServerConfigFromEnvironment()
	if err != nil {
		return nil, err
	}

	rabbitMQConfig, err := getRabbitMQConfigFromEnvironment()
	if err != nil {
		return nil, err
	}

	return &FunctionTriggerConfig{
		ApplicationConfig: ApplicationConfig{
			SystemConfig: systemConfig,
		},
		HttpServer: httpConfig,
		RabbitMQ:   rabbitMQConfig,
	}, nil
}
