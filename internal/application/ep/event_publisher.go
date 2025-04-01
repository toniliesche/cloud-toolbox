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

package ep

import (
	epinterfaces "cloud-toolbox/internal/application/ep/interfaces"
	"cloud-toolbox/internal/application/ep/models"
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	httpinterfaces "cloud-toolbox/internal/infrastructure/http/models/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"net/http"
)

const EventPublisherLogIdentifier = "EventPublisher"

type EventPublisher struct {
	backend epinterfaces.PublisherBackend
	cfg     *config.EventPublisherConfig
	logger  *zerolog.Logger
	context context.Context
}

func (p *EventPublisher) PublishEvent(payload []byte) httpinterfaces.Response {
	eventId := fmt.Sprintf("ctb:eb:%s:%s", p.cfg.PublisherName, uuid.New().String())

	p.logger.Trace().
		Str("event-id", eventId).
		Msgf("[%s] Trying to unmarshal payload: to EpRequest", EventPublisherLogIdentifier)

	event := &models.EpRequest{}
	if err := json.Unmarshal(payload, event); err != nil {
		p.logger.Warn().
			Err(err).
			Str("event-id", eventId).
			Msgf("[%s] Failed to unmarshal payload to FaasRequest", EventPublisherLogIdentifier)

		return models.NewEpErrorResponse("", "failed-send", errors.NewRequestParsingFailedError(err))
	}

	p.logger.Trace().
		Str("event-id", eventId).
		Msgf("[%s] Trying to validate FaasRequest", EventPublisherLogIdentifier)

	if err := event.Validate(); err != nil {
		p.logger.Warn().
			Err(err).
			Str("event-id", eventId).
			Msgf("[%s] Failed to validate FaasRequest", EventPublisherLogIdentifier)

		return models.NewEpErrorResponse("", "failed-send", err)
	}

	if err := p.backend.PublishEvent(event); err != nil {
		p.logger.Warn().
			Err(err).
			Str("event-id", eventId).
			Msgf("[%s] Failed to publish event", EventPublisherLogIdentifier)

		return models.NewEpErrorResponse("", "failed-send", err)
	}

	return models.NewEpResponse(eventId, http.StatusOK, map[string]string{})
}

func NewEventPublisher(container *di.Container) (*EventPublisher, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError(EventPublisherLogIdentifier)
	}

	if container.Context == nil {
		return nil, errors.NewResolveDependencyError(EventPublisherLogIdentifier, "Context")
	}

	if container.EventPublisherConfig == nil {
		return nil, errors.NewResolveDependencyError(EventPublisherLogIdentifier, "EventPublisherConfig")
	}

	if err := container.EventPublisherConfig.Validate(); err != nil {
		return nil, errors.NewInvalidConfigError(EventPublisherLogIdentifier, err)
	}

	if container.Logger == nil {
		return nil, errors.NewResolveDependencyError(EventPublisherLogIdentifier, "Logger")
	}

	if container.EventPublisherBackend == nil {
		return nil, errors.NewResolveDependencyError(EventPublisherLogIdentifier, "EventPublisherBackend")
	}

	return &EventPublisher{
		backend: container.EventPublisherBackend,
		context: container.Context,
		logger:  container.Logger,
		cfg:     container.EventPublisherConfig,
	}, nil
}
