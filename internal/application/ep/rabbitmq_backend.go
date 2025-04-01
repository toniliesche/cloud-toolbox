package ep

import (
	"cloud-toolbox/internal/application/ep/models"
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/rabbitmq/interfaces"
	"fmt"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"time"
)

const RabbitMQBackendLogIdentifier = "RabbitMQBackend"

type RabbitMQBackend struct {
	cfg       *config.EventPublisherConfig
	publisher interfaces.RabbitMQPublisher
	logger    *zerolog.Logger
}

func (r *RabbitMQBackend) PublishEvent(event *models.EpRequest) errors.ApplicationError {
	timestamp := event.Timestamp
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	msg := amqp091.Publishing{
		AppId:         fmt.Sprintf("ctb:ep:%s", r.cfg.PublisherName),
		Body:          []byte(event.Event),
		ContentType:   "application/json",
		DeliveryMode:  2,
		CorrelationId: event.CorrelationId,
		MessageId:     uuid.New().String(),
		Timestamp:     timestamp,
	}

	exchange, err := r.cfg.RabbitMQ.Publisher.GetExchange("default")
	if err != nil {
		return errors.NewGenericError(err)
	}

	return r.publisher.Publish(msg, exchange.Name, r.cfg.RabbitMQTopic)
}

func NewRabbitMQBackend(container *di.Container) (*RabbitMQBackend, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError(RabbitMQBackendLogIdentifier)
	}

	if container.EventPublisherConfig == nil {
		return nil, errors.NewResolveDependencyError(RabbitMQBackendLogIdentifier, "EventPublisherConfig")
	}

	if err := container.EventPublisherConfig.Validate(); err != nil {
		return nil, errors.NewApplicationSetupError(RabbitMQBackendLogIdentifier, err)
	}

	if container.RabbitMQPublisher == nil {
		return nil, errors.NewResolveDependencyError(RabbitMQBackendLogIdentifier, "RabbitMQPublisher")
	}

	if container.Logger == nil {
		return nil, errors.NewResolveDependencyError(RabbitMQBackendLogIdentifier, "Logger")
	}

	return &RabbitMQBackend{
		cfg:       container.EventPublisherConfig,
		publisher: container.RabbitMQPublisher,
		logger:    container.Logger,
	}, nil
}
