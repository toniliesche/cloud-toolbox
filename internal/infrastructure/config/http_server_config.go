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
	httpServerDefaultHost = "0.0.0.0"
	httpServerDefaultPort = 8080
)

type HttpServerConfig struct {
	Host string `yaml:"host"`
	Port int64  `yaml:"port"`
}

func (c *HttpServerConfig) Validate(path string) errors.ApplicationError {
	if c.Host == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.host", path))
	}

	if c.Port == 0 {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.port", path))
	}

	if c.Port < 0 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError(fmt.Sprintf("%s.port", path))
	}

	return nil
}

func getDefaultHttpServerConfig() *HttpServerConfig {
	return &HttpServerConfig{
		Host: httpServerDefaultHost,
		Port: httpServerDefaultPort,
	}
}

func getHttpServerConfigFromEnvironment() (*HttpServerConfig, errors.ApplicationError) {
	cfg := getDefaultHttpServerConfig()

	host := GetEnvironmentString("HTTP_SERVER_HOST", httpServerDefaultHost)
	if host != "" {
		cfg.Host = host
	}

	port, err := GetEnvironmentInt("HTTP_SERVER_PORT", httpServerDefaultPort)
	if err != nil {
		return nil, err
	}

	if port != 0 {
		cfg.Port = port
	}

	return cfg, nil
}
