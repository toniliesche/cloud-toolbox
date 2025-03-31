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

package config

import (
	"cloud-toolbox/internal/infrastructure/errors"
	"fmt"
)

type RabbitMQPublisherConfig struct {
	VHost     string                             `yaml:"vhost"`
	Host      string                             `yaml:"host"`
	Port      int64                              `yaml:"port"`
	Username  string                             `yaml:"username"`
	Password  string                             `yaml:"password"`
	Ssl       bool                               `yaml:"ssl"`
	Exchanges map[string]*RabbitMQExchangeConfig `yaml:"exchanges,omitempty"`
}

func (c *RabbitMQPublisherConfig) Addr() string {
	var proto string
	if c.Ssl {
		proto = "amqps"
	} else {
		proto = "amqp"
	}

	addr := fmt.Sprintf(
		"%s://%s:%s@%s:%d",
		proto,
		c.Username,
		c.Password,
		c.Host,
		c.Port,
	)

	if c.VHost != "/" {
		addr += fmt.Sprintf("/%s", c.VHost)
	}

	return addr
}

func (c *RabbitMQPublisherConfig) Validate(path string) errors.ApplicationError {
	if c.VHost == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.vhost", path))
	}

	if c.Host == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.host", path))
	}

	if c.Port < 1 {
		return errors.NewConfigValueNeedsToBeGreaterZeroError(fmt.Sprintf("%s.port", path))
	}

	if c.Username == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.username", path))
	}

	if c.Password == "" {
		return errors.NewMissingConfigValueError(fmt.Sprintf("%s.password", path))
	}

	if err := c.validateExchanges(fmt.Sprintf("%s.exchanges", path)); err != nil {
		return err
	}

	return nil
}

func (c *RabbitMQPublisherConfig) validateExchanges(path string) errors.ApplicationError {
	if c.Exchanges == nil {
		return errors.NewMissingConfigSectionError(path)
	}

	if len(c.Exchanges) == 0 {
		return errors.NewEmptyListConfigValueError(path)
	}

	for id, exchange := range c.Exchanges {
		if err := exchange.Validate(fmt.Sprintf("%s.%s", path, id)); err != nil {
			return err
		}
	}

	return nil
}

func (c *RabbitMQPublisherConfig) GetExchange(exchange string) (*RabbitMQExchangeConfig, errors.ApplicationError) {
	exchangeCfg, ok := c.Exchanges[exchange]
	if !ok {
		return nil, errors.NewMissingListConfigValueError("rabbitmq.cfg.exchanges", exchange)
	}

	return exchangeCfg, nil
}

func getDefaultRabbitMQPublisherConfig() *RabbitMQPublisherConfig {
	return &RabbitMQPublisherConfig{
		VHost:     rabbitMQDefaultVhost,
		Host:      rabbitMQDefaultHost,
		Port:      rabbitMQDefaultPort,
		Ssl:       rabbitMQDefaultSsl,
		Exchanges: make(map[string]*RabbitMQExchangeConfig),
	}
}

func getRabbitMQPublisherConfigFromEnvironment() (*RabbitMQPublisherConfig, errors.ApplicationError) {
	cfg := getDefaultRabbitMQPublisherConfig()

	host := GetEnvironmentString("RABBITMQ_PUBLISHER_HOST", "")
	if host == "" {
		host = GetEnvironmentString("RABBITMQ_HOST", rabbitMQDefaultHost)
	}

	if host == "" {
		return nil, errors.NewMissingEnvironmentVariableError("RABBITMQ_PUBLISHER_HOST")
	}

	cfg.Host = host

	port, err := GetEnvironmentInt("RABBITMQ_PUBLISHER_PORT", 0)
	if err != nil {
		return nil, err
	}

	if port < 0 {
		return nil, errors.NewConfigValueNeedsToBeGreaterZeroError("RABBITMQ_PUBLISHER_PORT")
	}

	if port == 0 {
		port, err = GetEnvironmentInt("RABBITMQ_PORT", rabbitMQDefaultPort)
		if err != nil {
			return nil, err
		}

		if port < 0 {
			return nil, errors.NewConfigValueNeedsToBeGreaterZeroError("RABBITMQ_PORT")
		}
	}
	cfg.Port = port

	vhost := GetEnvironmentString("RABBITMQ_PUBLISHER_VHOST", "")
	if vhost == "" {
		vhost = GetEnvironmentString("RABBITMQ_VHOST", rabbitMQDefaultVhost)
	}

	if vhost == "" {
		return nil, errors.NewMissingEnvironmentVariableError("RABBITMQ_PUBLISHER_VHOST")
	}

	cfg.VHost = vhost

	username := GetEnvironmentString("RABBITMQ_PUBLISHER_USERNAME", "")
	if username == "" {
		username = GetEnvironmentString("RABBITMQ_USERNAME", "")
	}

	if username == "" {
		return nil, errors.NewMissingEnvironmentVariableError("RABBITMQ_PUBLISHER_USERNAME")
	}

	cfg.Username = username

	password := GetEnvironmentString("RABBITMQ_PUBLISHER_PASSWORD", "")
	if password == "" {
		password = GetEnvironmentString("RABBITMQ_PASSWORD", "")
	}

	if password == "" {
		return nil, errors.NewMissingEnvironmentVariableError("RABBITMQ_PUBLISHER_PASSWORD")
	}

	cfg.Password = password

	if HasEnvironment("RABBITMQ_PUBLISHER_SSL") {
		ssl, err := GetEnvironmentBool("RABBITMQ_PUBLISHER_SSL", rabbitMQDefaultSsl)
		if err != nil {
			return nil, err
		}

		cfg.Ssl = ssl
	} else if HasEnvironment("RABBITMQ_SSL") {
		ssl, err := GetEnvironmentBool("RABBITMQ_SSL", rabbitMQDefaultSsl)
		if err != nil {
			return nil, err
		}

		cfg.Ssl = ssl
	}

	exchanges := GetEnvironmentString("RABBITMQ_EXCHANGES", "")
	if exchanges == "" {
		return nil, errors.NewMissingEnvironmentVariableError("RABBITMQ_EXCHANGES")
	}

	exchangeMap, err := parseRabbitMQExchanges(exchanges)
	if err != nil {
		return nil, err
	}

	cfg.Exchanges = exchangeMap

	return cfg, nil
}
