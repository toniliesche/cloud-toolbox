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

package http

import (
	"cloud-toolbox/internal/infrastructure/http/interfaces"
	"net/http"
)

type AuthMiddleware struct {
	JsonResponseHandler
	authenticator interfaces.RequestAuthenticator
}

func (m *AuthMiddleware) AuthenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := m.authenticator.AuthenticateRequest(request)
		if err != nil {
			m.SendErrorResponse(writer, request, http.StatusUnauthorized, err.Error())
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func NewAuthMiddleware(authenticator interfaces.RequestAuthenticator) *AuthMiddleware {
	return &AuthMiddleware{
		authenticator: authenticator,
	}
}
