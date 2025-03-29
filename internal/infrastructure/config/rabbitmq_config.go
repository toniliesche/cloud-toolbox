package config

import (
	"cloud-toolbox/internal/infrastructure/errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	rabbitMQDefaultHost             = "localhost"
	rabbitMQDefaultPort             = 5672
	rabbitMQDefaultSsl              = false
	rabbitMQRegexp                  = `^([a-zA-Z_-]+)=([^\x00\s]{1,255})\/([a-zA-Z0-9_.:#*-]{1,255})$`
	RabbitMQModePublisher           = 1
	RabbitMQModeSubscriber          = 2
	RabbitMQModePublisherSubscriber = 3
)

type RabbitMQConfig struct {
	Host      string                            `yaml:"host"`
	Port      int64                             `yaml:"port"`
	Username  string                            `yaml:"username"`
	Password  string                            `yaml:"password"`
	Ssl       bool                              `yaml:"ssl"`
	Exchanges map[string]RabbitMQExchangeConfig `yaml:"exchanges,omitempty"`
	Queues    map[string]RabbitMQQueueConfig    `yaml:"queues,omitempty"`
}

func (c *RabbitMQConfig) Validate(path string, mode int) errors.ApplicationError {
	if c.Host == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.host", path))
	}

	if c.Port == 0 {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.port", path))
	}

	if c.Port < 0 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError(fmt.Sprintf("%s.port", path))
	}

	if c.Username == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.username", path))
	}

	if c.Password == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.password", path))
	}

	if (mode & RabbitMQModePublisher) == RabbitMQModePublisher {
		if err := c.validateExchanges(path); err != nil {
			return err
		}
	}

	if (mode & RabbitMQModeSubscriber) == RabbitMQModeSubscriber {
		if err := c.validateQueues(path); err != nil {
			return err
		}
	}

	return nil
}

func (c *RabbitMQConfig) validateExchanges(path string) errors.ApplicationError {
	if c.Exchanges == nil {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.exchanges", path))
	}

	if len(c.Exchanges) == 0 {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.exchanges", path))
	}

	for id, exchange := range c.Exchanges {
		if err := exchange.Validate(fmt.Sprintf("%s.exchanges.%s", path, id)); err != nil {
			return err
		}
	}

	return nil
}

func (c *RabbitMQConfig) validateQueues(path string) errors.ApplicationError {
	if c.Queues == nil {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.queues", path))
	}

	if len(c.Queues) == 0 {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.queues", path))
	}

	for id, queue := range c.Queues {
		if err := queue.Validate(fmt.Sprintf("%s.queues.%s", path, id)); err != nil {
			return err
		}
	}

	return nil
}

type RabbitMQExchangeConfig struct {
	VHost string `yaml:"vhost"`
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
}

func (c *RabbitMQExchangeConfig) Validate(path string) errors.ApplicationError {
	if c.VHost == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.vhost", path))
	}

	if c.Name == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.name", path))
	}

	if c.Type == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.type", path))
	}

	return nil
}

type RabbitMQQueueConfig struct {
	VHost string `yaml:"vhost"`
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
}

func (c *RabbitMQQueueConfig) Validate(path string) errors.ApplicationError {
	if c.VHost == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.vhost", path))
	}

	if c.Name == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.name", path))
	}

	if c.Type == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.type", path))
	}

	return nil
}

func getDefaultRabbitMQConfig() *RabbitMQConfig {
	return &RabbitMQConfig{
		Host: rabbitMQDefaultHost,
		Port: rabbitMQDefaultPort,
		Ssl:  rabbitMQDefaultSsl,
	}
}

func getRabbitMQConfigFromEnvironment() (*RabbitMQConfig, errors.ApplicationError) {
	cfg := getDefaultRabbitMQConfig()

	host := GetEnvironmentString("RABBITMQ_HOST", rabbitMQDefaultHost)
	if host != "" {
		cfg.Host = host
	}

	port, err := GetEnvironmentInt("RABBITMQ_PORT", rabbitMQDefaultPort)
	if err != nil {
		return nil, err
	}

	if port != 0 {
		cfg.Port = port
	}

	username := GetEnvironmentString("RABBITMQ_USERNAME", "")
	if username != "" {
		cfg.Username = username
	}

	password := GetEnvironmentString("RABBITMQ_PASSWORD", "")
	if password != "" {
		cfg.Password = password
	}

	ssl, err := GetEnvironmentBool("RABBITMQ_SSL", rabbitMQDefaultSsl)
	if err != nil {
		return nil, err
	}

	if ssl != rabbitMQDefaultSsl {
		cfg.Ssl = ssl
	}

	queues := GetEnvironmentString("RABBITMQ_QUEUES", "")
	if queues != "" {
		cfg.Queues, err = parseRabbitMQQueues(queues)
	}
	if err != nil {
		return nil, err
	}

	exchanges := GetEnvironmentString("RABBITMQ_EXCHANGES", "")
	if exchanges != "" {
		cfg.Exchanges, err = parseRabbitMQExchanges(exchanges)
	}
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseRabbitMQQueues(queues string) (map[string]RabbitMQQueueConfig, errors.ApplicationError) {
	queueMap := make(map[string]RabbitMQQueueConfig)

	re := regexp.MustCompile(rabbitMQRegexp)

	queueList := strings.Split(queues, ",")
	for _, queue := range queueList {
		matches := re.FindStringSubmatch(queue)
		if matches == nil {
			return nil, errors.NewMalformedEnvironmentVariableError("RABBITMQ_QUEUES", queue, rabbitMQRegexp)
		}

		queue := RabbitMQQueueConfig{
			VHost: matches[2],
			Name:  matches[3],
			Type:  "direct",
		}

		queueMap[matches[1]] = queue
	}

	return queueMap, nil
}

func parseRabbitMQExchanges(exchanges string) (map[string]RabbitMQExchangeConfig, errors.ApplicationError) {
	exchangeMap := make(map[string]RabbitMQExchangeConfig)

	re := regexp.MustCompile(rabbitMQRegexp)

	exchangeList := strings.Split(exchanges, ",")
	for _, exchange := range exchangeList {
		matches := re.FindStringSubmatch(exchange)
		if matches == nil {
			return nil, errors.NewMalformedEnvironmentVariableError("RABBITMQ_EXCHANGES", exchange, rabbitMQRegexp)
		}

		exchange := RabbitMQExchangeConfig{
			VHost: matches[2],
			Name:  matches[3],
			Type:  "direct",
		}

		exchangeMap[matches[1]] = exchange
	}

	return exchangeMap, nil
}
