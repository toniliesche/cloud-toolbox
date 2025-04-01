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

func TestValidateEventPublisherConfigFailsOnMissingSystemConfig(t *testing.T) {
	t.Parallel()

	cfg := getValidEventPublisherConfig("rabbitmq")
	cfg.SystemConfig = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing system section error") {
		return
	}

	if !assert.Equal(t, "config section `system` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateEventPublisherConfigFailsOnInvalidSystemConfig(t *testing.T) {
	t.Parallel()

	cfg := getValidEventPublisherConfig("rabbitmq")
	cfg.SystemConfig.Log = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid system config error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `system`", "unexpected error message") {
		return
	}
}

func TestValidateEventPublisherConfigFailsOnMissingHttpConfig(t *testing.T) {
	t.Parallel()

	cfg := getValidEventPublisherConfig("rabbitmq")
	cfg.HttpServer = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing http section error") {
		return
	}

	if !assert.Equal(t, "config section `http` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateEventPublisherConfigFailsOnInvalidHttpConfig(t *testing.T) {
	t.Parallel()

	cfg := getValidEventPublisherConfig("rabbitmq")
	cfg.HttpServer.Port = 0

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid http config error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `http`", "unexpected error message") {
		return
	}
}

func TestValidateEventPublisherConfigFailsOnMissingPublisherName(t *testing.T) {
	t.Parallel()
	cfg := getValidEventPublisherConfig("rabbitmq")
	cfg.PublisherName = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing publisher name error") {
		return
	}

	if !assert.Equal(t, "config value `publisher_name` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateEventPublisherConfigFailsOnMissingPublisherBackend(t *testing.T) {
	t.Parallel()
	cfg := getValidEventPublisherConfig("rabbitmq")
	cfg.PublisherBackend = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing publisher backend error") {
		return
	}

	if !assert.Equal(t, "config value `publisher_backend` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateEventPublisherConfigFailsOnInvalidPublisherBackend(t *testing.T) {
	t.Parallel()
	cfg := getValidEventPublisherConfig("rabbitmq")
	cfg.PublisherBackend = "invalid-backend"

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid publisher backend error") {
		return
	}

	if !assert.Equal(t, "config value `publisher_backend` must be one of [rabbitmq], but is `invalid-backend`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateEventPublisherConfigFailsOnMissingRabbitMQConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidEventPublisherConfig("rabbitmq")
	cfg.RabbitMQ = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing rabbitmq section error") {
		return
	}

	if !assert.Equal(t, "config section `rabbitmq` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateEventPublisherConfigFailsOnInvalidRabbitMQConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidEventPublisherConfig("rabbitmq")
	cfg.RabbitMQ.Publisher.Host = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid rabbitmq config error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `rabbitmq`", "unexpected error message") {
		return
	}
}

func TestValidateEventPublisherConfigSucceedsOnValidConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidEventPublisherConfig("rabbitmq")

	err := cfg.Validate()
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func getValidEventPublisherConfig(backend string) *config.EventPublisherConfig {
	var rabbitMQConfig *config.RabbitMQConfig

	switch backend {
	case "rabbitmq":
		rabbitMQConfig = getValidRabbitMQConfig(config.RabbitMQModePublisher)
	}

	return &config.EventPublisherConfig{
		ApplicationConfig: config.ApplicationConfig{
			SystemConfig: getValidSystemConfig(),
		},
		HttpServer:       getValidHttpServerConfig(),
		RabbitMQ:         rabbitMQConfig,
		PublisherName:    "test-publisher",
		PublisherBackend: backend,
	}
}
