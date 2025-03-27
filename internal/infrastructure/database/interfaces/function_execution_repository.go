package interfaces

type FunctionExecutionRepository interface {
	GetFunction(executionId string) error
	SaveError(executionId string, error string) error
	SaveFunction() error
	SaveOutput(executionId string, output string) error
	UpdateStatus(executionId string, status int) error
}
