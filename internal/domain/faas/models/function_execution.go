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

package models

import (
	"bufio"
	"bytes"
	faaserrors "cloud-toolbox/internal/application/faas/errors"
	"cloud-toolbox/internal/domain/faas/interfaces"
	infrastructureerrors "cloud-toolbox/internal/infrastructure/errors"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"mvdan.cc/sh/v3/shell"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	FunctionExecutionLogIdentifier = "FunctionExecution"
	StatusWaiting                  = 0
	StatusStartingFailed           = 1
	StatusRejected                 = 2
	StatusRunning                  = 3
	StatusFinished                 = 4
	StatusError                    = 5
	StatusTimeout                  = 6
)

type FunctionExecution struct {
	id         string
	cmd        *exec.Cmd
	status     int
	err        infrastructureerrors.ApplicationError
	cancel     context.CancelFunc
	ctx        context.Context
	out        *bytes.Buffer
	logs       *bytes.Buffer
	body       string
	exitCode   syscall.WaitStatus
	wait       sync.WaitGroup
	failedInfo string
	logger     *zerolog.Logger
	listeners  []interfaces.UpdateListener
	logReader  *os.File
	logWriter  *os.File
}

func (e *FunctionExecution) updateStatus(status int) {
	e.status = status

	for _, listener := range e.listeners {
		listener.Notify(e.id, status)
	}
}

func (e *FunctionExecution) GetId() string {
	return e.id
}

func (e *FunctionExecution) AddListener(listener interfaces.UpdateListener) {
	e.listeners = append(e.listeners, listener)
}

func (e *FunctionExecution) Run(sync *chan uint) infrastructureerrors.ApplicationError {
	e.logger.Trace().
		Msgf("[%s] Trying to execute function: %s", FunctionExecutionLogIdentifier, e.cmd)

	if e.status > 0 {
		e.logger.Trace().
			Msgf("[%s] Function execution has already been run", FunctionExecutionLogIdentifier)

		return faaserrors.NewAlreadyStartedError()
	}

	select {
	case *sync <- 0:
		e.logger.Trace().
			Msgf("[%s] Launching function execution", FunctionExecutionLogIdentifier)

		return e.launch(sync)
	default:
		e.logger.Trace().
			Msgf("[%s] Function execution has been rejected", FunctionExecutionLogIdentifier)

		e.err = faaserrors.NewLimitExceededError()
		e.updateStatus(StatusRejected)

		return e.err
	}
}

func (e *FunctionExecution) Wait() infrastructureerrors.ApplicationError {
	e.logger.Trace().
		Msgf("[%s] Waiting for function execution to finish", FunctionExecutionLogIdentifier)

	e.wait.Wait()

	e.logger.Trace().
		Msgf("[%s] Finished waiting for function execution", FunctionExecutionLogIdentifier)

	return e.err
}

func (e *FunctionExecution) IsFinished() bool {
	return e.status == StatusRejected || e.status == StatusStartingFailed || e.status > StatusRunning
}

func (e *FunctionExecution) Terminate() infrastructureerrors.ApplicationError {
	e.logger.Trace().
		Msgf("[%s] Trying to terminate function execution", FunctionExecutionLogIdentifier)

	if e.status == StatusWaiting {
		e.logger.Trace().
			Msgf("[%s] Function execution has not been started yet", FunctionExecutionLogIdentifier)

		return faaserrors.NewNotStartedError()
	}

	if e.IsFinished() {
		e.logger.Trace().
			Msgf("[%s] Function execution has already been finished", FunctionExecutionLogIdentifier)

		return faaserrors.NewAlreadyFinishedError()
	}

	e.logger.Trace().
		Msgf("[%s] Terminating function execution", FunctionExecutionLogIdentifier)

	err := e.cmd.Process.Kill()
	if err != nil {
		return faaserrors.NewFaasCommandCouldNotBeKilledError(e.cmd.Process.Kill())
	}

	return nil
}

func (e *FunctionExecution) GetStatus() string {
	e.logger.Trace().
		Msgf("[%s] Getting function execution status", FunctionExecutionLogIdentifier)

	switch e.status {
	case StatusWaiting:
		e.logger.Trace().
			Msgf("[%s] Function execution is waiting", FunctionExecutionLogIdentifier)

		return "waiting"
	case StatusStartingFailed:
		e.logger.Trace().
			Msgf("[%s] Function execution has failed to start", FunctionExecutionLogIdentifier)

		return "failed-start"
	case StatusRejected:
		e.logger.Trace().
			Msgf("[%s] Function execution has been rejected", FunctionExecutionLogIdentifier)

		return "rejected"
	case StatusRunning:
		e.logger.Trace().
			Msgf("[%s] Function execution is running", FunctionExecutionLogIdentifier)

		return "running"
	case StatusTimeout:
		e.logger.Trace().
			Msgf("[%s] Function execution has timed out", FunctionExecutionLogIdentifier)

		return "timeout"
	case StatusFinished:
		e.logger.Trace().
			Msgf("[%s] Function execution has finished", FunctionExecutionLogIdentifier)

		return "success"
	default:
		e.logger.Trace().
			Msgf("[%s] Function execution has failed", FunctionExecutionLogIdentifier)

		return "failed-run"
	}
}

