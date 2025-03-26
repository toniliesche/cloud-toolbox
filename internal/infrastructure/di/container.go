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

package di

import (
	faasinterfaces "cloud-toolbox/internal/application/faas/interfaces"
	"cloud-toolbox/internal/infrastructure/config"
	httpinterfaces "cloud-toolbox/internal/infrastructure/http/interfaces"
	"github.com/rs/zerolog"
)

type Container struct {
	FunctionAsAServiceService faasinterfaces.FunctionAsAService
	FunctionAsAServiceHandler httpinterfaces.HttpHandler
	HttpServerConfig          *config.HttpServerConfig
	HttpServer                httpinterfaces.HttpServer
	SystemConfig              *config.SystemConfig
	Logger                    *zerolog.Logger
	FunctionAsAServiceConfig  *config.FunctionAsAServiceConfig
	FunctionRegistry          faasinterfaces.FunctionRegistry
}
