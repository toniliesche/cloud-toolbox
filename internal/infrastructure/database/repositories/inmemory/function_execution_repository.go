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

package inmemory

import (
	"cloud-toolbox/internal/application/faas/models"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
)

type FunctionExecutionRepository struct {
	storage map[string]*models.FunctionExecution
}

func (f *FunctionExecutionRepository) GetFunction(executionId string) (*models.FunctionExecution, errors.ApplicationError) {
	if function, ok := f.storage[executionId]; ok {
		return function, nil
	}

	return nil, errors.NewItemNotFoundError(executionId)
}

func (f *FunctionExecutionRepository) SaveError(executionId string, error string) errors.ApplicationError {
	if _, ok := f.storage[executionId]; !ok {
		return errors.NewItemNotFoundError(executionId)
	}

	f.storage[executionId].Error = error

	return nil
}

func (f *FunctionExecutionRepository) SaveFunction(function *models.FunctionExecution) errors.ApplicationError {
	f.storage[function.Id] = function

	return nil
}

func (f *FunctionExecutionRepository) SaveOutput(executionId string, output string) errors.ApplicationError {
	if _, ok := f.storage[executionId]; !ok {
		return errors.NewItemNotFoundError(executionId)
	}

	f.storage[executionId].Output = output

	return nil
}

func (f *FunctionExecutionRepository) UpdateStatus(executionId string, status string) errors.ApplicationError {
	if _, ok := f.storage[executionId]; !ok {
		return errors.NewItemNotFoundError(executionId)
	}

	f.storage[executionId].Status = status

	return nil
}

func NewFunctionExecutionRepository(container *di.Container) (*FunctionExecutionRepository, errors.ApplicationError) {
	return &FunctionExecutionRepository{
		storage: make(map[string]*models.FunctionExecution),
	}, nil
}
