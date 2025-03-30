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

func TestValidateRabbitMQProducerConfigFailsOnMissingExchangeConfig(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()
	cfg.Exchanges = nil

	err := cfg.Validate("rabbitmq.producer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Equal(t, "config section `rabbitmq.producer.exchanges` must exist", err.Error())
}

func TestValidateRabbitMQProducerConfigFailsOnEmptyExchangeListConfig(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()
	cfg.Exchanges = map[string]*config.RabbitMQExchangeConfig{}

	err := cfg.Validate("rabbitmq.producer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Equal(t, "config value `rabbitmq.producer.exchanges` must not be an empty list", err.Error())
}

func TestValidateRabbitMQProducerConfigFailsOnMissingVHost(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()
	cfg.VHost = ""

	err := cfg.Validate("rabbitmq.producer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.producer.vhost` must exist")
}

func TestValidateRabbitMQProducerConfigFailsOnMissingHost(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()
	cfg.Host = ""

	err := cfg.Validate("rabbitmq.producer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.producer.host` must exist")
}

func TestValidateRabbitMQProducerConfigFailsOnMissingPort(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()
	cfg.Port = 0

	err := cfg.Validate("rabbitmq.producer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.producer.port` must exist")
}

func TestValidateRabbitMQProducerConfigFailsOnInvalidPort(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()
	cfg.Port = -1

	err := cfg.Validate("rabbitmq.producer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.producer.port` must be greater than `0`")
}

func TestValidateRabbitMQProducerConfigFailsOnMissingUsername(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()
	cfg.Username = ""

	err := cfg.Validate("rabbitmq.producer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.producer.username` must exist")
}

func TestValidateRabbitMQProducerConfigFailsOnMissingPassword(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()
	cfg.Password = ""

	err := cfg.Validate("rabbitmq.producer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.producer.password` must exist")
}

func TestValidateRabbitMQProducerConfigFailsOnInvalidExchangeConfig(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()
	cfg.Exchanges["test"].Name = ""

	err := cfg.Validate("rabbitmq.producer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.producer.exchanges.test.name` must exist")
}

func TestValidateRabbitMQProducerConfigSucceedsOnValidConfig(t *testing.T) {
	cfg := getValidRabbitMQProducerConfig()

	err := cfg.Validate("rabbitmq.producer")
	if !assert.NoError(t, err, "expected no error") {
		return
	}
}

func getValidRabbitMQProducerConfig() *config.RabbitMQProducerConfig {
	return &config.RabbitMQProducerConfig{
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
