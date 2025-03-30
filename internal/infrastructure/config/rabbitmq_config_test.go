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

func TestValidateRabbitMQConfigFailsOnMissingConsumerConfig(t *testing.T) {
	cfg := getValidRabbitMQConfig(config.RabbitMQModeConsumer)
	cfg.Consumer = nil

	err := cfg.Validate("rabbitmq", config.RabbitMQModeConsumer)
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Equal(t, "config section `rabbitmq.consumer` must exist", err.Error())
}

func TestValidateRabbitMQConfigFailsOnInvalidConsumerConfig(t *testing.T) {
	cfg := getValidRabbitMQConfig(config.RabbitMQModeConsumer)
	cfg.Consumer.Queues = nil

	err := cfg.Validate("rabbitmq", config.RabbitMQModeConsumer)
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "Error in config section `rabbitmq.consumer`")
}

func TestValidateRabbitMQConfigFailsOnMissingPublisherConfig(t *testing.T) {
	cfg := getValidRabbitMQConfig(config.RabbitMQModePublisher)
	cfg.Producer = nil

	err := cfg.Validate("rabbitmq", config.RabbitMQModePublisher)
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Equal(t, "config section `rabbitmq.producer` must exist", err.Error())
}

func TestValidateRabbitMQConfigFailsOnInvalidPublisherConfig(t *testing.T) {
	cfg := getValidRabbitMQConfig(config.RabbitMQModePublisher)
	cfg.Producer.Exchanges = nil

	err := cfg.Validate("rabbitmq", config.RabbitMQModePublisher)
	if !assert.Error(t, err, "expected error") {
		return
	}
	assert.Contains(t, err.Error(), "Error in config section `rabbitmq.producer`")
}

func TestValidateRabbitMQConfigSucceedsOnValidConsumerConfig(t *testing.T) {
	cfg := getValidRabbitMQConfig(config.RabbitMQModeConsumer)

	err := cfg.Validate("rabbitmq", config.RabbitMQModeConsumer)
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func TestValidateRabbitMQConfigSucceedsOnValidPublisherConfig(t *testing.T) {
	cfg := getValidRabbitMQConfig(config.RabbitMQModePublisher)

	err := cfg.Validate("rabbitmq", config.RabbitMQModePublisher)
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func TestValidateRabbitMQConfigSucceedsOnValidConsumerAndPublisherConfig(t *testing.T) {
	cfg := getValidRabbitMQConfig(config.RabbitMQModePublisherConsumer)

	err := cfg.Validate("rabbitmq", config.RabbitMQModePublisherConsumer)
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func getValidRabbitMQConfig(mode int) *config.RabbitMQConfig {
	cfg := &config.RabbitMQConfig{}

	if (mode & config.RabbitMQModeConsumer) == config.RabbitMQModeConsumer {
		cfg.Consumer = getValidRabbitMQConsumerConfig()
	}

	if (mode & config.RabbitMQModePublisher) == config.RabbitMQModePublisher {
		cfg.Producer = getValidRabbitMQProducerConfig()
	}

	return cfg
}
