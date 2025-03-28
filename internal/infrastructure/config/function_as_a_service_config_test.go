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

package config_test

import (
	"cloud-toolbox/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunctionAsAServiceConfigFailsOnMissingSystemConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.SystemConfig = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing system section error") {
		return
	}

	if !assert.Equal(t, "config section `system` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestFunctionAsAServiceConfigFailsOnInvalidSystemConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.SystemConfig.Log = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid system config error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `system`", "unexpected error message") {
		return
	}
}

func TestFunctionAsAServiceConfigFailsOnMissingHttpConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.HttpServer = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing http section error") {
		return
	}

	if !assert.Equal(t, "config section `http` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestFunctionAsAServiceConfigFailsOnInvalidHttpConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.HttpServer.Port = 0

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid http config error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `http`", "unexpected error message") {
		return
	}
}

func TestFunctionAsAServiceConfigFailsOnMissingCommand(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()
	cfg.Command = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing command error") {
		return
	}

	if !assert.Equal(t, "config value `command` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestProvideFunctionAsAServiceConfigSucceedsOnValidConfig(t *testing.T) {
	cfg := getValidFunctionAsAServiceConfig()

	err := cfg.Validate()
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func getValidFunctionAsAServiceConfig() *config.FunctionAsAServiceConfig {
	return &config.FunctionAsAServiceConfig{
		ApplicationConfig: config.ApplicationConfig{
			SystemConfig: getValidSystemConfig(),
		},
		HttpServer:        getValidHttpServerConfig(),
		Command:           "echo 'Hello, World!'",
		ParallelExecution: 1,
		ExecutionTimeout:  30,
	}
}
