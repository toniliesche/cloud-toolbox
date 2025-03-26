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

func TestLogConfigFailsOnMissingLevel(t *testing.T) {
	cfg := getValidLogConfig()
	cfg.Level = ""

	err := cfg.Validate("log")
	if !assert.Error(t, err, "did not catch missing level error") {
		return
	}

	if !assert.Equal(t, "config value `log.level` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestLogConfigSucceeds(t *testing.T) {
	cfg := getValidLogConfig()

	err := cfg.Validate("log")
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func getValidLogConfig() *config.LogConfig {
	return &config.LogConfig{
		Level: "info",
	}
}
