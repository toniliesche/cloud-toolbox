package redis

import "cloud-toolbox/internal/infrastructure/di"

type FunctionExecutionRepository struct {
}

func NewFunctionExecutionRepository(container *di.Container) (*FunctionExecutionRepository, error) {
	return &FunctionExecutionRepository{}, nil
}
