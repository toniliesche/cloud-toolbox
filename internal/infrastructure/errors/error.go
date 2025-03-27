package errors

type InfrastructureError struct {
	message string
	code    int
}

func (e InfrastructureError) Error() string {
	return e.message
}

func (e InfrastructureError) Code() int {
	return e.code
}

func NewGenericError(err error) ApplicationError {
	return InfrastructureError{
		message: err.Error(),
		code:    ErrorCodeGenericError,
	}
}
