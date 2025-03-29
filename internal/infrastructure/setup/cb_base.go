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

package setup

import (
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/log"
)

func (b *ContainerBuilder) setupBasics(container *di.Container) errors.ApplicationError {

	return nil
}

func (b *ContainerBuilder) setupLogger(container *di.Container) errors.ApplicationError {
	var err errors.ApplicationError

	applicationConfig, err := b.getApplicationConfig()
	if err != nil {
		return errors.NewResolveDependencyError("Generic", "applicationConfig")
	}

	container.SystemConfig = applicationConfig.GetSystemConfig()
	container.SystemConfig.ComponentType = b.getComponentType()

	if err = container.SystemConfig.Validate("system"); err != nil {
		return errors.NewApplicationSetupError("Generic", err)
	}

	container.Logger, err = log.NewLogger(container)
	if err != nil {
		return errors.NewApplicationSetupError("Generic", err)
	}

	return nil
}
