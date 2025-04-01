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

package rabbitmq

import (
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"context"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

const PublisherLogIdentifier = "RabbitMQPublisher"

type Publisher struct {
	cfg      *config.RabbitMQConfig
	rabbitMQ *amqp091.Connection
	logger   *zerolog.Logger
	channel  *amqp091.Channel
	context  context.Context
}

func (p Publisher) Publish(msg amqp091.Publishing, destination string, routingKey string) errors.ApplicationError {
	err := p.channel.PublishWithContext(
		p.context,
		destination,
		routingKey,
		false,
		false,
		msg,
	)

	if err != nil {
		p.logger.Error().
			Err(err).
			Msgf("[%s] Error publishing message to RabbitMQ", PublisherLogIdentifier)
		return errors.NewGenericError(err)
	}

	p.logger.Info().
		Msgf("[%s] Message published to RabbitMQ with routing key: %s", PublisherLogIdentifier, routingKey)

	return nil
}

func NewPublisher(container *di.Container) (*Publisher, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError(PublisherLogIdentifier)
	}

	if container.RabbitMQConfig == nil {
		return nil, errors.NewResolveDependencyError(PublisherLogIdentifier, "RabbitMQConfig")
	}

	if err := container.RabbitMQConfig.Validate("", config.RabbitMQModePublisher); err != nil {
		return nil, errors.NewInvalidConfigError(PublisherLogIdentifier, err)
	}

	if container.Logger == nil {
		return nil, errors.NewResolveDependencyError(PublisherLogIdentifier, "Logger")
	}

	if container.RabbitMQConnectionPublisher == nil {
		return nil, errors.NewResolveDependencyError(PublisherLogIdentifier, "RabbitMQConnectionPublisher")
	}

	if container.Context == nil {
		container.RabbitMQConnectionPublisher.Close()
		return nil, errors.NewResolveDependencyError(PublisherLogIdentifier, "Context")
	}

	ch, err := container.RabbitMQConnectionPublisher.Channel()
	if err != nil {
		container.RabbitMQConnectionPublisher.Close()
		return nil, errors.NewGenericError(err)
	}

	go func() {
		<-container.Context.Done()
		if err := ch.Close(); err != nil {
			container.Logger.Error().
				Err(err).
				Msgf("[%s] Error closing channel", PublisherLogIdentifier)
		}
	}()

	return &Publisher{
		context:  container.Context,
		cfg:      container.RabbitMQConfig,
		rabbitMQ: container.RabbitMQConnectionPublisher,
		logger:   container.Logger,
		channel:  ch,
	}, nil
}
