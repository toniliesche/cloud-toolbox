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

type RabbitMQRecord struct {
	MessageId string `json:"messageId"`
	Body      string `json:"body"`
}

func (r *RabbitMQRecord) GetRecordIdentifier() string {
	return r.MessageId
}

func (r *RabbitMQRecord) GetPayload() string {
	return r.Body
}

func (r *RabbitMQRecord) Validate() errors.ApplicationError {
	if r.MessageId == "" {
		return errors.NewMissingRequestFieldError("messageId")
	}

	if r.Body == "" {
		return errors.NewMissingRequestFieldError("body")
	}

	return nil
}
