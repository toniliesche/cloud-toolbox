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
)

const (
	rabbitMQDefaultVhost          = "/"
	rabbitMQDefaultHost           = "localhost"
	rabbitMQDefaultPort           = 5672
	rabbitMQDefaultSsl            = false
	rabbitMQRegexp                = `^([a-zA-Z_-]+)=([a-zA-Z0-9_.:#*-]{1,255})$`
	RabbitMQModePublisher         = 1
	RabbitMQModeConsumer          = 2
	RabbitMQModePublisherConsumer = 3
)

type RabbitMQConfig struct {
	Publisher *RabbitMQPublisherConfig `yaml:"publisher"`
	Consumer  *RabbitMQConsumerConfig  `yaml:"consumer"`
}

func (c *RabbitMQConfig) Validate(path string, mode int) errors.ApplicationError {
	if (mode & RabbitMQModePublisher) == RabbitMQModePublisher {
		if c.Publisher == nil {
			return errors.NewMissingConfigSectionError(fmt.Sprintf("%s.publisher", path))
		}

		if err := c.Publisher.Validate(fmt.Sprintf("%s.publisher", path)); err != nil {
			return errors.NewValidateConfigSectionError(fmt.Sprintf("%s.publisher", path), err)
		}
	}

	if (mode & RabbitMQModeConsumer) == RabbitMQModeConsumer {
		if c.Consumer == nil {
			return errors.NewMissingConfigSectionError(fmt.Sprintf("%s.consumer", path))
		}

		if err := c.Consumer.Validate(fmt.Sprintf("%s.consumer", path)); err != nil {
			return errors.NewValidateConfigSectionError(fmt.Sprintf("%s.consumer", path), err)
		}
	}

	return nil
}

func (c *RabbitMQConfig) GetExchange(id string) (*RabbitMQExchangeConfig, errors.ApplicationError) {
	if c.Publisher == nil {
		return nil, errors.NewMissingConfigSectionError("rabbitmq.publisher")
	}

	return c.Publisher.GetExchange(id)
}

func (c *RabbitMQConfig) GetQueue(id string) (*RabbitMQQueueConfig, errors.ApplicationError) {
	if c.Consumer == nil {
		return nil, errors.NewMissingConfigSectionError("rabbitmq.consumer")
	}

	return c.Consumer.GetQueue(id)
}

func getDefaultRabbitMQConfig(mode int) *RabbitMQConfig {
	cfg := &RabbitMQConfig{}

	if (mode & RabbitMQModePublisher) == RabbitMQModePublisher {
		cfg.Publisher = getDefaultRabbitMQPublisherConfig()
	}

	if (mode & RabbitMQModeConsumer) == RabbitMQModeConsumer {
		cfg.Consumer = getDefaultRabbitMQConsumerConfig()
	}

	return cfg
}

func getRabbitMQConfigFromEnvironment(mode int) (*RabbitMQConfig, errors.ApplicationError) {
	cfg := getDefaultRabbitMQConfig(mode)

	if (mode & RabbitMQModePublisher) == RabbitMQModePublisher {
		publisherCfg, err := getRabbitMQPublisherConfigFromEnvironment()
		if err != nil {
			return nil, err
		}

		cfg.Publisher = publisherCfg
	}

	if (mode & RabbitMQModeConsumer) == RabbitMQModeConsumer {
		consumerCfg, err := getRabbitMQConsumerConfigFromEnvironment()
		if err != nil {
			return nil, err
		}

		cfg.Consumer = consumerCfg
	}

	return cfg, nil
}
