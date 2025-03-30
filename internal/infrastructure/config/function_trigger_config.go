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

const (
	functionTriggerFaasDefaultSsl = false
	functionTriggerBatchSize      = 10
	functionTriggerBatchTimeout   = 10
)

type FunctionTriggerConfig struct {
	ApplicationConfig
	HttpServer       *HttpServerConfig `yaml:"http"`
	RabbitMQ         *RabbitMQConfig   `yaml:"rabbitmq"`
	TriggerName      string            `yaml:"trigger_name"`
	BatchTimeout     int64             `yaml:"batch_timeout"`
	BatchSize        int64             `yaml:"batch_size"`
	FaasHost         string            `yaml:"faas_host"`
	FaasPort         int64             `yaml:"faas_port"`
	FaasSsl          bool              `yaml:"faas_ssl"`
	FaasTimeout      int64             `yaml:"faas_timeout"`
	FaasFunctionName string            `yaml:"faas_function_name"`
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

	if err := c.RabbitMQ.Validate("rabbitmq", RabbitMQModeConsumer); err != nil {
		return errors.NewValidateConfigSectionError("rabbitmq", err)
	}

	if c.TriggerName == "" {
		return errors.NewMissingConfigValueError("trigger_name")
	}

	if c.BatchTimeout < 1 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError("rabbitmq.timeout")
	}

	if c.BatchSize < 1 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError("rabbitmq.batch_size")
	}

	if c.FaasHost == "" {
		return errors.NewMissingConfigValueError("faas_host")
	}

	if c.FaasPort == 0 {
		return errors.NewMissingConfigValueError("faas_port")
	}

	if c.FaasPort < 0 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError("faas_port")
	}

	if c.FaasFunctionName == "" {
		return errors.NewMissingConfigValueError("faas_function_name")
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
		HttpServer:   getDefaultHttpServerConfig(),
		RabbitMQ:     getDefaultRabbitMQConfig(RabbitMQModeConsumer),
		BatchTimeout: 10,
		BatchSize:    10,
		FaasTimeout:  30,
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

	rabbitMQConfig, err := getRabbitMQConfigFromEnvironment(RabbitMQModeConsumer)
	if err != nil {
		return nil, err
	}

	triggerName := GetEnvironmentString("FUNCTION_TRIGGER_TRIGGER_NAME", "")
	if triggerName == "" {
		return nil, errors.NewMissingEnvironmentVariableError("FUNCTION_TRIGGER_TRIGGER_NAME")
	}

	batchTimeout, err := GetEnvironmentInt("FUNCTION_TRIGGER_BATCH_TIMEOUT", 10)
	if err != nil {
		return nil, err
	}

	if batchTimeout < 1 {
		return nil, errors.NewConfigValueNeedsToBeGreaterZeroError("FUNCTION_TRIGGER_BATCH_TIMEOUT")
	}

	batchSize, err := GetEnvironmentInt("FUNCTION_TRIGGER_BATCH_SIZE", 10)
	if err != nil {
		return nil, err
	}

	if batchSize < 1 {
		return nil, errors.NewConfigValueNeedsToBeGreaterZeroError("FUNCTION_TRIGGER_BATCH_SIZE")
	}

	faasHost := GetEnvironmentString("FUNCTION_TRIGGER_FAAS_HOST", "")
	if faasHost == "" {
		return nil, errors.NewMissingEnvironmentVariableError("FUNCTION_TRIGGER_FAAS_HOST")
	}

	faasPort, err := GetEnvironmentInt("FUNCTION_TRIGGER_FAAS_PORT", 8080)
	if err != nil {
		return nil, err
	}

	if faasPort < 1 {
		return nil, errors.NewConfigValueNeedsToBeGreaterZeroError("FUNCTION_TRIGGER_FAAS_PORT")
	}

	faasTimeout, err := GetEnvironmentInt("FUNCTION_TRIGGER_FAAS_TIMEOUT", 30)
	if err != nil {
		return nil, err
	}

	if faasTimeout < 1 {
		return nil, errors.NewConfigValueNeedsToBeGreaterZeroError("FUNCTION_TRIGGER_FAAS_TIMEOUT")
	}

	if faasTimeout > 900 {
		return nil, errors.NewConfigValueNeedsToBeLessThanOrEqualValueError("FUNCTION_TRIGGER_FAAS_TIMEOUT", 900)
	}

	faasSsl, err := GetEnvironmentBool("FUNCTION_TRIGGER_FAAS_SSL", functionTriggerFaasDefaultSsl)
	if err != nil {
		return nil, err
	}

	faasFunction := GetEnvironmentString("FUNCTION_TRIGGER_FAAS_FUNCTION_NAME", "")
	if faasFunction == "" {
		return nil, errors.NewMissingEnvironmentVariableError("FUNCTION_TRIGGER_FAAS_FUNCTION_NAME")
	}

	return &FunctionTriggerConfig{
		ApplicationConfig: ApplicationConfig{
			SystemConfig: systemConfig,
		},
		HttpServer:       httpConfig,
		RabbitMQ:         rabbitMQConfig,
		TriggerName:      triggerName,
		BatchTimeout:     batchTimeout,
		BatchSize:        batchSize,
		FaasHost:         faasHost,
		FaasPort:         faasPort,
		FaasTimeout:      faasTimeout,
		FaasSsl:          faasSsl,
		FaasFunctionName: faasFunction,
	}, nil
}
