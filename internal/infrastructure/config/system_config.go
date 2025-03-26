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
	"cloud-toolbox/internal/domain/errors"
	"fmt"
)

type SystemConfig struct {
	Log           *LogConfig `yaml:"log"`
	ComponentType string
	ComponentId   string
}

func (c *SystemConfig) Validate(path string) error {
	if c.Log == nil {
		return errors.NewMissingConfigSectionError(fmt.Sprintf("%s.log", path))
	}

	if err := c.Log.Validate(fmt.Sprintf("%s.log", path)); err != nil {
		return err
	}

	return nil
}

func getSystemConfigFromEnvironment() (*SystemConfig, error) {
	logConfig, err := getLogConfigFromEnvironment()
	if err != nil {
		return nil, err
	}

	return &SystemConfig{
		Log:         logConfig,
		ComponentId: GetEnvironmentString("CLOUD_TOOLBOX_COMPONENT_ID", "unknown"),
	}, nil
}
