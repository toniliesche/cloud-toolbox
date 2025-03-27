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
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
)

type FunctionRegistry struct {
	executions map[string]*models.ExecutionRecord
	statuses   map[string]string
	outputs    map[string]string
	errors     map[string]errors.ApplicationError
}

func (f *FunctionRegistry) AddRecord(function *faasmodels.FunctionExecution, request *models.FaasRequest[interfaces.FaasRecord]) errors.ApplicationError {
	record := models.NewExecutionRecord(
		function.GetId(),
		request,
		function,
	)

	f.executions[function.GetId()] = record
	function.AddListener(f)

	return nil
}

func (f *FunctionRegistry) GetRecord(executionId string) *models.ExecutionRecord {
	record, ok := f.executions[executionId]
	if !ok {
		return nil
	}

	return record
}

func (f *FunctionRegistry) GetStatus(executionId string) string {
	return f.statuses[executionId]
}

func (f *FunctionRegistry) GetOutput(executionId string) string {
	return f.outputs[executionId]
}

func (f *FunctionRegistry) GetError(executionId string) errors.ApplicationError {
	return f.errors[executionId]
}

func (f *FunctionRegistry) Notify(executionId string, status int) {
	record, _ := f.executions[executionId]
	if record == nil {
		return
	}

	f.statuses[executionId] = record.FunctionExecution.GetStatus()
	if record.FunctionExecution.IsFinished() {
		f.outputs[executionId], _ = record.FunctionExecution.GetOutput()
		f.errors[executionId] = record.FunctionExecution.GetError()
		delete(f.executions, executionId)
	}
}

func NewFunctionRegistry(container *di.Container) (*FunctionRegistry, errors.ApplicationError) {
	return &FunctionRegistry{
		executions: make(map[string]*models.ExecutionRecord),
		statuses:   make(map[string]string),
		outputs:    make(map[string]string),
		errors:     make(map[string]errors.ApplicationError),
	}, nil
}
