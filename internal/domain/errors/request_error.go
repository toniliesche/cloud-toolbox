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

type RequestError struct {
	message string
}

func (e RequestError) Error() string {
	return e.message
}

func NewMissingRequestFieldError(field string) error {
	return RequestError{
		message: fmt.Sprintf("missing field `%s` in request payload", field),
	}
}

func NewEmptyRequestFieldError(field string) error {
	return RequestError{
		message: fmt.Sprintf("field `%s` in request payload must contain at least one item", field),
	}
}

func NewRequestParsingError(err error) error {
	return RequestError{
		message: fmt.Sprintf("error parsing request: %s", err.Error()),
	}
}
