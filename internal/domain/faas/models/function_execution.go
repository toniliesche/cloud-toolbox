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
	faaserrors "cloud-toolbox/internal/domain/faas/errors"
	"cloud-toolbox/internal/domain/faas/interfaces"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"mvdan.cc/sh/v3/shell"
	"os/exec"
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
	err        error
	cancel     context.CancelFunc
	ctx        context.Context
	output     []byte
	body       string
	exitCode   syscall.WaitStatus
	wait       sync.WaitGroup
	failedInfo string
	logger     *zerolog.Logger
	listeners  []interfaces.UpdateListener
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

func (e *FunctionExecution) Run(sync *chan uint) error {
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

		e.updateStatus(StatusRejected)
		return faaserrors.NewLimitExceededError()
	}
}

func (e *FunctionExecution) Wait() {
	e.logger.Trace().
		Msgf("[%s] Waiting for function execution to finish", FunctionExecutionLogIdentifier)

	e.wait.Wait()

	e.logger.Trace().
		Msgf("[%s] Finished waiting for function execution", FunctionExecutionLogIdentifier)
}

func (e *FunctionExecution) IsFinished() bool {
	return e.status == StatusRejected || e.status == StatusStartingFailed || e.status > StatusRunning
}

func (e *FunctionExecution) Terminate() error {
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

	return e.cmd.Process.Kill()
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

func (e *FunctionExecution) GetOutput() (string, error) {
	e.logger.Trace().
		Msgf("[%s] Trying to get function execution output", FunctionExecutionLogIdentifier)

	if e.status == StatusWaiting || e.status == StatusRunning {
		e.logger.Trace().
			Msgf("[%s] Function execution has not finished yet", FunctionExecutionLogIdentifier)

		return "", faaserrors.NewNotFinishedError()
	}

	e.logger.Trace().
		Msgf("[%s] Returning function execution output", FunctionExecutionLogIdentifier)

	return string(e.output), nil
}

func (e *FunctionExecution) GetError() error {
	return e.err
}

func (e *FunctionExecution) launch(sync *chan uint) error {
	e.logger.Trace().
		Msgf("[%s] Launching function execution", FunctionExecutionLogIdentifier)

	e.wait.Add(1)
	go func() {
		defer e.wait.Done()
		defer e.cancel()
		defer func() { <-*sync }()

		e.logger.Trace().
			Msgf("[%s] Starting function execution in go routine", FunctionExecutionLogIdentifier)

		stdin, err := e.cmd.StdinPipe()
		if err != nil {
			e.logger.Trace().
				Err(err).
				Msgf("[%s] Error creating stdin pipe", FunctionExecutionLogIdentifier)

			e.err = err
			e.updateStatus(StatusStartingFailed)
			return
		}

		_, err = fmt.Fprintln(stdin, e.body)
		if err != nil {
			e.logger.Trace().
				Err(err).
				Msgf("[%s] Error writing to stdin", FunctionExecutionLogIdentifier)

			e.err = err
			e.updateStatus(StatusStartingFailed)
			return
		}

		err = stdin.Close()
		if err != nil {
			e.logger.Trace().
				Err(err).
				Msgf("[%s] Error closing stdin", FunctionExecutionLogIdentifier)

			e.err = err
			e.updateStatus(StatusStartingFailed)
			return
		}

		e.updateStatus(StatusRunning)
		e.output, e.err = e.cmd.CombinedOutput()

		if e.err != nil {
			e.logger.Trace().
				Err(e.err).
				Msgf("[%s] Function execution has failed", FunctionExecutionLogIdentifier)

			if errors.Is(e.ctx.Err(), context.DeadlineExceeded) {
				e.err = e.ctx.Err()
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
	logger *zerolog.Logger,
) (*FunctionExecution, error) {
	args, err := shell.Fields(command, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	executionLogger := logger.With().Str("execution-id", id).Logger()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	return &FunctionExecution{
		id:        id,
		logger:    &executionLogger,
		cmd:       cmd,
		ctx:       ctx,
		cancel:    cancel,
		body:      body,
		listeners: make([]interfaces.UpdateListener, 0),
	}, nil
}
