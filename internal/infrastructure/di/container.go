package di

import (
	epinterfaces "cloud-toolbox/internal/application/ep/interfaces"
	faasinterfaces "cloud-toolbox/internal/application/faas/interfaces"
	"cloud-toolbox/internal/infrastructure/config"
	dbinterfaces "cloud-toolbox/internal/infrastructure/database/repositories/interfaces"
	httpinterfaces "cloud-toolbox/internal/infrastructure/http/interfaces"
	rmqinterfaces "cloud-toolbox/internal/infrastructure/rabbitmq/interfaces"
	"context"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Container struct {
	Context                        context.Context
	EventPublisher                 epinterfaces.EventPublisher
	EventPublisherConfig           *config.EventPublisherConfig
	FunctionAsAServiceConfig       *config.FunctionAsAServiceConfig
	EventPublisherHttpHandler      httpinterfaces.HttpHandler
	FunctionAsAServiceHttpHandler  httpinterfaces.HttpHandler
	FunctionAsAService             faasinterfaces.FunctionAsAService
	FunctionExecutionRepository    dbinterfaces.FunctionExecutionRepository
	FunctionRegistry               faasinterfaces.FunctionRegistry
	FunctionTriggerConfig          *config.FunctionTriggerConfig
	FunctionTriggerRabbitMQHandler rmqinterfaces.RabbitMQMessageHandler
	HttpServer                     httpinterfaces.HttpServer
	HttpServerConfig               *config.HttpServerConfig
	Logger                         *zerolog.Logger
	RabbitMQConnectionConsumer     *amqp091.Connection
	RabbitMQConnectionPublisher    *amqp091.Connection
	RabbitMQConfig                 *config.RabbitMQConfig
	RabbitMQConsumer               rmqinterfaces.RabbitMQConsumer
	Redis                          *redis.Client
	RedisConfig                    *config.RedisConfig
	SystemConfig                   *config.SystemConfig
	Scylla                         *dynamodb.DynamoDB
	ScyllaConfig                   *config.ScyllaConfig
}
