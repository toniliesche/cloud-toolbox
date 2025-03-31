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

func TestValidateRabbitMQPublisherConfigFailsOnMissingExchangeConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQPublisherConfig()
	cfg.Exchanges = nil

	err := cfg.Validate("rabbitmq.publisher")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Equal(t, "config section `rabbitmq.publisher.exchanges` must exist", err.Error())
}

func TestValidateRabbitMQPublisherConfigFailsOnEmptyExchangeListConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQPublisherConfig()
	cfg.Exchanges = map[string]*config.RabbitMQExchangeConfig{}

	err := cfg.Validate("rabbitmq.publisher")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Equal(t, "config value `rabbitmq.publisher.exchanges` must not be an empty list", err.Error())
}

func TestValidateRabbitMQPublisherConfigFailsOnMissingVHost(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQPublisherConfig()
	cfg.VHost = ""

	err := cfg.Validate("rabbitmq.publisher")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.publisher.vhost` must exist")
}

func TestValidateRabbitMQPublisherConfigFailsOnMissingHost(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQPublisherConfig()
	cfg.Host = ""

	err := cfg.Validate("rabbitmq.publisher")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.publisher.host` must exist")
}

func TestValidateRabbitMQPublisherConfigFailsOnInvalidPort(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQPublisherConfig()
	cfg.Port = -1

	err := cfg.Validate("rabbitmq.publisher")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.publisher.port` must be greater than `0`")
}

func TestValidateRabbitMQPublisherConfigFailsOnMissingUsername(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQPublisherConfig()
	cfg.Username = ""

	err := cfg.Validate("rabbitmq.publisher")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.publisher.username` must exist")
}

func TestValidateRabbitMQPublisherConfigFailsOnMissingPassword(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQPublisherConfig()
	cfg.Password = ""

	err := cfg.Validate("rabbitmq.publisher")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.publisher.password` must exist")
}

func TestValidateRabbitMQPublisherConfigFailsOnInvalidExchangeConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQPublisherConfig()
	cfg.Exchanges["test"].Name = ""

	err := cfg.Validate("rabbitmq.publisher")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.publisher.exchanges.test.name` must exist")
}

func TestValidateRabbitMQPublisherConfigSucceedsOnValidConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQPublisherConfig()

	err := cfg.Validate("rabbitmq.publisher")
	if !assert.NoError(t, err, "expected no error") {
		return
	}
}

func getValidRabbitMQPublisherConfig() *config.RabbitMQPublisherConfig {
	return &config.RabbitMQPublisherConfig{
		VHost:    "/",
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		Ssl:      false,
		Exchanges: map[string]*config.RabbitMQExchangeConfig{
			"test": getValidRabbitMQExchangeConfig(),
		},
	}
}
