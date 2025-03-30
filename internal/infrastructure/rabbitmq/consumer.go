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

package rabbitmq

import (
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/rabbitmq/interfaces"
	"cloud-toolbox/internal/infrastructure/rabbitmq/services"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"sync"
)

const (
	ConsumerLogIdentifier = "RabbitMQConsumer"
)

type Consumer struct {
	cfg          *config.RabbitMQConfig
	rabbitMQ     *amqp.Connection
	logger       *zerolog.Logger
	wait         sync.WaitGroup
	subscribers  []*services.Subscriber
	BatchSize    int64
	BatchTimeout int64
}

func (c *Consumer) WithBatchSize(size int64) interfaces.RabbitMQConsumer {
	c.BatchSize = size

	return c
}

func (c *Consumer) WithBatchTimeout(timeout int64) interfaces.RabbitMQConsumer {
	c.BatchTimeout = timeout

	return c
}

func (c *Consumer) Run() errors.ApplicationError {
	defer c.rabbitMQ.Close()
	ch, err := c.rabbitMQ.Channel()
	if err != nil {
		c.logger.Error().
			Err(err).
			Msgf("[%s] Error creating channel", ConsumerLogIdentifier)
		return errors.NewGenericError(err)
	}
	defer ch.Close()

	for _, subscriber := range c.subscribers {
		msgs, err := ch.Consume(
			subscriber.Queue.Name,
			"",
			subscriber.Queue.AutoAck,
			subscriber.Queue.Exclusive,
			subscriber.Queue.NoLocal,
			subscriber.Queue.NoWait,
			subscriber.Queue.Args,
		)

		if err != nil {
			c.logger.Error().
				Err(err).
				Msgf("[%s] Error consuming messages from queue: %s", ConsumerLogIdentifier, subscriber.Queue.Name)
			c.wait.Done()
			return errors.NewGenericError(err)
		}

		c.wait.Add(1)
		go func(subscriber *services.Subscriber) {
			processor := services.NewBatchProcessor(
				c.BatchTimeout,
				c.BatchSize,
				subscriber.Handler,
				c.logger,
			)

			for msg := range msgs {
				msgContent, _ := json.Marshal(msg)

				c.logger.Trace().
					Str("msg", string(msgContent)).
					Msgf("[%s] Received message from queue: %s", ConsumerLogIdentifier, subscriber.Queue.Name)

				if err := processor.AddMessage(msg); err != nil {
					break
				}
			}

			c.wait.Done()
		}(subscriber)
	}

	c.wait.Wait()

	return nil
}

func (c *Consumer) RegisterHandler(handler interfaces.RabbitMQMessageHandler) errors.ApplicationError {
	c.logger.Trace().
		Msgf("[%s] Registering new handler", ConsumerLogIdentifier)
	if handler == nil {
		c.logger.Trace().
			Msgf("[%s] RabbitMQMessageHandler is nil", ConsumerLogIdentifier)
		return errors.NewGenericError(fmt.Errorf("handler is nil"))
	}

	id := handler.QueueIdentifier()
	if id == "" {
		c.logger.Trace().
			Msgf("[%s] RabbitMQMessageHandler has no identifier", ConsumerLogIdentifier)
		return errors.NewGenericError(fmt.Errorf("handler has no identifier"))
	}

	c.logger.Trace().
		Msgf("[%s] Retrieving config for queue: %s", ConsumerLogIdentifier, id)
	queue, err := c.cfg.GetQueue(id)
	if err != nil {
		c.logger.Trace().
			Msgf("[%s] Error retrieving config for queue: %s", ConsumerLogIdentifier, id)
		return errors.NewGenericError(err)
	}

	subscriber := &services.Subscriber{
		Queue:   queue,
		Handler: handler,
	}

	c.subscribers = append(c.subscribers, subscriber)

	return nil
}

func NewConsumer(container *di.Container) (*Consumer, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError(ConsumerLogIdentifier)
	}

	if container.RabbitMQConfig == nil {
		return nil, errors.NewResolveDependencyError(ConsumerLogIdentifier, "RabbitMQConfig")
	}

	if err := container.RabbitMQConfig.Validate("", config.RabbitMQModeConsumer); err != nil {
		return nil, errors.NewInvalidConfigError(ConsumerLogIdentifier, err)
	}

	if container.Logger == nil {
		return nil, errors.NewResolveDependencyError(ConsumerLogIdentifier, "Logger")
	}

	if container.RabbitMQConnectionConsumer == nil {
		return nil, errors.NewResolveDependencyError(ConsumerLogIdentifier, "RabbitMQConnectionConsumer")
	}

	return &Consumer{
		cfg:         container.RabbitMQConfig,
		rabbitMQ:    container.RabbitMQConnectionConsumer,
		logger:      container.Logger,
		subscribers: make([]*services.Subscriber, 0),
	}, nil
}
