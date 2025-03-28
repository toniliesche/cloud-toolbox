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

package faas

import (
	faasapperrors "cloud-toolbox/internal/application/faas/errors"
	faasinterfaces "cloud-toolbox/internal/application/faas/interfaces"
	"cloud-toolbox/internal/application/faas/models"
	faasmodels "cloud-toolbox/internal/domain/faas/models"
	modelinterfaces "cloud-toolbox/internal/domain/models/interfaces"
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/di"
	domainerrors "cloud-toolbox/internal/infrastructure/errors"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"net/http"
)

const FunctionAsAServiceLogIdentifier = "FunctionAsAService"

type FunctionAsAService struct {
	registry faasinterfaces.FunctionRegistry
	cfg      *config.FunctionAsAServiceConfig
	logger   *zerolog.Logger
	limiter  chan uint
}

func (f *FunctionAsAService) GetExecutionStatus(executionId string) modelinterfaces.Response {
	status := f.registry.GetStatus(executionId)
	if status == "" {
		return models.NewFaasErrorResponse("", "not-found", faasapperrors.NewExecutionNotFoundError(executionId))
	}

	data := map[string]interface{}{
		"status": status,
	}

	output := f.registry.GetOutput(executionId)
	if output != "" {
		data["output"] = output
	}

	if status == "failed" || status == "timeout" || status == "rejected" {
		err := f.registry.GetError(executionId)

		if err != "" {
			data["error"] = err
		}
	}

	return &models.FaasResponse{
		Status:      http.StatusOK,
		ExecutionId: executionId,
		Data:        data,
	}
}

