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

func TestValidateRabbitMQConsumerConfigFailsOnMissingQueueConfig(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()
	cfg.Queues = nil

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Equal(t, "config section `rabbitmq.consumer.queues` must exist", err.Error())
}

func TestValidateRabbitMQConsumerConfigFailsOnEmptyQueueListConfig(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()
	cfg.Queues = map[string]*config.RabbitMQQueueConfig{}

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Equal(t, "config value `rabbitmq.consumer.queues` must not be an empty list", err.Error())
}

func TestValidateRabbitMQConsumerConfigFailsOnMissingVHost(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()
	cfg.VHost = ""

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.consumer.vhost` must exist")
}

func TestValidateRabbitMQConsumerConfigFailsOnMissingHost(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()
	cfg.Host = ""

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.consumer.host` must exist")
}

func TestValidateRabbitMQConsumerConfigFailsOnMissingPort(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()
	cfg.Port = 0

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.consumer.port` must exist")
}

func TestValidateRabbitMQConsumerConfigFailsOnInvalidPort(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()
	cfg.Port = -1

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.consumer.port` must be greater than `0`")
}

func TestValidateRabbitMQConsumerConfigFailsOnMissingUsername(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()
	cfg.Username = ""

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.consumer.username` must exist")
}

func TestValidateRabbitMQConsumerConfigFailsOnMissingPassword(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()
	cfg.Password = ""

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.consumer.password` must exist")
}

func TestValidateRabbitMQConsumerConfigFailsOnInvalidQueueConfig(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()
	cfg.Queues["test"].Name = ""

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "config value `rabbitmq.consumer.queues.test.name` must exist")
}

func TestValidateRabbitMQConsumerConfigSucceedsOnValidConfig(t *testing.T) {
	cfg := getValidRabbitMQConsumerConfig()

	err := cfg.Validate("rabbitmq.consumer")
	if !assert.NoError(t, err, "expected no error") {
		return
	}
}

func getValidRabbitMQConsumerConfig() *config.RabbitMQConsumerConfig {
	return &config.RabbitMQConsumerConfig{
		VHost:    "/",
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		Ssl:      false,
		Queues: map[string]*config.RabbitMQQueueConfig{
			"test": getValidRabbitMQQueueConfig(),
		},
	}
}
