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

import "fmt"

type InvalidPayloadError struct {
	message string
}

func (e InvalidPayloadError) Error() string {
	return e.message
}

func NewMissingPayloadFieldError(field string) error {
	return InvalidPayloadError{
		message: fmt.Sprintf("Request payload is missing required field: `%s`", field),
	}
}

func NewUnknownPayloadTypeError(payloadType string) error {
	return InvalidPayloadError{
		message: fmt.Sprintf("Unknown payload type: `%s`", payloadType),
	}
}

func NewPayloadParsingError(err error) error {
	return InvalidPayloadError{
		message: fmt.Sprintf("Error parsing request payload: %v", err),
	}
}
