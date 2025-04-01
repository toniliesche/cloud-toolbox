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

package setup

import (
	"cloud-toolbox/internal/application/ep"
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/http"
	event_publisher "cloud-toolbox/internal/infrastructure/http/event-publisher"
	"cloud-toolbox/internal/infrastructure/rabbitmq"
	"cloud-toolbox/internal/infrastructure/rabbitmq/connectors"
)

func (b *ContainerBuilder) setupEp(container *di.Container) errors.ApplicationError {
	var err errors.ApplicationError

	if b.epConfig == nil {
		b.logger.Trace().
			Msgf("[%s] `Endpoint` config is nil", ContainerBuilderLogIdentifier)
		return errors.NewResolveDependencyError("Endpoint", "epConfig")
	}

	b.logger.Trace().
		Msgf("[%s] Validating `Endpoint` config", ContainerBuilderLogIdentifier)
	if err = b.epConfig.Validate(); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] `Endpoint` config is invalid", ContainerBuilderLogIdentifier)
		return errors.NewApplicationSetupError("Endpoint", err)
	}
	container.EventPublisherConfig = b.epConfig

	switch b.epConfig.PublisherBackend {
	case "rabbitmq":
		err := b.setupEpRabbitMQ(container)
		if err != nil {
			return err
		}
	}

	b.logger.Trace().
		Msgf("[%s] Initializing `EventPublisher` service", ContainerBuilderLogIdentifier)
	container.EventPublisher, err = ep.NewEventPublisher(container)
	if err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `EventPublisher` service", ContainerBuilderLogIdentifier)
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Retrieving HttpServerConfig from `EventPublisher` config", ContainerBuilderLogIdentifier)
	container.HttpServerConfig = b.epConfig.HttpServer

	b.logger.Trace().
		Msgf("[%s] Initializing `HttpServer`", ContainerBuilderLogIdentifier)
	if container.HttpServer, err = http.NewServer(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `HttpServer`", ContainerBuilderLogIdentifier)
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Initializing `EventPublisher` http handler", ContainerBuilderLogIdentifier)
	if container.EventPublisherHttpHandler, err = event_publisher.NewEventPublisherHandler(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `EventPublisher` http handler", ContainerBuilderLogIdentifier)
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Registering routes for `EventPublisher` http handler", ContainerBuilderLogIdentifier)
	if err = container.HttpServer.RegisterHandler(container.EventPublisherHttpHandler); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error registering routes for `EventPublisher` http handler", ContainerBuilderLogIdentifier)
		return err
	}

	b.logger.Trace().
		Msgf("[%s] `EventPublisher` setup completed", ContainerBuilderLogIdentifier)

	return nil
}

func (b *ContainerBuilder) setupEpRabbitMQ(container *di.Container) errors.ApplicationError {
	var err errors.ApplicationError

	b.logger.Trace().
		Msgf("[%s] Retrieving RabbitMQConfig from `EventPublisher` config", ContainerBuilderLogIdentifier)
	container.RabbitMQConfig = b.epConfig.RabbitMQ

	if container.RabbitMQConnectionPublisher, err = connectors.NewRabbitMQ(container, config.RabbitMQModePublisher); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `RabbitMQ` connection", ContainerBuilderLogIdentifier)
		return errors.NewApplicationSetupError("RabbitMQ", err)
	}

	if container.RabbitMQPublisher, err = rabbitmq.NewPublisher(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `RabbitMQ` publisher", ContainerBuilderLogIdentifier)
		return errors.NewApplicationSetupError("RabbitMQ", err)
	}

	if container.EventPublisherBackend, err = ep.NewRabbitMQBackend(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `RabbitMQ` backend", ContainerBuilderLogIdentifier)
		return errors.NewApplicationSetupError("RabbitMQ", err)
	}

	return nil
}

func NewBuilder() *ContainerBuilder {
	return &ContainerBuilder{}
}
