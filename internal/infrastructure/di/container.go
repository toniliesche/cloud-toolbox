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
	dbinterfaces "cloud-toolbox/internal/infrastructure/database/repositories/interfaces"
	httpinterfaces "cloud-toolbox/internal/infrastructure/http/interfaces"
	rmqinterfaces "cloud-toolbox/internal/infrastructure/rabbitmq/interfaces"
	"context"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

type Container struct {
	Context                     context.Context
	FunctionAsAServiceConfig    *config.FunctionAsAServiceConfig
	FunctionAsAServiceHandler   httpinterfaces.HttpHandler
	FunctionAsAServiceService   faasinterfaces.FunctionAsAService
	FunctionExecutionRepository dbinterfaces.FunctionExecutionRepository
	FunctionRegistry            faasinterfaces.FunctionRegistry
	FunctionTriggerConfig       *config.FunctionTriggerConfig
	FunctionTriggerHandler      rmqinterfaces.RabbitMQMessageHandler
	HttpServer                  httpinterfaces.HttpServer
	HttpServerConfig            *config.HttpServerConfig
	Logger                      *zerolog.Logger
	RabbitMQConnectionConsumer  *amqp.Connection
	RabbitMQConnectionProducer  *amqp.Connection
	RabbitMQConfig              *config.RabbitMQConfig
	RabbitMQConsumer            rmqinterfaces.RabbitMQConsumer
	Redis                       *redis.Client
	RedisConfig                 *config.RedisConfig
	SystemConfig                *config.SystemConfig
	Scylla                      *dynamodb.DynamoDB
	ScyllaConfig                *config.ScyllaConfig
}
