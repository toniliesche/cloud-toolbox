package errors

type ApplicationError interface {
	Code() int
	error
}
