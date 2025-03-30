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

package errors

import "fmt"

func NewResolveDependencyError(service string, dependency string) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("can't initalize service `%s`: can't resolve dependency `%s`", service, dependency),
		code:    ErrorCodeContainerMissingDependency,
	}
}

func NewInvalidConfigError(service string, err error) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("can't initalize service `%s`: invalid application config: %v", service, err),
		code:    ErrorCodeContainerInvalidConfig,
	}
}

func NewApplicationSetupError(application string, err error) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("can't setup application `%s`: %v", application, err),
		code:    ErrorCodeContainerApplicationSetup,
	}
}

func NewContainerMissingError(service string) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("can't initalize service `%s`: container needed, got `nil` instead", service),
		code:    ErrorCodeContainerMissing,
	}
}

func NewContainerConfigMissingError() ApplicationError {
	return InfrastructureError{
		message: "no application config provided",
		code:    ErrorCodeContainerConfigMissing,
	}
}

func NewConstructionFailedError(service string, err error) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("can't initalize service `%s`: construction failed: %v", service, err),
		code:    ErrorCodeContainerConstructionFailed,
	}
}
