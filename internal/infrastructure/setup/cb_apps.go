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
	"cloud-toolbox/internal/infrastructure/config/interfaces"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"fmt"
)

func (b *ContainerBuilder) SetEpConfig(cfg *config.EventPublisherConfig) *ContainerBuilder {
	b.epConfig = cfg
	b.application = "ep"

	return b
}

func (b *ContainerBuilder) SetFaasConfig(cfg *config.FunctionAsAServiceConfig) *ContainerBuilder {
	b.faasConfig = cfg
	b.application = "faas"

	return b
}

func (b *ContainerBuilder) SetFtConfig(cfg *config.FunctionTriggerConfig) *ContainerBuilder {
	b.ftConfig = cfg
	b.application = "ft"

	return b
}

func (b *ContainerBuilder) getApplicationConfig() (interfaces.ApplicationConfig, errors.ApplicationError) {
	switch b.application {
	case "ep":
		return b.epConfig, nil
	case "faas":
		return b.faasConfig, nil
	case "ft":
		return b.ftConfig, nil
	default:
		return nil, errors.NewContainerConfigMissingError()
	}
}

func (b *ContainerBuilder) getComponentType() string {
	switch b.application {
	case "ep":
		return "event-publisher"
	case "faas":
		return "function-as-a-service"
	case "ft":
		return "function-trigger"
	default:
		return "unknown"
	}
}

func (b *ContainerBuilder) setupApplication(container *di.Container) errors.ApplicationError {
	switch b.application {
	case "ep":
		b.logger.Debug().
			Msgf("[%s] Setting up `EventPublisher` components", ContainerBuilderLogIdentifier)
		if err := b.setupEp(container); err != nil {
			b.logger.Trace().
				Err(err).
				Msgf("[%s] Error during setup of `EventPublisher` components", ContainerBuilderLogIdentifier)
			return b.logError(errors.NewApplicationSetupError("Event Publisher", err))
		}
		b.logger.Debug().
			Msgf("[%s] `EventPublisher` component setup complete", ContainerBuilderLogIdentifier)
	case "faas":
		b.logger.Debug().
			Msgf("[%s] Setting up `FunctionAsAService` components", ContainerBuilderLogIdentifier)
		if err := b.setupFaas(container); err != nil {
			b.logger.Trace().
				Err(err).
				Msgf("[%s] Error during setup of `FunctionAsAService` components", ContainerBuilderLogIdentifier)
			return b.logError(errors.NewApplicationSetupError("Function as a Service", err))
		}
		b.logger.Debug().
			Msgf("[%s] `FunctionAsAService` component setup complete", ContainerBuilderLogIdentifier)
	case "ft":
		b.logger.Debug().
			Msgf("[%s] Setting up `FunctionTrigger` components", ContainerBuilderLogIdentifier)
		if err := b.setupFt(container); err != nil {
			b.logger.Trace().
				Err(err).
				Msgf("[%s] Error during setup of `FunctionTrigger` components", ContainerBuilderLogIdentifier)
			return b.logError(errors.NewApplicationSetupError("Function Trigger", err))
		}
		b.logger.Debug().
			Msgf("[%s] `FunctionTrigger` component setup complete", ContainerBuilderLogIdentifier)
	default:
		b.logger.Trace().
			Msgf("[%s] Unknown application type: %s", ContainerBuilderLogIdentifier, b.application)
		return b.logError(errors.NewApplicationSetupError("Generic", fmt.Errorf("no application config provided")))
	}

	return nil
}
