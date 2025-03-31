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

type EventPublisherConfig struct {
	ApplicationConfig
	HttpServer         *HttpServerConfig `yaml:"http"`
	RabbitMQ           *RabbitMQConfig   `yaml:"rabbitmq"`
	EventPublisherName string            `yaml:"event_publisher_name"`
}

func (c *EventPublisherConfig) Validate() errors.ApplicationError {
	if c.SystemConfig == nil {
		return errors.NewMissingConfigSectionError("system")
	}

	if err := c.SystemConfig.Validate("system"); err != nil {
		return errors.NewValidateConfigSectionError("system", err)
	}

	if c.HttpServer == nil {
		return errors.NewMissingConfigSectionError("http")
	}

	if err := c.HttpServer.Validate("http"); err != nil {
		return errors.NewValidateConfigSectionError("http", err)
	}

	if c.RabbitMQ == nil {
		return errors.NewMissingConfigSectionError("rabbitmq")
	}

	if err := c.RabbitMQ.Validate("rabbitmq", RabbitMQModePublisher); err != nil {
		return errors.NewValidateConfigSectionError("rabbitmq", err)
	}

	if c.EventPublisherName == "" {
		return errors.NewMissingConfigValueError("event_publisher_name")
	}

	return nil
}

func ProvideEventPublisherConfig() (*EventPublisherConfig, errors.ApplicationError) {
	configFile := GetEnvironmentString("EVENT_PUBLISHER_CONFIG_FILE", "")

	var cfg *EventPublisherConfig
	var err errors.ApplicationError

	if configFile != "" {
		cfg, err = getEventPublisherConfigFromFile(configFile)
	} else {
		cfg, err = getEventPublisherConfigFromEnvironment()
	}

	if err != nil {
		return nil, errors.NewConfigCreationError(err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, errors.NewConfigValidationError(err)
	}

	return cfg, nil
}

func getDefaultEventPublisherConfig() *EventPublisherConfig {
	return &EventPublisherConfig{
		HttpServer: getDefaultHttpServerConfig(),
		RabbitMQ:   getDefaultRabbitMQConfig(RabbitMQModePublisher),
	}
}

func getEventPublisherConfigFromFile(file string) (*EventPublisherConfig, errors.ApplicationError) {
	fileContents, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.NewConfigFileReadingError(err)
	}

	cfg := getDefaultEventPublisherConfig()
	if err = yaml.Unmarshal(fileContents, cfg); err != nil {
		return nil, errors.NewConfigFileParsingError(err)
	}

	return cfg, nil
}

func getEventPublisherConfigFromEnvironment() (*EventPublisherConfig, errors.ApplicationError) {
	systemConfig, err := getSystemConfigFromEnvironment()
	if err != nil {
		return nil, err
	}

	httpConfig, err := getHttpServerConfigFromEnvironment()
	if err != nil {
		return nil, err
	}

	rabbitMQConfig, err := getRabbitMQConfigFromEnvironment(RabbitMQModePublisher)
	if err != nil {
		return nil, err
	}

	eventPublisherName := GetEnvironmentString("EVENT_PUBLISHER_PUBLISHER_NAME", "")
	if eventPublisherName == "" {
		return nil, errors.NewMissingEnvironmentVariableError("EVENT_PUBLISHER_PUBLISHER_NAME")
	}

	return &EventPublisherConfig{
		ApplicationConfig: ApplicationConfig{
			SystemConfig: systemConfig,
		},
		HttpServer:         httpConfig,
		RabbitMQ:           rabbitMQConfig,
		EventPublisherName: eventPublisherName,
	}, nil
}
