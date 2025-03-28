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

package services

import (
	"cloud-toolbox/internal/application/faas/models"
	"cloud-toolbox/internal/application/faas/models/interfaces"
	faasmodels "cloud-toolbox/internal/domain/faas/models"
	dbinterfaces "cloud-toolbox/internal/infrastructure/database/repositories/interfaces"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
)

type FunctionRegistry struct {
	executions map[string]*models.ExecutionRecord
	repository dbinterfaces.FunctionExecutionRepository
}

func (f *FunctionRegistry) GetError(executionId string) string {
	fn, err := f.repository.GetFunction(executionId)
	if err != nil {
		return ""
	}

	return fn.Error
}

func (f *FunctionRegistry) GetRecord(executionId string) *models.ExecutionRecord {
	return f.executions[executionId]
}

func (f *FunctionRegistry) AddRecord(function *faasmodels.FunctionExecution, request *models.FaasRequest[interfaces.FaasRecord]) errors.ApplicationError {
	record := models.NewExecutionRecord(
		function.GetId(),
		request,
		function,
	)

	fn := models.FunctionExecutionFromRecord(record)
	err := f.repository.SaveFunction(fn)

	f.executions[function.GetId()] = record
	function.AddListener(f)

	return err
}

func (f *FunctionRegistry) GetStatus(executionId string) string {
	fn, err := f.repository.GetFunction(executionId)
	if err != nil {
		return ""
	}

	return fn.Status
}

func (f *FunctionRegistry) GetOutput(executionId string) string {
	fn, err := f.repository.GetFunction(executionId)
	if err != nil {
		return ""
	}

	return fn.Output
}

func (f *FunctionRegistry) Notify(executionId string, status int) {
	record, _ := f.executions[executionId]
	if record == nil {
		return
	}

	newStatus := record.FunctionExecution.GetStatus()
	f.repository.UpdateStatus(executionId, newStatus)
	if record.FunctionExecution.IsFinished() {
		err := record.FunctionExecution.GetError()
		if err != nil {
			f.repository.SaveError(executionId, err.Error())
		} else {
			output, _ := record.FunctionExecution.GetOutput()
			f.repository.SaveOutput(executionId, output)
		}

		delete(f.executions, executionId)
	}
}

func NewFunctionRegistry(container *di.Container) (*FunctionRegistry, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError("FunctionRegistry")
	}

	if container.FunctionExecutionRepository == nil {
		return nil, errors.NewResolveDependencyError("FunctionRegistry", "FunctionExecutionRepository")
	}

	return &FunctionRegistry{
		executions: make(map[string]*models.ExecutionRecord),
		repository: container.FunctionExecutionRepository,
	}, nil
}
