package interfaces

import "cloud-toolbox/internal/infrastructure/errors"

type FunctionExecutionRepository interface {
	GetFunction(executionId string) errors.ApplicationError
	SaveError(executionId string, error string) errors.ApplicationError
	SaveFunction() errors.ApplicationError
	SaveOutput(executionId string, output string) errors.ApplicationError
	UpdateStatus(executionId string, status int) errors.ApplicationError
}
