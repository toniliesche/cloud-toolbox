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

func NewConfigCreationError(err error) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("Error while creating config: %v", err),
		code:    ErrorCodeConfigCreationFailedError,
	}
}

func NewConfigValidationError(err error) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("Error while validating config: %v", err),
		code:    ErrorCodeConfigValidationFailedError,
	}
}

func NewConfigFileReadingError(err error) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("Error while reading config file: %v", err),
		code:    ErrorCodeConfigFileReadingFailedError,
	}
}

func NewConfigFileParsingError(err error) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("Error while parsing config file: %v", err),
		code:    ErrorCodeConfigFileParsingFailedError,
	}
}

func NewConfigValueNeedsToBeGreaterZeroError(field string) ApplicationError {
	return InfrastructureError{
		message: "config value `" + field + "` must be greater than `0`",
		code:    ErrorCodeConfigValueMustBeGreaterThanZero,
	}
}

func NewMissingConfigValueError(field string) ApplicationError {
	return InfrastructureError{
		message: "config value `" + field + "` must exist",
		code:    ErrorCodeConfigValueMissing,
	}
}

func NewEmptyListConfigValueError(field string) ApplicationError {
	return InfrastructureError{
		message: "config value `" + field + "` must not be an empty list",
		code:    ErrorCodeConfigValueEmptyList,
	}
}

func NewMissingListConfigValueError(field string, id string) ApplicationError {
	return InfrastructureError{
		message: "config value `" + field + "` does not have an entry with id `" + id + "`",
		code:    ErrorCodeConfigValueMissingListEntry,
	}
}

func NewMissingConfigSectionError(section string) ApplicationError {
	return InfrastructureError{
		message: "config section `" + section + "` must exist",
		code:    ErrorCodeConfigSectionMissing,
	}
}

func NewValidateConfigSectionError(path string, err error) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("Error in config section `%s`: %v", path, err),
		code:    ErrorCodeConfigSectionInvalid,
	}
}

func NewConfigValueNeedsToBeGreaterThanOrEqualValueError(field string, limit int) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("config value `%s` must be greater than or equal to `%d`", field, limit),
		code:    ErrorCodeConfigValueMustBeGreaterThanOrEqual,
	}
}

func NewConfigValueNeedsToBeLessThanOrEqualValueError(field string, value int) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("config value `%s` must be less than or equal to `%d`", field, value),
		code:    ErrorCodeConfigValueMustBeLessThanOrEqual,
	}
}

func NewInvalidConfigValueError(field string, allowedValues []string, value string) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("config value `%s` must be one of %v, but is `%s`", field, allowedValues, value),
		code:    ErrorCodeConfigValueInvalid,
	}
}
