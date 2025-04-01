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
	"cloud-toolbox/internal/infrastructure/log"
	"context"
	"github.com/rs/zerolog"
)

const ContainerBuilderLogIdentifier = "ContainerBuilder"

type ContainerBuilder struct {
	application string
	epConfig    *config.EventPublisherConfig
	faasConfig  *config.FunctionAsAServiceConfig
	ftConfig    *config.FunctionTriggerConfig
	logger      *zerolog.Logger
	context     context.Context
}

func (b *ContainerBuilder) Build() (*di.Container, errors.ApplicationError) {
	b.logger = log.NewTempLogger()
	b.logger.Info().
		Msgf("[%s] Starting container setup", ContainerBuilderLogIdentifier)

	if b.context == nil {
		b.logger.Error().
			Msgf("[%s] Context is not set", ContainerBuilderLogIdentifier)
		return nil, b.logError(errors.NewContainerConfigMissingError())
	}

	container := &di.Container{}
	container.Context = b.context

	b.logger.Debug().
		Msgf("[%s] Setting up logger", ContainerBuilderLogIdentifier)
	if err := b.setupLogger(container); err != nil {
		return nil, b.logError(err)
	}

	b.logger = container.Logger

	err := b.setupApplication(container)
	if err != nil {
		return nil, err
	}

	b.logger.Debug().
		Msgf("[%s] Setting up basic components", ContainerBuilderLogIdentifier)
	if err := b.setupBasics(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error during setup of basic components", ContainerBuilderLogIdentifier)
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
