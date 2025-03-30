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

package services

import (
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/rabbitmq/interfaces"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"slices"
	"sync"
	"time"
)

const (
	BatchProcessorLogIdentifier = "BatchProcessor"
)

type BatchProcessor struct {
	mu       sync.Mutex
	buffer   []amqp.Delivery
	limit    int
	timer    *time.Timer
	timeout  time.Duration
	handler  interfaces.RabbitMQMessageHandler
	logger   *zerolog.Logger
	batchErr errors.ApplicationError
}

func (p *BatchProcessor) AddMessage(msg amqp.Delivery) errors.ApplicationError {
	if p.batchErr != nil {
		p.logger.Error().
			Err(p.batchErr).
			Msgf("[%s/%s] Batch processor is in error state, rejecting message (AddMessage)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())

		return p.batchErr
	}

	p.logger.Trace().
		Msgf("[%s/%s] Adding message to batch, current size: %d (AddMessage)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier(), len(p.buffer))

	p.logger.Trace().
		Msgf("[%s/%s] Locking batch processor (AddMessage)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
	p.mu.Lock()
	p.logger.Trace().
		Msgf("[%s/%s] Locked batch processor (AddMessage)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())

	p.buffer = append(p.buffer, msg)
	if len(p.buffer) >= p.limit {
		p.logger.Trace().
			Msgf("[%s/%s] Batch limit reached, processing batch of %d messages (AddMessage)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier(), len(p.buffer))

		p.logger.Trace().
			Msgf("[%s/%s] Unlocking batch processor (AddMessage)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
		p.mu.Unlock()
		err := p.processBatch()
		if err != nil {
			p.logger.Error().
				Err(err).
				Msgf("[%s/%s] Error processing batch (AddMessage)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
			return errors.NewGenericError(err)
		}

		return err
	}

	if len(p.buffer) == 1 {
		p.logger.Trace().
			Msgf("[%s/%s] First message added to batch, starting timer", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
		p.timer.Reset(p.timeout)
	}

	p.logger.Trace().
		Msgf("[%s/%s] Unlocking batch processor (AddMessage)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
	p.mu.Unlock()

	return nil
}

func (p *BatchProcessor) processBatch() errors.ApplicationError {
	p.logger.Trace().
		Msgf("[%s/%s] Locking batch processor (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
	p.mu.Lock()
	p.logger.Trace().
		Msgf("[%s/%s] Locked batch processor (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
	defer func() {
		p.logger.Trace().
			Msgf("[%s/%s] Unlocking batch processor (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
		p.mu.Unlock()
	}()

	p.logger.Trace().
		Msgf("[%s/%s] Processing batch of %d messages (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier(), len(p.buffer))

	if len(p.buffer) == 0 {
		p.logger.Trace().
			Msgf("[%s/%s] No messages to process (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())

		return nil
	}

	failed, err := p.handler.HandleMessageBatch(p.buffer)
	if err != nil {
		p.logger.Error().
			Err(err).
			Msgf("[%s/%s] Error processing batch (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())

		return errors.NewGenericError(err)
	}

	p.logger.Trace().
		Msgf("[%s/%s] Batch processed, %d messages failed (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier(), len(failed))
	for _, msg := range p.buffer {
		p.logger.Trace().
			Msgf("[%s/%s] Processing message %s (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier(), msg.MessageId)
		if slices.Contains(failed, msg.MessageId) {
			p.logger.Trace().
				Msgf("[%s/%s] Message %s failed, rejecting (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier(), msg.MessageId)
			if err := msg.Nack(false, false); err != nil {
				return errors.NewGenericError(err)
			}
		} else {
			p.logger.Trace().
				Msgf("[%s/%s] Message %s succeeded, acknowledging (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier(), msg.MessageId)
			if err := msg.Ack(false); err != nil {
				return errors.NewGenericError(err)
			}
		}
	}

	p.logger.Trace().
		Msgf("[%s/%s] Batch processed, %d messages acknowledged (processBatch)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier(), len(p.buffer)-len(failed))

	p.buffer = make([]amqp.Delivery, 0, p.limit)

	return nil
}
func (p *BatchProcessor) startTimeoutListener() {
	for {
		<-p.timer.C
		p.logger.Trace().
			Msgf("[%s/%s] Locking batch processor (startTimeoutListener)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
		p.mu.Lock()
		p.logger.Trace().
			Msgf("[%s/%s] Locked batch processor (startTimeoutListener)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
		if len(p.buffer) > 0 {
			p.logger.Trace().
				Msgf("[%s/%s] Timeout reached, processing batch of %d messages (startTimeoutListener)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier(), len(p.buffer))
			p.logger.Trace().
				Msgf("[%s/%s] Unlocking batch processor (startTimeoutListener)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
			p.mu.Unlock()
			if p.batchErr = p.processBatch(); p.batchErr != nil {
				p.logger.Error().
					Err(p.batchErr).
					Msgf("[%s/%s] Error processing batch (startTimeoutListener)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
				break
			}
		} else {
			p.logger.Trace().
				Msgf("[%s/%s] Timeout reached, but no messages to process (startTimeoutListener)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
			p.logger.Trace().
				Msgf("[%s/%s] Unlocking batch processor (startTimeoutListener)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
			p.mu.Unlock()
		}

		p.logger.Trace().
			Msgf("[%s/%s] Resetting timer (startTimeoutListener)", BatchProcessorLogIdentifier, p.handler.QueueIdentifier())
		p.timer.Reset(p.timeout)
	}
}

func NewBatchProcessor(timeout int64, limit int64, handler interfaces.RabbitMQMessageHandler, logger *zerolog.Logger) *BatchProcessor {
	timeoutObj := time.Duration(timeout) * time.Second
	processor := &BatchProcessor{
		buffer:  make([]amqp.Delivery, 0, limit),
		limit:   int(limit),
		timer:   time.NewTimer(timeoutObj),
		timeout: timeoutObj,
		handler: handler,
		logger:  logger,
	}

	go processor.startTimeoutListener()

	return processor
}
