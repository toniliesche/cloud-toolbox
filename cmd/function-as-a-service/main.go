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

package main

import (
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/setup"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.ProvideFunctionAsAServiceConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	container, err := setup.NewBuilder().
		SetFaasConfig(cfg).
		SetContext(ctx).
		Build()

	if err != nil {
		log.Fatal(err)
	}

	logger := container.Logger
	server := container.HttpServer

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = server.Run(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logger.Fatal().Err(err).
				Msg("failed to run server")
		}
	}()

	<-stop
	cancel()
	logger.Info().Msg("Received shutdown signal, initiating shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal().Err(err).
			Err(err).
			Msg("failed to shutdown server")
	}

	logger.Info().
		Msg("Server gracefully stopped")
}
