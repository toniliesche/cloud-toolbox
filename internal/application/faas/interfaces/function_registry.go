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

package interfaces

import (
	faasappmodels "cloud-toolbox/internal/application/faas/models"
	faasinterfaces "cloud-toolbox/internal/application/faas/models/interfaces"
	"cloud-toolbox/internal/domain/faas/interfaces"
	faasdomainmodels "cloud-toolbox/internal/domain/faas/models"
)

type FunctionRegistry interface {
	interfaces.UpdateListener
	AddRecord(function *faasdomainmodels.FunctionExecution, request *faasappmodels.FaasRequest[faasinterfaces.FaasRecord]) error
	GetRecord(executionId string) *faasappmodels.ExecutionRecord
	GetStatus(executionId string) string
	GetOutput(id string) string
	GetError(id string) error
}
