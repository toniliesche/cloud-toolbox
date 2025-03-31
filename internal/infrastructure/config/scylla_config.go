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

const (
	scyllaServerDefaultRegion = "eu-central-1"
	scyllaServerDefaultHost   = "localhost"
	scyllaServerDefaultPort   = 8000
)

type ScyllaConfig struct {
	Region   string `yaml:"region"`
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
	Table    string `yaml:"table"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSL      bool   `yaml:"ssl"`
}

func (c *ScyllaConfig) Validate(path string) error {
	if c.Region == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.region", path))
	}

	if c.Host == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.host", path))
	}

	if c.Port < 1 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError(fmt.Sprintf("%s.port", path))
	}

	if c.Table == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.table", path))
	}

	return nil
}

func getDefaultScyllaConfig() *ScyllaConfig {
	return &ScyllaConfig{
		Region: scyllaServerDefaultRegion,
		Host:   scyllaServerDefaultHost,
		Port:   scyllaServerDefaultPort,
	}
}

func getScyllaConfigFromEnvironment() (*ScyllaConfig, errors.ApplicationError) {
	cfg := getDefaultScyllaConfig()

	region := GetEnvironmentString("SCYLLA_REGION", scyllaServerDefaultRegion)
	if region != "" {
		cfg.Region = region
	}

	host := GetEnvironmentString("SCYLLA_HOST", scyllaServerDefaultHost)
	if host != "" {
		cfg.Host = host
	}

	port, err := GetEnvironmentInt("SCYLLA_PORT", scyllaServerDefaultPort)
	if err != nil {
		return nil, err
	}

	if port != 0 {
		cfg.Port = port
	}

	table := GetEnvironmentString("SCYLLA_TABLE", "")
	if table == "" {
		return nil, errors.NewMissingEnvironmentVariableError("SCYLLA_TABLE")
	}

	cfg.Table = table

	username := GetEnvironmentString("SCYLLA_USERNAME", "")
	if username != "" {
		cfg.Username = username
	}

	password := GetEnvironmentString("SCYLLA_PASSWORD", "")
	if password != "" {
		cfg.Password = password
	}

	ssl, err := GetEnvironmentBool("SCYLLA_SSL", false)
	if err != nil {
		return nil, err
	}

	if ssl {
		cfg.SSL = ssl
	}

	return cfg, nil
}
