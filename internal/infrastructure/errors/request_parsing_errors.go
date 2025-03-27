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

package errors

import (
	"fmt"
)

func NewRequestParsingFailedError(err error) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("Error parsing request payload: %v", err),
		code:    ErrorCodeRequestParsingFailed,
	}
}

func NewMissingRequestFieldError(field string) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("missing field `%s` in request payload", field),
		code:    ErrorCodeRequestParsingMissingPayloadField,
	}
}

func NewEmptyRequestFieldError(field string) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("field `%s` in request payload must contain at least one item", field),
		code:    ErrorCodeRequestParsingEmptyPayloadField,
	}
}

func NewUnknownPayloadTypeError(payloadType string) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("Unknown payload type: `%s`", payloadType),
		code:    ErrorCodeRequestParsingUnknownPayloadType,
	}
}
