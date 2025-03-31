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

func TestValidateScyllaConfigFailsOnEmptyHost(t *testing.T) {
	t.Parallel()
	cfg := getValidScyllaConfig()
	cfg.Host = ""

	err := cfg.Validate("scylla")
	if assert.Error(t, err, "expected error") {
		assert.Equal(t, "config value `scylla.host` must exist", err.Error())
	}
}

func TestValidateScyllaConfigFailsOnEmptyTable(t *testing.T) {
	t.Parallel()
	cfg := getValidScyllaConfig()
	cfg.Table = ""

	err := cfg.Validate("scylla")
	if assert.Error(t, err, "expected error") {
		assert.Equal(t, "config value `scylla.table` must exist", err.Error())
	}
}

func TestValidateScyllaConfigFailsOnEmptyRegion(t *testing.T) {
	t.Parallel()
	cfg := getValidScyllaConfig()
	cfg.Region = ""

	err := cfg.Validate("scylla")
	if assert.Error(t, err, "expected error") {
		assert.Equal(t, "config value `scylla.region` must exist", err.Error())
	}
}

func TestValidateScyllaConfigSucceedsOnValidConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidScyllaConfig()

	err := cfg.Validate("scylla")
	if !assert.NoError(t, err, "unexpected error") {
		return
	}
}

func getValidScyllaConfig() *config.ScyllaConfig {
	return &config.ScyllaConfig{
		Host:     "localhost",
		Port:     8000,
		Table:    "test_table",
		Region:   "eu-central-1",
		Username: "test_user",
		Password: "test_password",
		SSL:      false,
	}
}
