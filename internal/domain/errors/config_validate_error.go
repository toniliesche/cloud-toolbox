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

type ConfigValidateError struct {
	message string
}

func (e ConfigValidateError) Error() string {
	return e.message
}

func NewConfigValueNeedsToBeGreaterZeroError(field string) error {
	return ConfigValidateError{
		message: "config value `" + field + "` must be greater than 0",
	}
}

func NewMissingConfigValueError(field string) error {
	return ConfigValidateError{
		message: "config value `" + field + "` must exist",
	}
}

func NewMissingConfigSectionError(section string) error {
	return ConfigValidateError{
		message: "config section `" + section + "` must exist",
	}
}

func NewValidateConfigSectionError(path string, err error) error {
	return ConfigValidateError{
		message: fmt.Sprintf("Error in config section `%s`: %v", path, err),
	}
}

func NewConfigValueNeedsToBeGreaterThanOrEqualValueError(field string, limit int) error {
	return ConfigValidateError{
		message: fmt.Sprintf("config value `%s` must be greater than or equal to `%d`", field, limit),
	}
}

func NewConfigValueNeedsToBeLessThanOrEqualValueError(field string, value int) error {
	return ConfigValidateError{
		message: fmt.Sprintf("config value `%s` must be less than or equal to `%d`", field, value),
	}
}
