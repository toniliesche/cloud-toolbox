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
	"cloud-toolbox/internal/application/faas"
	"cloud-toolbox/internal/application/faas/services"
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/config/interfaces"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/http"
	"cloud-toolbox/internal/infrastructure/log"
	"fmt"
	"github.com/rs/zerolog"
)

const ContainerBuilderLogIdentifier = "ContainerBuilder"

type ContainerBuilder struct {
	application     string
	containerConfig *di.ContainerConfig
	faasConfig      *config.FunctionAsAServiceConfig
	logger          *zerolog.Logger
}

func (b *ContainerBuilder) SetFaasConfig(cfg *config.FunctionAsAServiceConfig) *ContainerBuilder {
	b.faasConfig = cfg
	b.application = "faas"

	return b
}

func (b *ContainerBuilder) Build() (*di.Container, errors.ApplicationError) {
	b.logger = log.NewTempLogger()
	b.logger.Info().
		Msgf("[%s] Starting container setup", ContainerBuilderLogIdentifier)

	container := &di.Container{}

	b.logger.Debug().
		Msgf("[%s] Setting up logger", ContainerBuilderLogIdentifier)
	if err := b.setupLogger(container); err != nil {
		return nil, b.logError(err)
	}

	b.logger = container.Logger

	switch b.application {
	case "faas":
		b.logger.Debug().
			Msgf("[%s] Setting up `FunctionAsAService` components", ContainerBuilderLogIdentifier)
		if err := b.setupFaas(container); err != nil {
			return nil, b.logError(errors.NewApplicationSetupError("Function as a Service", err))
		}
		b.logger.Debug().
			Msgf("[%s] `FunctionAsAService` component setup complete", ContainerBuilderLogIdentifier)
	default:
		return nil, b.logError(errors.NewApplicationSetupError("Generic", fmt.Errorf("no application config provided")))
	}

	b.logger.Debug().
		Msgf("[%s] Setting up basic components", ContainerBuilderLogIdentifier)
	if err := b.setupBasics(container); err != nil {
		return nil, b.logError(errors.NewApplicationSetupError("Generic", err))
	}
	b.logger.Debug().
		Msgf("[%s] Basic components setup complete", ContainerBuilderLogIdentifier)

	b.logger.Info().
		Msgf("[%s] Container setup complete", ContainerBuilderLogIdentifier)

	return container, nil
}

func (b *ContainerBuilder) logError(err errors.ApplicationError) errors.ApplicationError {
	b.logger.Error().
		Err(err).
		Msgf("[%s] Error during container setup", ContainerBuilderLogIdentifier)

	return err
}

func (b *ContainerBuilder) setupBasics(container *di.Container) errors.ApplicationError {

	return nil
}

func (b *ContainerBuilder) setupFaas(container *di.Container) errors.ApplicationError {
	var err errors.ApplicationError

	if b.faasConfig == nil {
		return errors.NewResolveDependencyError("Function as a Service", "faasConfig")
	}

	container.FunctionAsAServiceConfig = b.faasConfig

	if err = b.faasConfig.Validate(); err != nil {
		return errors.NewApplicationSetupError("Function as a Service", err)
	}

	b.logger.Trace().
		Msgf("[%s] Initializing `FunctionAsAService` registry", ContainerBuilderLogIdentifier)
	if container.FunctionRegistry, err = services.NewFunctionRegistry(container); err != nil {
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Initializing `FunctionAsAService` service", ContainerBuilderLogIdentifier)
	if container.FunctionAsAServiceService, err = faas.NewFunctionAsAServiceService(container); err != nil {
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Retrieving HttpServerConfig from `FunctionAsAService` config", ContainerBuilderLogIdentifier)
	container.HttpServerConfig = b.faasConfig.HttpServer

	b.logger.Trace().
		Msgf("[%s] Initializing `FunctionAsAService` http handler", ContainerBuilderLogIdentifier)
	if container.FunctionAsAServiceHandler, err = http.NewFunctionAsAServiceHandler(container); err != nil {
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Initializing `HttpServer`", ContainerBuilderLogIdentifier)
	if container.HttpServer, err = http.NewServer(container); err != nil {
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Registering routes for `FunctionAsAService` http handler", ContainerBuilderLogIdentifier)
	if err = container.HttpServer.RegisterRoutes(container.FunctionAsAServiceHandler); err != nil {
		return err
	}

	return nil
}

func (b *ContainerBuilder) getApplicationConfig() (interfaces.ApplicationConfig, errors.ApplicationError) {
	switch b.application {
	case "faas":
		return b.faasConfig, nil
	default:
		return nil, errors.NewContainerConfigMissingError()
	}
}

func (b *ContainerBuilder) getComponentType() string {
	switch b.application {
	case "faas":
		return "function-as-a-service"
	default:
		return "unknown"
	}
}

func (b *ContainerBuilder) setupLogger(container *di.Container) errors.ApplicationError {
	var err errors.ApplicationError

	applicationConfig, err := b.getApplicationConfig()
	if err != nil {
		return errors.NewResolveDependencyError("Generic", "applicationConfig")
	}

	container.SystemConfig = applicationConfig.GetSystemConfig()
	container.SystemConfig.ComponentType = b.getComponentType()

	if err = container.SystemConfig.Validate("system"); err != nil {
		return errors.NewApplicationSetupError("Generic", err)
	}

	container.Logger, err = log.NewLogger(container)
	if err != nil {
		return errors.NewApplicationSetupError("Generic", err)
	}

	return nil
}

func NewBuilder(cfg *di.ContainerConfig) *ContainerBuilder {
	return &ContainerBuilder{
		containerConfig: cfg,
	}
}