func (f *FunctionAsAService) RunFunction(payload []byte) modelinterfaces.Response {
	executionId := uuid.New().String()

	base := &models.FaasRequestSource{}
	f.logger.Debug().
		Msgf("[%s] Determining payload type by code", FunctionAsAServiceLogIdentifier)

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Trying to unmarshal payload to FaasRequestSource", FunctionAsAServiceLogIdentifier)
	if err := json.Unmarshal(payload, base); err != nil {
		f.logger.Warn().
			Err(err).
			Msgf("[%s] Failed to unmarshal payload to FaasRequestSource", FunctionAsAServiceLogIdentifier)

		return models.NewFaasErrorResponse("", "failed-create", domainerrors.NewRequestParsingFailedError(err))
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Trying to validate FaasRequestSource payload", FunctionAsAServiceLogIdentifier)
	if err := base.Validate(); err != nil {
		f.logger.Warn().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Failed to validate FaasRequestSource payload", FunctionAsAServiceLogIdentifier)

		return models.NewFaasErrorResponse("", "failed-create", err)
	}

	f.logger.Debug().
		Str("execution-id", executionId).
		Str("payload-type", base.Source).
		Msgf("[%s] Payload type has been determined", FunctionAsAServiceLogIdentifier)

	switch base.Source {
	case "ctb:rmq":
		return f.runRabbitMQSourceFunction(executionId, payload)
	default:
		return models.NewFaasErrorResponse("", "failed-create", domainerrors.NewUnknownPayloadTypeError(base.Source))
	}
}

func (f *FunctionAsAService) runRabbitMQSourceFunction(executionId string, payload []byte) modelinterfaces.Response {
	f.logger.Debug().
		Str("execution-id", executionId).
		Msgf("[%s] Running RabbitMQ sourced function", FunctionAsAServiceLogIdentifier)

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Trying to unmarshal payload to FaasRequest", FunctionAsAServiceLogIdentifier)

	request := &models.FaasRequest[*models.RabbitMQRecord]{}
	if err := json.Unmarshal(payload, request); err != nil {
		f.logger.Warn().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Failed to unmarshal payload to FaasRequest", FunctionAsAServiceLogIdentifier)

		return models.NewFaasErrorResponse("", "failed-create", domainerrors.NewRequestParsingFailedError(err))
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Trying to validate FaasRequest", FunctionAsAServiceLogIdentifier)

	if err := request.Validate(); err != nil {
		f.logger.Warn().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Failed to validate FaasRequest", FunctionAsAServiceLogIdentifier)

		return models.NewFaasErrorResponse("", "failed-create", err)
	}

	return f.runFunction(executionId, request)
}

func (f *FunctionAsAService) runFunction(executionId string, request *models.FaasRequest[*models.RabbitMQRecord]) modelinterfaces.Response {
	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Trying to marshal records to JSON", FunctionAsAServiceLogIdentifier)

	recordsJson, jsonErr := json.Marshal(request.Records)
	if jsonErr != nil {
		f.logger.Warn().
			Err(jsonErr).
			Str("execution-id", executionId).
			Msgf("[%s] Failed to marshal records to JSON", FunctionAsAServiceLogIdentifier)

		return models.NewFaasErrorResponse("", "failed-create", domainerrors.NewRequestParsingFailedError(jsonErr))
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Creating function execution", FunctionAsAServiceLogIdentifier)

	functionExecution, err := faasmodels.NewFunctionExecution(executionId, f.cfg.Command, f.cfg.ExecutionTimeout, string(recordsJson), f.logger)
	if err != nil {
		f.logger.Error().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Failed to create function execution", FunctionAsAServiceLogIdentifier)

		return models.NewFaasErrorResponse("", "failed-create", err)
	}

	convertedRequest := request.ToAbstractRequest()
	err = f.registry.AddRecord(functionExecution, &convertedRequest)
	if err != nil {
		f.logger.Error().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Failed to add function execution to registry", FunctionAsAServiceLogIdentifier)

		return models.NewFaasErrorResponse(executionId, "failed-create", err)
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Running function execution", FunctionAsAServiceLogIdentifier)

	err = functionExecution.Run(&f.limiter)
	if err != nil {
		f.logger.Error().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Failed to run function execution", FunctionAsAServiceLogIdentifier)

		fmt.Printf("%T\n", err)

		if err.Code() == domainerrors.ErrorCodeFaasLimitExceeded {
			f.logger.Trace().
				Str("execution-id", executionId).
				Msgf("[%s] Function execution has been rejected", FunctionAsAServiceLogIdentifier)

			return models.NewFaasErrorResponse(executionId, "rejected", err)
		}

		f.logger.Trace().
			Str("execution-id", executionId).
			Msgf("[%s] Function execution has failed", FunctionAsAServiceLogIdentifier)

		return models.NewFaasErrorResponse(executionId, "failed-start", err)
	}

	if !request.Async {
		f.logger.Trace().
			Str("execution-id", executionId).
			Msgf("[%s] Waiting for function execution to finish", FunctionAsAServiceLogIdentifier)

		err := functionExecution.Wait()

		f.logger.Debug().
			Str("execution-id", executionId).
			Msgf("[%s] Function execution has finished", FunctionAsAServiceLogIdentifier)

		var httpStatus int
		status := functionExecution.GetStatus()

		output := f.registry.GetOutput(executionId)
		if status == "success" {
			f.logger.Debug().
				Str("execution-id", executionId).
				Msgf("[%s] Function execution has been successful", FunctionAsAServiceLogIdentifier)

			httpStatus = http.StatusOK
		} else {
			f.logger.Debug().
				Str("execution-id", executionId).
				Str("status", status).
				Msgf("[%s] Function execution has failed", FunctionAsAServiceLogIdentifier)

			httpStatus = http.StatusInternalServerError
		}

		if err != nil {
			return models.NewFaasErrorResponse(executionId, functionExecution.GetStatus(), err)
		}

		data := map[string]interface{}{
			"status": status,
		}

		if output != "" {
			data["output"] = output
		}

		return &models.FaasResponse{
			Status:      httpStatus,
			ExecutionId: executionId,
			Data:        data,
		}
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Function execution has been accepted", FunctionAsAServiceLogIdentifier)

	data := map[string]interface{}{
		"status": "created",
	}

	return &models.FaasResponse{
		Status:      http.StatusAccepted,
		ExecutionId: executionId,
		Data:        data,
	}
}

func NewFunctionAsAServiceService(container *di.Container) (faasinterfaces.FunctionAsAService, domainerrors.ApplicationError) {
	if container == nil {
		return nil, domainerrors.NewContainerMissingError("FunctionAsAServiceService")
	}

	if container.FunctionAsAServiceConfig == nil {
		return nil, domainerrors.NewResolveDependencyError("FunctionAsAServiceService", "FunctionAsAServiceConfig")
	}

	if container.Logger == nil {
		return nil, domainerrors.NewResolveDependencyError("FunctionAsAServiceService", "Logger")
	}

	if container.FunctionRegistry == nil {
		return nil, domainerrors.NewResolveDependencyError("FunctionAsAServiceService", "FunctionRegistry")
	}

	return &FunctionAsAService{
		registry: container.FunctionRegistry,
		cfg:      container.FunctionAsAServiceConfig,
		logger:   container.Logger,
		limiter:  make(chan uint, container.FunctionAsAServiceConfig.ParallelExecution),
	}, nil
}
