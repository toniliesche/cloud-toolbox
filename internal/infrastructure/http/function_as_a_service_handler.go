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

package http

import (
	"cloud-toolbox/internal/application/faas/interfaces"
	"cloud-toolbox/internal/infrastructure/di"
	domainerrors "cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/http/models"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

const FunctionAsAServiceHandlerLogIdentifier = "FunctionAsAServiceHandler"

type FunctionAsAServiceHandler struct {
	Handler
	JsonResponseHandler
	service interfaces.FunctionAsAService
	logger  *zerolog.Logger
}

func (h *FunctionAsAServiceHandler) GetRoutes() []*models.Route {
	return []*models.Route{
		{
			"/cloud-toolbox/faas",
			[]string{"POST"},
			h.runFunction,
		},
		{
			"/cloud-toolbox/faas/status/{executionId}",
			[]string{"GET"},
			h.queryFunctionStatus,
		},
	}
}

func (h *FunctionAsAServiceHandler) runFunction(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		h.sendErrorResponse(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info().
		Bytes("body", body).
		Msgf("[%s] Handling run function request", FunctionAsAServiceHandlerLogIdentifier)

	response := h.service.RunFunction(body)

	h.sendJsonResponse(writer, request, response.GetStatusCode(), response.GetBody())
}

func (h *FunctionAsAServiceHandler) queryFunctionStatus(writer http.ResponseWriter, request *http.Request) {
	executionId, ok := h.GetPathVar(request, "executionId")
	if !ok {
		h.sendErrorResponse(writer, request, http.StatusBadRequest, "Execution ID missing")
		return
	}

	h.logger.Info().
		Str("executionId", executionId).
		Msgf("[%s] Handling get function status request", FunctionAsAServiceHandlerLogIdentifier)

	response := h.service.GetExecutionStatus(executionId)

	h.sendJsonResponse(writer, request, response.GetStatusCode(), response.GetBody())
}

func NewFunctionAsAServiceHandler(container *di.Container) (*FunctionAsAServiceHandler, domainerrors.ApplicationError) {
	if container == nil {
		return nil, domainerrors.NewContainerMissingError("FunctionAsAServiceHandler")
	}

	if container.FunctionAsAServiceService == nil {
		return nil, domainerrors.NewResolveDependencyError("FunctionAsAServiceHandler", "FunctionAsAServiceService")
	}

	if container.Logger == nil {
		return nil, domainerrors.NewResolveDependencyError("FunctionAsAServiceHandler", "Logger")
	}

	return &FunctionAsAServiceHandler{
		JsonResponseHandler: JsonResponseHandler{},
		service:             container.FunctionAsAServiceService,
		logger:              container.Logger,
	}, nil
}
