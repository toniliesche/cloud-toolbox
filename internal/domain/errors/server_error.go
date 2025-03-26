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

package errors

import "fmt"

type ServerError struct {
	message string
}

func (e ServerError) Error() string {
	return e.message
}

func NewServerError(message string) error {
	return ServerError{
		message: fmt.Sprintf("Error while registering routes: %s", message),
	}
}

func NewMissingRoutesError() error {
	return NewServerError("HttpHandler did not provide any routes")
}

func NewRouteRegisterError(err error) error {
	return NewServerError(err.Error())
}
