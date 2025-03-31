package errors

import "fmt"

func NewEndpointNotImplementedError(route string) ApplicationError {
	return InfrastructureError{
		message: fmt.Sprintf("Endpoint '%s' is not implemented", route),
		code:    ErrorCodeEndpointNotImplemented,
	}
}
