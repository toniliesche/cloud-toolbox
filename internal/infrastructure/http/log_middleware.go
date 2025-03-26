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
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

const LogMiddlewareLogIdentifier = "LogMiddleware"

type LogMiddleware struct {
	logger *zerolog.Logger
}

func (m *LogMiddleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestID := uuid.New().String()

		m.logger.Info().
			Str("request-id", requestID).
			Msgf("[%s] Received request: %s %s", LogMiddlewareLogIdentifier, request.Method, request.URL.Path)

		sizeAwareWriter := &SizeAwareWriter{
			ResponseWriter: writer,
		}

		start := time.Now()

		ctx := context.WithValue(request.Context(), "request_id", requestID)
		next.ServeHTTP(sizeAwareWriter, request.WithContext(ctx))

		duration := time.Since(start)
		m.logger.Info().
			Str("request-id", requestID).
			Msgf("[%s] Completed request: %s %s, Status: %d, Size: %d bytes, Duration: %s",
				LogMiddlewareLogIdentifier, request.Method, request.URL.Path, sizeAwareWriter.statusCode, sizeAwareWriter.size, duration)
	})
}
