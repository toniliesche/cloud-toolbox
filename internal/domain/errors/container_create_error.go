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

type ContainerCreateError struct {
	message string
}

func (e ContainerCreateError) Error() string {
	return e.message
}

func NewResolveDependencyError(service string, dependency string) error {
	return ContainerCreateError{
		message: fmt.Sprintf("can't initalize service `%s`: can't resolve dependency `%s`", service, dependency),
	}
}

func NewInvalidConfigError(service string, err error) error {
	return ContainerCreateError{
		message: fmt.Sprintf("can't initalize service `%s`: invalid application config: %v", service, err),
	}
}

func NewApplicationSetupError(application string, err error) error {
	return ContainerCreateError{
		message: fmt.Sprintf("can't setup application `%s`: %v", application, err),
	}
}

func NewContainerMissingError(service string) error {
	return ContainerCreateError{
		message: fmt.Sprintf("can't initalize service `%s`: container needed, got `nil` instead", service),
	}
}
