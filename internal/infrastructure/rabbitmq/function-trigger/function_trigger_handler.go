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

package function_trigger

import (
	"bytes"
	"cloud-toolbox/internal/application/faas/models"
	"cloud-toolbox/internal/infrastructure/config"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"cloud-toolbox/internal/infrastructure/http/client"
	"cloud-toolbox/internal/infrastructure/rabbitmq/model"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"io"
	"time"
)

const (
	FunctionTriggerHandlerLogIdentifier = "FunctionTriggerHandler"
)

type FunctionTriggerHandler struct {
	logger *zerolog.Logger
	cfg    *config.FunctionTriggerConfig
	client *client.HttpClient
}

func (f *FunctionTriggerHandler) QueueIdentifier() string {
	return "trigger"
}

func (f *FunctionTriggerHandler) HandleMessageBatch(msg []amqp091.Delivery) ([]string, error) {
	records := make([]*model.RabbitMQRecord, 0, len(msg))

	for _, m := range msg {
		f.logger.Trace().
			Msgf("[%s] Received message: %s", FunctionTriggerHandlerLogIdentifier, m.MessageId)

		record := model.NewRabbitMQRecordFromDelivery(m)
		records = append(records, record)

		f.logger.Trace().
			Msgf("[%s] Processing message: %v", FunctionTriggerHandlerLogIdentifier, m.MessageId)
	}

	faasRequest := &models.FaasRequest[*model.RabbitMQRecord]{
		FaasRequestSource: models.FaasRequestSource{
			Source: "ctb:rmq",
		},
		Async:   false,
		Records: records,
	}

	requestBytes, _ := json.Marshal(faasRequest)
	buffer := bytes.NewBuffer(requestBytes)
	requestUri := fmt.Sprintf("/cloud-toolbox/faas/%s", f.cfg.FaasFunctionName)

	f.logger.Trace().
		Str("uri", requestUri).
		Bytes("request", requestBytes).
		Msgf("[%s] Sending request to function: %s", FunctionTriggerHandlerLogIdentifier, f.cfg.FaasFunctionName)

	req, err := f.client.NewRequest("POST", requestUri, buffer)
	if err != nil {
		f.logger.Error().
			Err(err).
			Msgf("[%s] Error creating request: %s", FunctionTriggerHandlerLogIdentifier, err.Error())
		return nil, errors.NewGenericError(err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		f.logger.Error().
			Err(err).
			Msgf("[%s] Error sending request: %s", FunctionTriggerHandlerLogIdentifier, err.Error())
		return nil, errors.NewGenericError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		f.logger.Error().
			Msgf("[%s] Error response from function: %s", FunctionTriggerHandlerLogIdentifier, resp.Status)
		return nil, errors.NewGenericError(fmt.Errorf("error response from function: %s", resp.Status))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		f.logger.Error().
			Err(err).
			Msgf("[%s] Error reading response: %s", FunctionTriggerHandlerLogIdentifier, err.Error())
		return nil, errors.NewGenericError(err)
	}

	f.logger.Trace().
		Bytes("response", respBody).
		Msgf("[%s] Received response from function", FunctionTriggerHandlerLogIdentifier)

	var faasResponse models.FaasResponse
	if err := json.Unmarshal(respBody, &faasResponse); err != nil {
		f.logger.Error().
			Err(err).
			Msgf("[%s] Error unmarshalling response: %s", FunctionTriggerHandlerLogIdentifier, err.Error())
		return nil, errors.NewGenericError(err)
	}

	body := faasResponse.GetBody()
	bodyMap := body.(map[string]interface{})

	output := bodyMap["output"].(string)
	responseOutput := &models.FaasResponseOutput{}
	if err := json.Unmarshal([]byte(output), responseOutput); err != nil {
		f.logger.Error().
			Err(err).
			Msgf("[%s] Error unmarshalling response output: %s", FunctionTriggerHandlerLogIdentifier, err.Error())
		return nil, errors.NewGenericError(err)
	}

	return responseOutput.FailedRecords, nil
}

func NewFunctionTriggerHandler(container *di.Container) (*FunctionTriggerHandler, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError(FunctionTriggerHandlerLogIdentifier)
	}

	if container.Logger == nil {
		return nil, errors.NewResolveDependencyError(FunctionTriggerHandlerLogIdentifier, "Logger")
	}

	if container.FunctionTriggerConfig == nil {
		return nil, errors.NewResolveDependencyError(FunctionTriggerHandlerLogIdentifier, "FunctionTriggerConfig")
	}

	ftConfig := container.FunctionTriggerConfig
	if err := ftConfig.Validate(); err != nil {
		return nil, errors.NewInvalidConfigError(FunctionTriggerHandlerLogIdentifier, err)
	}

	var proto string
	if ftConfig.FaasSsl {
		proto = "https"
	} else {
		proto = "http"
	}

	addr := fmt.Sprintf("%s://%s:%d", proto, ftConfig.FaasHost, ftConfig.FaasPort)

	timeout := time.Duration(ftConfig.FaasTimeout) * time.Second
	httpClient, err := client.NewHttpClient(addr, timeout)
	if err != nil {
		return nil, errors.NewResolveDependencyError(FunctionTriggerHandlerLogIdentifier, "HttpClient")
	}

	return &FunctionTriggerHandler{
		cfg:    ftConfig,
		logger: container.Logger,
		client: httpClient,
	}, nil
}
