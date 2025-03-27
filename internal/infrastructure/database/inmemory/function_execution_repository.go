package inmemory

import "cloud-toolbox/internal/infrastructure/di"

type FunctionExecutionRepository struct {
}

func NewFunctionExecutionRepository(container *di.Container) (*FunctionExecutionRepository, errors.ApplicationError) {
	return &FunctionExecutionRepository{}, nil
}
