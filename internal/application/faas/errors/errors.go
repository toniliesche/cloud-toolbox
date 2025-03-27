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
	"cloud-toolbox/internal/infrastructure/errors"
	"fmt"
)

type FaasError struct {
	message string
	code    int
}

func (e FaasError) Error() string {
	return e.message
}

func (e FaasError) Code() int {
	return e.code
}

func NewExecutionFailedError(err error) errors.ApplicationError {
	return FaasError{
		message: fmt.Sprintf("Execution failed: %s", err.Error()),
		code:    errors.ErrorCodeFaasExecutionFailed,
	}
}

func NewExecutionTimeoutError(err error) errors.ApplicationError {
	return FaasError{
		message: fmt.Sprintf("Execution timed out: %s", err.Error()),
		code:    errors.ErrorCodeFaasExecutionTimedOut,
	}
}

func NewFaasCommandCouldNotBeParsedError(err error) errors.ApplicationError {
	return FaasError{
		message: fmt.Sprintf("Command could not be parsed: %s", err.Error()),
		code:    errors.ErrorCodeFaasCommandCouldNotBeParsed,
	}
}

func NewFaasCommandCouldNotBeKilledError(err error) errors.ApplicationError {
	return FaasError{
		message: fmt.Sprintf("Command could not be killed: %s", err.Error()),
		code:    errors.ErrorCodeFaasCommandCouldNotBeKilled,
	}
}

func NewExecutionNotFoundError(executionId string) errors.ApplicationError {
	return FaasError{
		message: fmt.Sprintf("Execution not found: `%s`", executionId),
		code:    errors.ErrorCodeFaasExecutionNotFound,
	}
}

func NewLimitExceededError() errors.ApplicationError {
	return FaasError{
		message: "Limit of parallel executions exceeded",
		code:    errors.ErrorCodeFaasLimitExceeded,
	}
}

func NewNotStartedError() errors.ApplicationError {
	return FaasError{
		message: "Function not started, yet",
	}
}

func NewAlreadyStartedError() errors.ApplicationError {
	return FaasError{
		message: "Function already started",
	}
}

func NewAlreadyFinishedError() errors.ApplicationError {
	return FaasError{
		message: "Function already finished",
	}
}

func NewNotFinishedError() errors.ApplicationError {
	return FaasError{
		message: "Function not finished, yet",
	}
}
