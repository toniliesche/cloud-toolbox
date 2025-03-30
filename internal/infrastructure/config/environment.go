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

package config

import (
	"cloud-toolbox/internal/infrastructure/errors"
	"os"
	"strconv"
)

func GetEnvironmentString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetEnvironmentInt(key string, defaultValue int64) (int64, errors.ApplicationError) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.NewInvalidIntegerEnvironmentVariableError(key, value)
	}

	return intValue, nil
}

func HasEnvironment(key string) bool {
	value := os.Getenv(key)
	return value != ""
}

func GetEnvironmentBool(key string, defaultValue bool) (bool, errors.ApplicationError) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, errors.NewInvalidBooleanEnvironmentVariableError(key, value)
	}

	return boolValue, nil
}
