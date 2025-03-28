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
	redis  *redis.Client
	logger *zerolog.Logger
}

func (f *FunctionExecutionRepository) GetFunction(executionId string) (*models.FunctionExecution, errors.ApplicationError) {
	f.logger.Trace().
		Str("executionId", executionId).
		Msgf("[%s] Getting function execution", FunctionExecutionRepositoryLogIdentifier)

	fn := f.redis.JSONGet(context.Background(), fmt.Sprintf("faas-%s", executionId), "$")
	if fn.Err() != nil {
		f.logger.Trace().
			Err(fn.Err()).
			Str("executionId", executionId).
			Msgf("[%s] Error getting function execution", FunctionExecutionRepositoryLogIdentifier)

		return nil, errors.NewGenericError(fn.Err())
	}

	function := make([]*FunctionExecution, 0)
	f.logger.Trace().
		Str("executionId", executionId).
		Str("json", fn.Val()).
		Msgf("[%s] Function execution retrieved", FunctionExecutionRepositoryLogIdentifier)
	if err := json.Unmarshal([]byte(fn.Val()), &function); err != nil {
		f.logger.Trace().
			Err(err).
			Str("executionId", executionId).
			Str("json", fn.Val()).
			Msgf("[%s] Error unmarshalling function execution", FunctionExecutionRepositoryLogIdentifier)

		return nil, errors.NewGenericError(err)
	}

	if len(function) == 0 {
		f.logger.Trace().
			Str("executionId", executionId).
			Msgf("[%s] Function execution not found", FunctionExecutionRepositoryLogIdentifier)

		return nil, errors.NewItemNotFoundError(executionId)
	}

	f.logger.Trace().
		Str("executionId", executionId).
		Msgf("[%s] Function execution found", FunctionExecutionRepositoryLogIdentifier)

	return function[0].ToModel(), nil
}

func (f *FunctionExecutionRepository) SaveError(executionId string, error string) errors.ApplicationError {
	timeObj := time.Now()

	f.logger.Trace().
		Str("executionId", executionId).
		Msgf("[%s] Saving error", FunctionExecutionRepositoryLogIdentifier)

	result := f.redis.JSONSet(context.Background(), fmt.Sprintf("faas-%s", executionId), "$.error", fmt.Sprintf(`"%s"`, error))
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("executionId", executionId).
			Msgf("[%s] Error saving error", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	result = f.redis.JSONSet(context.Background(), fmt.Sprintf("faas-%s", executionId), "$.updated", timeObj)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("executionId", executionId).
			Msgf("[%s] Error saving error", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	f.logger.Trace().
		Str("executionId", executionId).
		Msgf("[%s] Error saved", FunctionExecutionRepositoryLogIdentifier)

	return nil
}

func (f *FunctionExecutionRepository) SaveFunction(function *models.FunctionExecution) errors.ApplicationError {
	f.logger.Trace().
		Str("executionId", function.Id).
		Msgf("[%s] Saving function", FunctionExecutionRepositoryLogIdentifier)

	result := f.redis.JSONSet(context.Background(), fmt.Sprintf("faas-%s", function.Id), "$", function)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("executionId", function.Id).
			Msgf("[%s] Error saving function", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	f.logger.Trace().
		Str("executionId", function.Id).
		Msgf("[%s] Function saved", FunctionExecutionRepositoryLogIdentifier)

	return nil
}

func (f *FunctionExecutionRepository) SaveOutput(executionId string, output string) errors.ApplicationError {
	timeObj := time.Now()

	f.logger.Trace().
		Str("executionId", executionId).
		Msgf("[%s] Saving output", FunctionExecutionRepositoryLogIdentifier)

	var js interface{}
	if json.Unmarshal([]byte(output), &js) != nil {
		f.logger.Trace().
			Str("executionId", executionId).
			Msgf("[%s] Output is not JSON", FunctionExecutionRepositoryLogIdentifier)

		output = fmt.Sprintf(`"%s"`, output)
	}

	result := f.redis.JSONSet(context.Background(), fmt.Sprintf("faas-%s", executionId), "$.output", output)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("executionId", executionId).
			Msgf("[%s] Error saving output", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	result = f.redis.JSONSet(context.Background(), fmt.Sprintf("faas-%s", executionId), "$.updated", timeObj)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("executionId", executionId).
			Msgf("[%s] Error saving output", FunctionExecutionRepositoryLogIdentifier)
		return errors.NewGenericError(result.Err())
	}

	f.logger.Trace().
		Str("executionId", executionId).
		Msgf("[%s] Output saved", FunctionExecutionRepositoryLogIdentifier)

	return nil
}

func (f *FunctionExecutionRepository) UpdateStatus(executionId string, status string) errors.ApplicationError {
	timeObj := time.Now()

	f.logger.Trace().
		Str("executionId", executionId).
		Str("status", status).
		Msgf("[%s] Updating status", FunctionExecutionRepositoryLogIdentifier)

	result := f.redis.JSONSet(context.Background(), fmt.Sprintf("faas-%s", executionId), "$.updated", timeObj)
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("executionId", executionId).
			Msgf("[%s] Error updating updated date", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	statusUpdate := &models.StatusUpdate{
		Status: status,
		Time:   timeObj,
	}
	appendResult := f.redis.JSONArrAppend(context.Background(), fmt.Sprintf("faas-%s", executionId), "$.status_updates", statusUpdate)
	if appendResult.Err() != nil {
		f.logger.Trace().
			Err(appendResult.Err()).
			Str("executionId", executionId).
			Msgf("[%s] Error updating status list", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(appendResult.Err())
	}

	result = f.redis.JSONSet(context.Background(), fmt.Sprintf("faas-%s", executionId), "$.status", fmt.Sprintf(`"%s"`, status))
	if result.Err() != nil {
		f.logger.Trace().
			Err(result.Err()).
			Str("executionId", executionId).
			Msgf("[%s] Error updating status", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(result.Err())
	}

	f.logger.Trace().
		Str("executionId", executionId).
		Msgf("[%s] Status updated", FunctionExecutionRepositoryLogIdentifier)

	return nil
}

func NewFunctionExecutionRepository(container *di.Container) (*FunctionExecutionRepository, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError("FunctionExecutionRepository")
	}

	if container.Redis == nil {
		return nil, errors.NewResolveDependencyError("FunctionExecutionRepository", "Redis")
	}

	if container.Logger == nil {
		return nil, errors.NewResolveDependencyError("FunctionExecutionRepository", "Logger")
	}

	return &FunctionExecutionRepository{
		redis:  container.Redis,
		logger: container.Logger,
	}, nil
}
