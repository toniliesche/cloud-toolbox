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
)

func (b *ContainerBuilder) setupFt(container *di.Container) errors.ApplicationError {
	var err errors.ApplicationError

	if b.ftConfig == nil {
		return errors.NewResolveDependencyError("Function trigger", "ftConfig")
	}

	if err = b.ftConfig.Validate(); err != nil {
		return errors.NewApplicationSetupError("Function trigger", err)
	}
	container.FunctionTriggerConfig = b.ftConfig

	return nil
}