func (e *FunctionExecution) GetOutput() (string, infrastructureerrors.ApplicationError) {
	e.logger.Trace().
		Msgf("[%s] Trying to get function execution output", FunctionExecutionLogIdentifier)

	if e.status == StatusWaiting || e.status == StatusRunning {
		e.logger.Trace().
			Msgf("[%s] Function execution has not finished yet", FunctionExecutionLogIdentifier)

		return "", faaserrors.NewNotFinishedError()
	}

	e.logger.Trace().
		Msgf("[%s] Returning function execution output", FunctionExecutionLogIdentifier)

	return e.out.String(), nil
}

func (e *FunctionExecution) GetError() infrastructureerrors.ApplicationError {
	return e.err
}

func (e *FunctionExecution) launch(sync *chan uint) infrastructureerrors.ApplicationError {
	e.logger.Trace().
		Msgf("[%s] Launching function execution", FunctionExecutionLogIdentifier)

	e.wait.Add(1)
	go func() {
		defer e.wait.Done()
		defer e.cancel()
		defer e.logWriter.Close()
		defer func() { <-*sync }()

		e.logger.Trace().
			Msgf("[%s] Starting function execution in go routine", FunctionExecutionLogIdentifier)

		e.cmd.Stdout = e.out
		e.cmd.Stderr = e.logWriter

		go func() {
			scanner := bufio.NewScanner(e.logReader)
			for scanner.Scan() {
				text := strings.TrimSpace(scanner.Text())

				if text == "" {
					continue
				}

				e.logger.Trace().
					Msgf("[%s] Function execution log: %s", FunctionExecutionLogIdentifier, scanner.Text())
			}
		}()

		stdin, err := e.cmd.StdinPipe()
		if err != nil {
			e.logger.Trace().
				Err(err).
				Msgf("[%s] Error creating stdin pipe", FunctionExecutionLogIdentifier)

			e.err = faaserrors.NewExecutionFailedError(err)
			e.updateStatus(StatusStartingFailed)
			return
		}

		_, err = fmt.Fprintln(stdin, e.body)
		if err != nil {
			e.logger.Trace().
				Err(err).
				Msgf("[%s] Error writing to stdin", FunctionExecutionLogIdentifier)

			e.err = faaserrors.NewExecutionFailedError(err)
			e.updateStatus(StatusStartingFailed)
			return
		}

		err = stdin.Close()
		if err != nil {
			e.logger.Trace().
				Err(err).
				Msgf("[%s] Error closing stdin", FunctionExecutionLogIdentifier)

			e.err = faaserrors.NewExecutionFailedError(err)
			e.updateStatus(StatusStartingFailed)
			return
		}

		e.updateStatus(StatusRunning)

		err = e.cmd.Run()
		if err != nil {
			e.err = faaserrors.NewExecutionFailedError(err)
		}

		if e.err != nil {
			e.logger.Trace().
				Err(e.err).
				Msgf("[%s] Function execution has failed", FunctionExecutionLogIdentifier)

			if errors.Is(e.ctx.Err(), context.DeadlineExceeded) {
				e.err = faaserrors.NewExecutionTimeoutError(e.ctx.Err())
				e.updateStatus(StatusTimeout)
			} else {
				e.updateStatus(StatusError)
			}

			var exitError *exec.ExitError
			if errors.As(e.err, &exitError) {
				e.exitCode = exitError.Sys().(syscall.WaitStatus)

				e.logger.Trace().
					Msgf("[%s] Function execution exit code: %d", FunctionExecutionLogIdentifier, e.exitCode.ExitStatus())
			} else {
				e.logger.Trace().
					Msgf("[%s] Function execution error type: %T", FunctionExecutionLogIdentifier, e.err)
			}
		} else {
			e.logger.Trace().
				Msgf("[%s] Function execution has finished successfully", FunctionExecutionLogIdentifier)

			e.logger.Trace().
				Msgf("[%s] Function execution exit code: 0", FunctionExecutionLogIdentifier)

			e.updateStatus(StatusFinished)
		}

		e.logger.Trace().
			Msgf("[%s] Function execution go routine has finished", FunctionExecutionLogIdentifier)
	}()

	return nil
}

func NewFunctionExecution(
	id string,
	command string,
	timeout int64,
	body string,
	ctx context.Context,
	logger *zerolog.Logger,
) (*FunctionExecution, infrastructureerrors.ApplicationError) {
	args, err := shell.Fields(command, nil)
	if err != nil {
		return nil, faaserrors.NewFaasCommandCouldNotBeParsedError(err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	executionLogger := logger.With().Str("execution-id", id).Logger()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	logReader, logWriter, _ := os.Pipe()

	return &FunctionExecution{
		id:        id,
		logger:    &executionLogger,
		cmd:       cmd,
		out:       &bytes.Buffer{},
		logReader: logReader,
		logWriter: logWriter,
		ctx:       ctx,
		cancel:    cancel,
		body:      body,
		listeners: make([]interfaces.UpdateListener, 0),
	}, nil
}
