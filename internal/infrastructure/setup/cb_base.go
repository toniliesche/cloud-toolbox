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
	"cloud-toolbox/internal/infrastructure/config/interfaces"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/log"
)

func (b *ContainerBuilder) setupBasics(container *di.Container) errors.ApplicationError {

	return nil
}

func (b *ContainerBuilder) setupLogger(container *di.Container) errors.ApplicationError {
	var applicationConfig interfaces.ApplicationConfig
	var err errors.ApplicationError

	b.logger.Trace().
		Msgf("[%s] Retrieving application config", ContainerBuilderLogIdentifier)
	if applicationConfig, err = b.getApplicationConfig(); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error retrieving application config", ContainerBuilderLogIdentifier)
		return errors.NewResolveDependencyError("Generic", "applicationConfig")
	}

	b.logger.Trace().
		Msgf("[%s] Retrieving SystemConfig from application config", ContainerBuilderLogIdentifier)
	container.SystemConfig = applicationConfig.GetSystemConfig()

	if container.SystemConfig == nil {
		b.logger.Trace().
			Msgf("[%s] SystemConfig is nil", ContainerBuilderLogIdentifier)
		return errors.NewResolveDependencyError("Generic", "SystemConfig")
	}

	b.logger.Trace().
		Msgf("[%s] Validating SystemConfig", ContainerBuilderLogIdentifier)
	if err = container.SystemConfig.Validate("system"); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] SystemConfig is invalid", ContainerBuilderLogIdentifier)
		return errors.NewApplicationSetupError("Generic", err)
	}

	container.SystemConfig.ComponentType = b.getComponentType()

	b.logger.Trace().
		Msgf("[%s] Initializing logger", ContainerBuilderLogIdentifier)
	if container.Logger, err = log.NewLogger(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing logger", ContainerBuilderLogIdentifier)
		return errors.NewApplicationSetupError("Generic", err)
	}

	b.logger.Trace().
		Msgf("[%s] Logger initialized", ContainerBuilderLogIdentifier)

	return nil
}
