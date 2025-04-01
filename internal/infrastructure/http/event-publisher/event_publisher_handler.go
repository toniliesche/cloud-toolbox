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

package event_publisher

import (
	"cloud-toolbox/internal/application/ep/interfaces"
	"cloud-toolbox/internal/infrastructure/di"
	domainerrors "cloud-toolbox/internal/infrastructure/errors"
	infrastructurehttp "cloud-toolbox/internal/infrastructure/http"
	httpinterfaces "cloud-toolbox/internal/infrastructure/http/interfaces"
	"cloud-toolbox/internal/infrastructure/http/models"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

const EventPublisherHandlerLogIdentifier = "EventPublisherHttpHandler"

type EventPublisherHandler struct {
	infrastructurehttp.Handler
	infrastructurehttp.JsonResponseHandler
	service            interfaces.EventPublisher
	logger             *zerolog.Logger
	eventPublisherName string
}

func (h *EventPublisherHandler) GetAuthenticator() httpinterfaces.RequestAuthenticator {
	return nil
}

func (h *EventPublisherHandler) GetRoutes() []*models.Route {
	return []*models.Route{
		{
			fmt.Sprintf("/cloud-toolbox/ep/%s", h.eventPublisherName),
			[]string{"POST"},
			h.publishEvent,
			"publish-event",
		},
	}
}

func (h *EventPublisherHandler) publishEvent(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		h.SendErrorResponse(writer, request, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info().
		Bytes("body", body).
		Msgf("[%s] Handling publish event rquest", EventPublisherHandlerLogIdentifier)

	response := h.service.PublishEvent(body)

	h.SendJsonResponse(writer, request, response.GetStatusCode(), response.GetBody())
}

func NewEventPublisherHandler(container *di.Container) (*EventPublisherHandler, domainerrors.ApplicationError) {
	if container == nil {
		return nil, domainerrors.NewContainerMissingError(EventPublisherHandlerLogIdentifier)
	}

	if container.EventPublisher == nil {
		return nil, domainerrors.NewResolveDependencyError(EventPublisherHandlerLogIdentifier, "EventPublisher")
	}

	if container.EventPublisherConfig == nil {
		return nil, domainerrors.NewResolveDependencyError(EventPublisherHandlerLogIdentifier, "EventPublisherConfig")
	}

	if err := container.EventPublisherConfig.Validate(); err != nil {
		return nil, domainerrors.NewInvalidConfigError(EventPublisherHandlerLogIdentifier, err)
	}

	if container.Logger == nil {
		return nil, domainerrors.NewResolveDependencyError(EventPublisherHandlerLogIdentifier, "Logger")
	}

	if container.EventPublisher == nil {
		return nil, domainerrors.NewResolveDependencyError(EventPublisherHandlerLogIdentifier, "EventPublisher")
	}

	return &EventPublisherHandler{
		JsonResponseHandler: infrastructurehttp.JsonResponseHandler{},
		service:             container.EventPublisher,
		logger:              container.Logger,
		eventPublisherName:  container.EventPublisherConfig.PublisherName,
	}, nil
}
