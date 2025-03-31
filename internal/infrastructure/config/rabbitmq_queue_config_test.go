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
	"testing"
)

func TestValidateRabbitMQQueueConfigFailsOnMissingName(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQQueueConfig()
	cfg.Name = ""

	if err := cfg.Validate("rabbitmq"); err == nil {
		t.Fatal("Expected error, got none")
	}
}

func TestValidateRabbitMQQueueConfigFailsOnMissingType(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQQueueConfig()
	cfg.Type = ""

	if err := cfg.Validate("rabbitmq"); err == nil {
		t.Fatal("Expected error, got none")
	}
}

func TestValidateRabbitMQQueueConfigSucceedsOnValidConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidRabbitMQQueueConfig()

	if err := cfg.Validate("rabbitmq"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func getValidRabbitMQQueueConfig() *config.RabbitMQQueueConfig {
	return &config.RabbitMQQueueConfig{
		Name:       "test-queue",
		Type:       "direct",
		AutoAck:    false,
		AutoDelete: false,
		Durable:    true,
		Exclusive:  false,
		NoLocal:    false,
		NoWait:     false,
		Args:       nil,
	}
}
