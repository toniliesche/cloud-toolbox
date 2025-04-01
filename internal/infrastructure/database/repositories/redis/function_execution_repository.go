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

package redis

import (
	"cloud-toolbox/internal/application/faas/models"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"time"
)

const FunctionExecutionRepositoryLogIdentifier = "FunctionExecutionRepository/Redis"

type FunctionExecutionRepository struct {
	redis      *redis.Client
	logger     *zerolog.Logger
	timeToLive time.Duration
	context    context.Context
}

func (f *FunctionExecutionRepository) GetFunction(executionId string) (*models.FunctionExecution, errors.ApplicationError) {
	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Getting function execution", FunctionExecutionRepositoryLogIdentifier)

	fn := f.redis.JSONGet(f.context, fmt.Sprintf("faas-%s", executionId), "$")
	if fn.Err() != nil {
		f.logger.Trace().
			Err(fn.Err()).
			Str("execution-id", executionId).
			Msgf("[%s] Error getting function execution", FunctionExecutionRepositoryLogIdentifier)

		return nil, errors.NewGenericError(fn.Err())
	}

	function := make([]*FunctionExecution, 0)
	f.logger.Trace().
		Str("execution-id", executionId).
		Str("json", fn.Val()).
		Msgf("[%s] Function execution retrieved", FunctionExecutionRepositoryLogIdentifier)
	if err := json.Unmarshal([]byte(fn.Val()), &function); err != nil {
		f.logger.Trace().
			Err(err).
			Str("execution-id", executionId).
			Str("json", fn.Val()).
			Msgf("[%s] Error unmarshalling function execution", FunctionExecutionRepositoryLogIdentifier)

		return nil, errors.NewGenericError(err)
	}

	if len(function) == 0 {
		f.logger.Trace().
			Str("execution-id", executionId).
			Msgf("[%s] Function execution not found", FunctionExecutionRepositoryLogIdentifier)

		return nil, errors.NewItemNotFoundError(executionId)
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Function execution found", FunctionExecutionRepositoryLogIdentifier)

	return function[0].ToModel(), nil
}

func (f *FunctionExecutionRepository) SaveError(executionId string, error string) errors.ApplicationError {
	timeObj := time.Now()

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Saving error", FunctionExecutionRepositoryLogIdentifier)

	result := f.redis.JSONSet(f.context, fmt.Sprintf("faas-%s", executionId), "$.error", fmt.Sprintf(`"%s"`, error))
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("execution-id", executionId).
			Msgf("[%s] Error saving error", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	result = f.redis.JSONSet(f.context, fmt.Sprintf("faas-%s", executionId), "$.updated", timeObj)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("execution-id", executionId).
			Msgf("[%s] Error saving error", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Error saved", FunctionExecutionRepositoryLogIdentifier)

	f.redis.Expire(f.context, fmt.Sprintf("faas-%s", executionId), f.timeToLive)

	return nil
}

func (f *FunctionExecutionRepository) SaveFunction(function *models.FunctionExecution) errors.ApplicationError {
	f.logger.Trace().
		Str("execution-id", function.Id).
		Msgf("[%s] Saving function", FunctionExecutionRepositoryLogIdentifier)

	result := f.redis.JSONSet(f.context, fmt.Sprintf("faas-%s", function.Id), "$", function)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("execution-id", function.Id).
			Msgf("[%s] Error saving function", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	f.logger.Trace().
		Str("execution-id", function.Id).
		Msgf("[%s] Function saved", FunctionExecutionRepositoryLogIdentifier)

	f.redis.Expire(f.context, fmt.Sprintf("faas-%s", function.Id), f.timeToLive)

	return nil
}

func (f *FunctionExecutionRepository) SaveOutput(executionId string, output string) errors.ApplicationError {
	timeObj := time.Now()

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Saving output", FunctionExecutionRepositoryLogIdentifier)

	var js interface{}
	if json.Unmarshal([]byte(output), &js) != nil {
		f.logger.Trace().
			Str("execution-id", executionId).
			Msgf("[%s] Output is not JSON", FunctionExecutionRepositoryLogIdentifier)

		output = fmt.Sprintf(`"%s"`, output)
	}

	result := f.redis.JSONSet(f.context, fmt.Sprintf("faas-%s", executionId), "$.output", output)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("execution-id", executionId).
			Msgf("[%s] Error saving output", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	result = f.redis.JSONSet(f.context, fmt.Sprintf("faas-%s", executionId), "$.updated", timeObj)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("execution-id", executionId).
			Msgf("[%s] Error saving output", FunctionExecutionRepositoryLogIdentifier)
		return errors.NewGenericError(result.Err())
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Output saved", FunctionExecutionRepositoryLogIdentifier)

	f.redis.Expire(f.context, fmt.Sprintf("faas-%s", executionId), f.timeToLive)

	return nil
}

func (f *FunctionExecutionRepository) UpdateStatus(executionId string, status string) errors.ApplicationError {
	timeObj := time.Now()

	f.logger.Trace().
		Str("execution-id", executionId).
		Str("status", status).
		Msgf("[%s] Updating status", FunctionExecutionRepositoryLogIdentifier)

	result := f.redis.JSONSet(f.context, fmt.Sprintf("faas-%s", executionId), "$.updated", timeObj)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("execution-id", executionId).
			Msgf("[%s] Error updating updated date", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	statusUpdate := &models.StatusUpdate{
		Status: status,
		Time:   timeObj,
	}
	appendResult := f.redis.JSONArrAppend(f.context, fmt.Sprintf("faas-%s", executionId), "$.status_updates", statusUpdate)
	if appendResult.Err() != nil {
		f.logger.Trace().
			Err(appendResult.Err()).
			Str("execution-id", executionId).
			Msgf("[%s] Error updating status list", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(appendResult.Err())
	}

	result = f.redis.JSONSet(f.context, fmt.Sprintf("faas-%s", executionId), "$.status", fmt.Sprintf(`"%s"`, status))
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("execution-id", executionId).
			Msgf("[%s] Error updating status", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Status updated", FunctionExecutionRepositoryLogIdentifier)

	f.redis.Expire(f.context, fmt.Sprintf("faas-%s", executionId), f.timeToLive)

	return nil
}

func NewFunctionExecutionRepository(container *di.Container) (*FunctionExecutionRepository, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError(FunctionExecutionRepositoryLogIdentifier)
	}

	if container.Context == nil {
		return nil, errors.NewResolveDependencyError(FunctionExecutionRepositoryLogIdentifier, "Context")
	}

	if container.Redis == nil {
		return nil, errors.NewResolveDependencyError(FunctionExecutionRepositoryLogIdentifier, "Redis")
	}

	if container.Logger == nil {
		return nil, errors.NewResolveDependencyError(FunctionExecutionRepositoryLogIdentifier, "Logger")
	}

	if container.FunctionAsAServiceConfig == nil {
		return nil, errors.NewResolveDependencyError(FunctionExecutionRepositoryLogIdentifier, "FunctionAsAServiceConfig")
	}

	if err := container.FunctionAsAServiceConfig.Validate(); err != nil {
		return nil, errors.NewInvalidConfigError(FunctionExecutionRepositoryLogIdentifier, err)
	}

	return &FunctionExecutionRepository{
		context:    container.Context,
		redis:      container.Redis,
		logger:     container.Logger,
		timeToLive: time.Duration(container.FunctionAsAServiceConfig.StorageTtl) * time.Second,
	}, nil
}
