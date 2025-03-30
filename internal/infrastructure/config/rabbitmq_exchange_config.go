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
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"regexp"
	"strings"
)

type RabbitMQExchangeConfig struct {
	Name       string        `yaml:"name"`
	Type       string        `yaml:"type"`
	Passive    bool          `yaml:"passive"`
	Durable    bool          `yaml:"durable"`
	AutoDelete bool          `yaml:"auto_delete"`
	Internal   bool          `yaml:"internal"`
	NoWait     bool          `yaml:"no_wait"`
	Args       amqp091.Table `yaml:"args"`
}

func (c *RabbitMQExchangeConfig) Validate(path string) errors.ApplicationError {
	if c.Name == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.name", path))
	}

	if c.Type == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.type", path))
	}

	return nil
}

func parseRabbitMQExchanges(exchanges string) (map[string]*RabbitMQExchangeConfig, errors.ApplicationError) {
	exchangeMap := make(map[string]*RabbitMQExchangeConfig)

	re := regexp.MustCompile(rabbitMQRegexp)

	exchangeList := strings.Split(exchanges, ",")
	for _, exchange := range exchangeList {
		matches := re.FindStringSubmatch(exchange)
		if matches == nil {
			return nil, errors.NewMalformedEnvironmentVariableError("RABBITMQ_EXCHANGES", exchange, rabbitMQRegexp)
		}

		exchange := &RabbitMQExchangeConfig{
			Name:       matches[2],
			Type:       "direct",
			Passive:    false,
			Durable:    true,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Args:       nil,
		}

		exchangeMap[matches[1]] = exchange
	}

	return exchangeMap, nil
}
