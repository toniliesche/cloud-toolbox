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

func TestHttpServerConfigFailsOnEmptyHost(t *testing.T) {
	cfg := getValidHttpServerConfig()
	cfg.Host = ""

	err := cfg.Validate("http")
	if !assert.Error(t, err, "did not catch missing host error") {
		return
	}

	if !assert.Equal(t, "config value `http.host` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestHttpServerConfigFailsOnZeroPort(t *testing.T) {
	cfg := getValidHttpServerConfig()
	cfg.Port = 0

	err := cfg.Validate("http")
	if !assert.Error(t, err, "did not catch missing port error") {
		return
	}

	if !assert.Equal(t, "config value `http.port` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestHttpServerConfigFailsOnNegativePort(t *testing.T) {
	cfg := getValidHttpServerConfig()
	cfg.Port = -1

	err := cfg.Validate("http")
	if !assert.Error(t, err, "did not catch negative port error") {
		return
	}

	if !assert.Equal(t, "config value `http.port` must be greater than 0", err.Error(), "unexpected error message") {
		return
	}
}

func TestHttpServerConfigSucceedsOnValidConfig(t *testing.T) {
	cfg := getValidHttpServerConfig()

	err := cfg.Validate("http")
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func getValidHttpServerConfig() *config.HttpServerConfig {
	return &config.HttpServerConfig{
		Host: "localhost",
		Port: 8080,
	}
}
