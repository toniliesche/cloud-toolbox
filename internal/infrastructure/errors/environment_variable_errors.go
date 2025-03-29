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

func NewMissingEnvironmentVariableError(variable string) ApplicationError {
	return InfrastructureError{
		message: "environment variable `" + variable + "` must exist",
		code:    ErrorCodeEnvironmentVariableMissing,
	}
}

func NewInvalidIntegerEnvironmentVariableError(variable string, value string) ApplicationError {
	return InfrastructureError{
		message: "environment variable `" + variable + "` must be an integer, got `" + value + "`",
		code:    ErrorCodeEnvironmentVariableMustBeInteger,
	}
}

func NewInvalidBooleanEnvironmentVariableError(variable string, value string) ApplicationError {
	return InfrastructureError{
		message: "environment variable `" + variable + "` must be a boolean, got `" + value + "`",
		code:    ErrorCodeEnvironmentVariableMustBeBoolean,
	}
}

func NewMalformedEnvironmentVariableError(variable string, value string, regexp string) ApplicationError {
	return InfrastructureError{
		message: "environment variable `" + variable + "` must match the regexp `" + regexp + "`, got `" + value + "`",
		code:    ErrorCodeEnvironmentVariableMalformed,
	}
}
