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

package connectors

import (
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/config/interfaces"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(container *di.Container, mode int) (*amqp091.Connection, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError("RabbitMQ")
	}

	cfg := container.RabbitMQConfig
	if cfg == nil {
		return nil, errors.NewResolveDependencyError("RabbitMQ", "RabbitMQConfig")
	}

	if err := cfg.Validate("rabbitmq", mode); err != nil {
		return nil, errors.NewInvalidConfigError("RabbitMQConfig", err)
	}

	var serverCfg interfaces.RabbitMQServerConfig
	if mode == config.RabbitMQModeConsumer {
		serverCfg = container.RabbitMQConfig.Consumer
	} else {
		serverCfg = container.RabbitMQConfig.Producer
	}

	conn, err := amqp091.Dial(serverCfg.Addr())
	if err != nil {
		return nil, errors.NewConstructionFailedError("RabbitMQ", err)
	}

	return conn, nil
}
