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
	"cloud-toolbox/internal/infrastructure/database/connectors"
	"cloud-toolbox/internal/infrastructure/database/repositories/inmemory"
	"cloud-toolbox/internal/infrastructure/database/repositories/redis"
	"cloud-toolbox/internal/infrastructure/database/repositories/scylla"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/http"
)

func (b *ContainerBuilder) setupFaas(container *di.Container) errors.ApplicationError {
	var err errors.ApplicationError

	if b.faasConfig == nil {
		b.logger.Trace().
			Msgf("[%s] `FunctionAsAService` config is nil", ContainerBuilderLogIdentifier)
		return errors.NewResolveDependencyError("Function as a Service", "faasConfig")
	}

	b.logger.Trace().
		Msgf("[%s] Validating `FunctionAsAService` config", ContainerBuilderLogIdentifier)
	if err = b.faasConfig.Validate(); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] `FunctionAsAService` config is invalid", ContainerBuilderLogIdentifier)
		return errors.NewApplicationSetupError("Function as a Service", err)
	}
	container.FunctionAsAServiceConfig = b.faasConfig

	if b.faasConfig.StorageBackend == "redis" {
		b.logger.Trace().
			Msgf("[%s] Retrieving RedisConfig from `FunctionAsAService` config", ContainerBuilderLogIdentifier)
		container.RedisConfig = b.faasConfig.Redis

		b.logger.Trace().
			Msgf("[%s] Initializing `Redis` storage backend", ContainerBuilderLogIdentifier)
		if container.Redis, err = connectors.NewRedis(container); err != nil {
			return err
		}

		b.logger.Trace().
			Msgf("[%s] Initializing `FunctionExecutionRepository` with redis backend", ContainerBuilderLogIdentifier)
		if container.FunctionExecutionRepository, err = redis.NewFunctionExecutionRepository(container); err != nil {
			return err
		}
	} else if b.faasConfig.StorageBackend == "scylla" {
		b.logger.Trace().
			Msgf("[%s] Retrieving ScyllaConfig from `FunctionAsAService` config", ContainerBuilderLogIdentifier)
		container.ScyllaConfig = b.faasConfig.Scylla

		b.logger.Trace().
			Msgf("[%s] Initializing `Scylla` storage backend", ContainerBuilderLogIdentifier)
		if container.Scylla, err = connectors.NewScylla(container); err != nil {
			return err
		}

		b.logger.Trace().
			Msgf("[%s] Initializing `FunctionExecutionRepository` with scylla backend", ContainerBuilderLogIdentifier)
		if container.FunctionExecutionRepository, err = scylla.NewFunctionExecutionRepository(container); err != nil {
			b.logger.Trace().
				Err(err).
				Msgf("[%s] Error initializing `FunctionExecutionRepository` with scylla backend", ContainerBuilderLogIdentifier)
			return err
		}
	} else {
		b.logger.Trace().
			Msgf("[%s] Initializing `FunctionExecutionRepository` with storage backend", ContainerBuilderLogIdentifier)
		if container.FunctionExecutionRepository, err = inmemory.NewFunctionExecutionRepository(container); err != nil {
			b.logger.Trace().
				Err(err).
				Msgf("[%s] Error initializing `FunctionExecutionRepository` with inmemory backend", ContainerBuilderLogIdentifier)
			return err
		}
	}

	b.logger.Trace().
		Msgf("[%s] Initializing `FunctionAsAService` registry", ContainerBuilderLogIdentifier)
	if container.FunctionRegistry, err = services.NewFunctionRegistry(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `FunctionRegistry`", ContainerBuilderLogIdentifier)
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Initializing `FunctionAsAService` service", ContainerBuilderLogIdentifier)
	if container.FunctionAsAServiceService, err = faas.NewFunctionAsAServiceService(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `FunctionAsAService` service", ContainerBuilderLogIdentifier)
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Retrieving HttpServerConfig from `FunctionAsAService` config", ContainerBuilderLogIdentifier)
	container.HttpServerConfig = b.faasConfig.HttpServer

	b.logger.Trace().
		Msgf("[%s] Initializing `HttpServer`", ContainerBuilderLogIdentifier)
	if container.HttpServer, err = http.NewServer(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `HttpServer`", ContainerBuilderLogIdentifier)
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Initializing `FunctionAsAService` http handler", ContainerBuilderLogIdentifier)
	if container.FunctionAsAServiceHandler, err = http.NewFunctionAsAServiceHandler(container); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error initializing `FunctionAsAService` http handler", ContainerBuilderLogIdentifier)
		return err
	}

	b.logger.Trace().
		Msgf("[%s] Registering routes for `FunctionAsAService` http handler", ContainerBuilderLogIdentifier)
	if err = container.HttpServer.RegisterRoutes(container.FunctionAsAServiceHandler); err != nil {
		b.logger.Trace().
			Err(err).
			Msgf("[%s] Error registering routes for `FunctionAsAService` http handler", ContainerBuilderLogIdentifier)
		return err
	}

	b.logger.Trace().
		Msgf("[%s] `FunctionAsAService` setup completed", ContainerBuilderLogIdentifier)

	return nil
}
