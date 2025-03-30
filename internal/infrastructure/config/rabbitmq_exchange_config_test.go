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

func TestValidateRabbitMQExchangeConfigFailsOnMissingName(t *testing.T) {
	cfg := getValidRabbitMQExchangeConfig()
	cfg.Name = ""

	if err := cfg.Validate("rabbitmq"); err == nil {
		t.Fatal("Expected error, got none")
	}
}

func TestValidateRabbitMQExchangeConfigFailsOnMissingType(t *testing.T) {
	cfg := getValidRabbitMQExchangeConfig()
	cfg.Type = ""

	if err := cfg.Validate("rabbitmq"); err == nil {
		t.Fatal("Expected error, got none")
	}
}

func TestValidateRabbitMQExchangeConfigSucceedsOnValidConfig(t *testing.T) {
	cfg := getValidRabbitMQExchangeConfig()

	if err := cfg.Validate("rabbitmq"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func getValidRabbitMQExchangeConfig() *config.RabbitMQExchangeConfig {
	return &config.RabbitMQExchangeConfig{
		Name:       "test",
		Type:       "direct",
		Passive:    false,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}
