package errors

import "fmt"

func NewItemNotFoundError(item string) ApplicationError {
	return &InfrastructureError{
		message: fmt.Sprintf("Item `%s` not found", item),
		code:    ErrorCodeDatabaseItemNotFound,
	}
}
