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
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/http/interfaces"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"net/http"
)

const ServerLogIdentifier = "Server"

type Server struct {
	cfg    *config.HttpServerConfig
	logger *zerolog.Logger
	router *mux.Router
	server *http.Server
}

func (s *Server) Run() errors.ApplicationError {
	s.logger.Info().
		Msgf("[%s] Starting http server", ServerLogIdentifier)

	s.logger.Debug().
		Msgf("[%s] Listening on %s:%d", ServerLogIdentifier, s.cfg.Host, s.cfg.Port)

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		Handler: s.router,
	}

	err := s.server.ListenAndServe()
	if err != nil {
		return errors.NewGenericError(err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) errors.ApplicationError {
	s.logger.Info().
		Msgf("[%s] Shutting down http server", ServerLogIdentifier)

	err := s.server.Shutdown(ctx)
	if err != nil {
		return errors.NewGenericError(err)
	}

	return nil
}

func (s *Server) RegisterRoutes(handler interfaces.HttpHandler) errors.ApplicationError {
	if handler == nil {
		return errors.NewRouteRegisterError(fmt.Errorf("passed HttpHandler is `nil`"))
	}

	routes := handler.GetRoutes()
	if routes == nil {
		return errors.NewMissingRoutesError()
	}

	if len(routes) == 0 {
		return errors.NewMissingRoutesError()
	}

	for _, route := range routes {
		if len(route.Methods) == 0 {
			s.logger.Debug().
				Msgf("[%s] Route %s registered", ServerLogIdentifier, route.Path)
		} else {
			for _, method := range route.Methods {
				s.logger.Debug().
					Msgf("[%s] Route %s %s registered", ServerLogIdentifier, method, route.Path)
			}
		}

		r := s.router.HandleFunc(route.Path, route.Handler).Methods(route.Methods...)

		if err := r.GetError(); err != nil {
			return errors.NewRouteRegisterError(err)
		}
	}

	return nil
}

func NewServer(container *di.Container) (*Server, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError("HttpServer")
	}

	if container.HttpServerConfig == nil {
		return nil, errors.NewResolveDependencyError("HttpServer", "HttpServerConfig")
	}

	if err := container.HttpServerConfig.Validate("http"); err != nil {
		return nil, errors.NewInvalidConfigError("HttpServerConfig", err)
	}

	if container.Logger == nil {
		return nil, errors.NewResolveDependencyError("HttpServer", "Logger")
	}

	router := mux.NewRouter()

	logMiddleware := &LogMiddleware{
		logger: container.Logger,
	}

	router.Use(logMiddleware.LogRequest)

	return &Server{
		cfg:    container.HttpServerConfig,
		logger: container.Logger,
		router: router,
	}, nil
}
