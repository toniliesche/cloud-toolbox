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
	"fmt"
)

type LogConfig struct {
	DevMode bool   `yaml:"dev_mode"`
	Level   string `yaml:"level"`
	Path    string `yaml:"path"`
}

func (c *LogConfig) Validate(path string) errors.ApplicationError {
	if c.Level == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.level", path))
	}

	if c.Path == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.path", path))
	}

	return nil
}

func getLogConfigFromEnvironment() (*LogConfig, errors.ApplicationError) {
	devMode, err := GetEnvironmentBool("LOG_DEV_MODE", false)
	if err != nil {
		return nil, err
	}

	return &LogConfig{
		DevMode: devMode,
		Path:    GetEnvironmentString("LOG_PATH", "/dev/stdout"),
		Level:   GetEnvironmentString("LOG_LEVEL", "info"),
	}, nil
}
