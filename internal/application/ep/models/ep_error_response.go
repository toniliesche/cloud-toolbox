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

package models

import (
	"cloud-toolbox/internal/infrastructure/errors"
)

type EpErrorResponse struct {
	EventId string
	Status  string
	Err     errors.ApplicationError
}

func (r *EpErrorResponse) GetStatusCode() int {
	return errors.MapToStatusCode(r.Err.Code())
}

func (r *EpErrorResponse) GetBody() interface{} {
	body := map[string]interface{}{
		"status": r.Status,
		"error":  r.Err.Error(),
	}

	if r.EventId != "" {
		body["event_id"] = r.EventId
	}

	return body
}

func NewEpErrorResponse(eventId string, errorCode string, err errors.ApplicationError) *EpErrorResponse {
	return &EpErrorResponse{
		EventId: eventId,
		Status:  errorCode,
		Err:     err,
	}
}
