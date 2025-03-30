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
	"github.com/streadway/amqp"
	"regexp"
	"strings"
)

type RabbitMQQueueConfig struct {
	Name       string     `yaml:"name"`
	Type       string     `yaml:"type"`
	AutoAck    bool       `yaml:"auto_ack"`
	AutoDelete bool       `yaml:"auto_delete"`
	Durable    bool       `yaml:"durable"`
	Exclusive  bool       `yaml:"exclusive"`
	NoLocal    bool       `yaml:"no_local"`
	NoWait     bool       `yaml:"no_wait"`
	Args       amqp.Table `yaml:"args"`
}

func (c *RabbitMQQueueConfig) Validate(path string) errors.ApplicationError {
	if c.Name == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.name", path))
	}

	if c.Type == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.type", path))
	}

	return nil
}

func parseRabbitMQQueues(queues string) (map[string]*RabbitMQQueueConfig, errors.ApplicationError) {
	queueMap := make(map[string]*RabbitMQQueueConfig)

	re := regexp.MustCompile(rabbitMQRegexp)

	queueList := strings.Split(queues, ",")
	for _, queue := range queueList {
		matches := re.FindStringSubmatch(queue)
		if matches == nil {
			return nil, errors.NewMalformedEnvironmentVariableError("RABBITMQ_QUEUES", queue, rabbitMQRegexp)
		}

		queue := RabbitMQQueueConfig{
			Name:       matches[2],
			Type:       "direct",
			AutoAck:    false,
			AutoDelete: false,
			Durable:    true,
			Exclusive:  false,
			NoLocal:    false,
			NoWait:     false,
			Args:       nil,
		}

		queueMap[matches[1]] = &queue
	}

	return queueMap, nil
}
