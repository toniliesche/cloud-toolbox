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

package function_as_a_service

import (
	"cloud-toolbox/internal/application/faas/interfaces"
	"cloud-toolbox/internal/infrastructure/di"
	domainerrors "cloud-toolbox/internal/infrastructure/errors"
	infrastructurehttp "cloud-toolbox/internal/infrastructure/http"
	"cloud-toolbox/internal/infrastructure/http/models"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

const FunctionAsAServiceHandlerLogIdentifier = "FunctionAsAServiceHttpHandler"

type FunctionAsAServiceHandler struct {
	infrastructurehttp.Handler
	infrastructurehttp.JsonResponseHandler
	functionAsAService interfaces.FunctionAsAService
	logger             *zerolog.Logger
	functionName       string
}

func (h *FunctionAsAServiceHandler) GetRoutes() []*models.Route {
	return []*models.Route{
		{
			fmt.Sprintf("/cloud-toolbox/faas/%s", h.functionName),
			[]string{"POST"},
			h.runFunction,
		},
		{
			fmt.Sprintf("/cloud-toolbox/faas/%s/status/{executionId}", h.functionName),
			[]string{"GET"},
			h.queryFunctionStatus,
		},
	}
}

func (h *FunctionAsAServiceHandler) runFunction(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		h.SendErrorResponse(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info().
		Bytes("body", body).
		Msgf("[%s] Handling run function request", FunctionAsAServiceHandlerLogIdentifier)

	response := h.functionAsAService.RunFunction(body)

	h.SendJsonResponse(writer, request, response.GetStatusCode(), response.GetBody())
}

func (h *FunctionAsAServiceHandler) queryFunctionStatus(writer http.ResponseWriter, request *http.Request) {
	executionId, ok := h.GetPathVar(request, "executionId")
	if !ok {
		h.SendErrorResponse(writer, request, http.StatusBadRequest, "Execution ID missing")
		return
	}

	h.logger.Info().
		Str("executionId", executionId).
		Msgf("[%s] Handling get function status request", FunctionAsAServiceHandlerLogIdentifier)

	response := h.functionAsAService.GetExecutionStatus(executionId)

	h.SendJsonResponse(writer, request, response.GetStatusCode(), response.GetBody())
}

func NewFunctionAsAServiceHandler(container *di.Container) (*FunctionAsAServiceHandler, domainerrors.ApplicationError) {
	if container == nil {
		return nil, domainerrors.NewContainerMissingError(FunctionAsAServiceHandlerLogIdentifier)
	}

	if container.FunctionAsAService == nil {
		return nil, domainerrors.NewResolveDependencyError(FunctionAsAServiceHandlerLogIdentifier, "FunctionAsAService")
	}

	if container.FunctionAsAServiceConfig == nil {
		return nil, domainerrors.NewResolveDependencyError(FunctionAsAServiceHandlerLogIdentifier, "FunctionAsAServiceConfig")
	}

	if err := container.FunctionAsAServiceConfig.Validate(); err != nil {
		return nil, domainerrors.NewInvalidConfigError(FunctionAsAServiceHandlerLogIdentifier, err)
	}

	if container.Logger == nil {
		return nil, domainerrors.NewResolveDependencyError(FunctionAsAServiceHandlerLogIdentifier, "Logger")
	}

	return &FunctionAsAServiceHandler{
		JsonResponseHandler: infrastructurehttp.JsonResponseHandler{},
		functionAsAService:  container.FunctionAsAService,
		logger:              container.Logger,
		functionName:        container.FunctionAsAServiceConfig.FunctionName,
	}, nil
}
