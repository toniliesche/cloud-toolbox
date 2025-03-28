package interfaces

import (
	"cloud-toolbox/internal/application/faas/models"
	"cloud-toolbox/internal/infrastructure/errors"
)

type FunctionExecutionRepository interface {
	GetFunction(executionId string) (*models.FunctionExecution, errors.ApplicationError)
	SaveError(executionId string, error string) errors.ApplicationError
	SaveFunction(function *models.FunctionExecution) errors.ApplicationError
	SaveOutput(executionId string, output string) errors.ApplicationError
	UpdateStatus(executionId string, status string) errors.ApplicationError
}
