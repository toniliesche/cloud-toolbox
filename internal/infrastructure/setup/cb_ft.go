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
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/rabbitmq"
	"cloud-toolbox/internal/infrastructure/rabbitmq/connectors"
	"cloud-toolbox/internal/infrastructure/rabbitmq/function-trigger"
)

func (b *ContainerBuilder) setupFt(container *di.Container) errors.ApplicationError {
	var err errors.ApplicationError

	if b.ftConfig == nil {
		b.logger.Trace().
			Msgf("[%s] `FunctionTrigger` config is nil", ContainerBuilderLogIdentifier)
		return errors.NewResolveDependencyError("Function trigger", "ftConfig")
	}

	b.logger.Trace().
		Msgf("[%s] Validating `FunctionTrigger` config", ContainerBuilderLogIdentifier)
	if err = b.ftConfig.Validate(); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] `FunctionTrigger` config is invalid", ContainerBuilderLogIdentifier)
		return errors.NewApplicationSetupError("Function trigger", err)
	}
	container.FunctionTriggerConfig = b.ftConfig

	switch b.ftConfig.TriggerSource {
	case "rabbitmq":
		err := b.setupFtRabbitMQ(container)
		if err != nil {
			return err
		}
	}

	b.logger.Trace().
		Msgf("[%s] `FunctionTrigger` setup complete", ContainerBuilderLogIdentifier)

	return nil
}

func (b *ContainerBuilder) setupFtRabbitMQ(container *di.Container) errors.ApplicationError {
	var err errors.ApplicationError
	b.logger.Trace().
		Msgf("[%s] Retrieving RabbitMQConfig from `FunctionTrigger` config", ContainerBuilderLogIdentifier)
	container.RabbitMQConfig = b.ftConfig.RabbitMQ

	b.logger.Trace().
		Msgf("[%s] Initializing `RabbitMQ` consumer connection", ContainerBuilderLogIdentifier)
	if container.RabbitMQConnectionConsumer, err = connectors.NewRabbitMQ(container, config.RabbitMQModeConsumer); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing RabbitMQ consumer connection", ContainerBuilderLogIdentifier)
		return errors.NewContainerMissingError("Redis")
	}

	b.logger.Trace().
		Msgf("[%s] Initializing `RabbitMQ` consumer", ContainerBuilderLogIdentifier)
	if container.RabbitMQConsumer, err = rabbitmq.NewConsumer(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing RabbitMQ consumer", ContainerBuilderLogIdentifier)
		return errors.NewContainerMissingError("RabbitMQConsumer")
	}

	container.RabbitMQConsumer.
		WithBatchSize(b.ftConfig.BatchSize).
		WithBatchTimeout(b.ftConfig.BatchTimeout)

	b.logger.Trace().
		Msgf("[%s] Initializing `FunctionTrigger` RabbitMQ handler", ContainerBuilderLogIdentifier)
	if container.FunctionTriggerRabbitMQHandler, err = function_trigger.NewFunctionTriggerHandler(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `FunctionTrigger` RabbitMQ handler", ContainerBuilderLogIdentifier)
		return errors.NewContainerMissingError("RabbitMQFunctionTriggerHandler")
	}

	b.logger.Trace().
		Msgf("[%s] Registering `FunctionTrigger` RabbitMQ handler", ContainerBuilderLogIdentifier)
	if err = container.RabbitMQConsumer.RegisterHandler(container.FunctionTriggerRabbitMQHandler); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error registering `FunctionTrigger` RabbitMQ handler", ContainerBuilderLogIdentifier)
		return errors.NewGenericError(err)
	}

	return nil
}
