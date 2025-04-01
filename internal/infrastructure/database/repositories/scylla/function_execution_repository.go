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

package scylla

import (
	"cloud-toolbox/internal/application/faas/models"
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
	"time"
)

const FunctionExecutionRepositoryLogIdentifier = "FunctionExecutionRepository/Scylla"

type FunctionExecutionRepository struct {
	table      string
	logger     *zerolog.Logger
	timeToLive time.Duration
	context    context.Context
	scylla     *dynamodb.DynamoDB
}

func (f *FunctionExecutionRepository) GetFunction(executionId string) (*models.FunctionExecution, errors.ApplicationError) {
	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Getting function execution", FunctionExecutionRepositoryLogIdentifier)

	get := &dynamodb.GetItemInput{
		TableName: aws.String(f.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(executionId),
			},
		},
	}

	req, out := f.scylla.GetItemRequest(get)
	if err := req.Send(); err != nil {
		f.logger.Trace().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Error getting function execution", FunctionExecutionRepositoryLogIdentifier)

		return nil, errors.NewGenericError(err)
	}

	if out.Item == nil {
		f.logger.Trace().
			Str("execution-id", executionId).
			Msgf("[%s] Function execution not found", FunctionExecutionRepositoryLogIdentifier)

		return nil, errors.NewItemNotFoundError(executionId)
	}

	fn := NewFunctionExecutionFromScyllaItem(out.Item)

	result, _ := json.Marshal(fn)

	f.logger.Trace().
		Str("execution-id", executionId).
		Str("json", string(result)).
		Msgf("[%s] Function execution found", FunctionExecutionRepositoryLogIdentifier)

	return fn, nil
}

func (f *FunctionExecutionRepository) SaveError(executionId string, error string) errors.ApplicationError {
	timeObj := time.Now()

	ttl := fmt.Sprintf("%d", timeObj.Add(f.timeToLive).Unix())
	f.logger.Trace().
		Str("execution-id", executionId).
		Str("ttl", ttl).
		Msgf("[%s] Saving error", FunctionExecutionRepositoryLogIdentifier)

	update := dynamodb.UpdateItemInput{
		TableName: aws.String(f.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(executionId),
			},
		},
		UpdateExpression: aws.String("SET #error = :error, #updated = :updated, #ttl = :ttl"),
		ExpressionAttributeNames: map[string]*string{
			"#error":   aws.String("error"),
			"#updated": aws.String("updated"),
			"#ttl":     aws.String("ttl"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":error": {
				S: aws.String(error),
			},
			":updated": {
				S: aws.String(timeObj.Format(time.RFC3339Nano)),
			},
			":ttl": {
				N: aws.String(ttl),
			},
		},
	}

	_, err := f.scylla.UpdateItem(&update)
	if err != nil {
		f.logger.Trace().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Error saving error", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(err)
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Error saved", FunctionExecutionRepositoryLogIdentifier)

	return nil
}

func (f *FunctionExecutionRepository) SaveFunction(function *models.FunctionExecution) errors.ApplicationError {
	ttl := fmt.Sprintf("%d", time.Now().Add(f.timeToLive).Unix())
	f.logger.Trace().
		Str("execution-id", function.Id).
		Str("ttl", ttl).
		Msgf("[%s] Saving function", FunctionExecutionRepositoryLogIdentifier)

	item := NewScyllaItemFromFunctionExecution(function)
	item["ttl"] = &dynamodb.AttributeValue{
		N: aws.String(ttl),
	}

	put := &dynamodb.PutItemInput{
		TableName: aws.String(f.table),
		Item:      item,
	}

	_, err := f.scylla.PutItem(put)
	if err != nil {
		f.logger.Trace().
			Err(err).
			Str("execution-id", function.Id).
			Msgf("[%s] Error saving function", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(err)
	}

	f.logger.Trace().
		Str("execution-id", function.Id).
		Msgf("[%s] Function saved", FunctionExecutionRepositoryLogIdentifier)

	return nil
}

func (f *FunctionExecutionRepository) SaveOutput(executionId string, output string) errors.ApplicationError {
	timeObj := time.Now()

	ttl := fmt.Sprintf("%d", timeObj.Add(f.timeToLive).Unix())
	f.logger.Trace().
		Str("execution-id", executionId).
		Str("ttl", ttl).
		Msgf("[%s] Saving output", FunctionExecutionRepositoryLogIdentifier)

	update := dynamodb.UpdateItemInput{
		TableName: aws.String(f.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(executionId),
			},
		},
		UpdateExpression: aws.String("SET #output = :output, #updated = :updated, #ttl = :ttl"),
		ExpressionAttributeNames: map[string]*string{
			"#output":  aws.String("output"),
			"#updated": aws.String("updated"),
			"#ttl":     aws.String("ttl"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":output": {
				S: aws.String(output),
			},
			":updated": {
				S: aws.String(timeObj.Format(time.RFC3339Nano)),
			},
			":ttl": {
				N: aws.String(ttl),
			},
		},
	}

	_, err := f.scylla.UpdateItem(&update)
	if err != nil {
		f.logger.Trace().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Error saving output", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(err)
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Output saved", FunctionExecutionRepositoryLogIdentifier)

	return nil
}

func (f *FunctionExecutionRepository) UpdateStatus(executionId string, status string) errors.ApplicationError {
	timeObj := time.Now()

	ttl := fmt.Sprintf("%d", timeObj.Add(f.timeToLive).Unix())
	f.logger.Trace().
		Str("execution-id", executionId).
		Str("ttl", ttl).
		Str("status", status).
		Msgf("[%s] Updating status", FunctionExecutionRepositoryLogIdentifier)

	statusUpdate := map[string]*dynamodb.AttributeValue{
		"time": {
			S: aws.String(timeObj.Format(time.RFC3339Nano)),
		},
		"status": {
			S: aws.String(status),
		},
	}

	update := dynamodb.UpdateItemInput{
		TableName: aws.String(f.table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(executionId),
			},
		},
		UpdateExpression: aws.String("SET #status = :status, #updated = :updated, #status_updates = list_append(#status_updates, :status_updates), #ttl = :ttl"),
		ExpressionAttributeNames: map[string]*string{
			"#status":         aws.String("status"),
			"#updated":        aws.String("updated"),
			"#status_updates": aws.String("status_updates"),
			"#ttl":            aws.String("ttl"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":status": {
				S: aws.String(status),
			},
			":updated": {
				S: aws.String(timeObj.Format(time.RFC3339Nano)),
			},
			":status_updates": {
				L: []*dynamodb.AttributeValue{
					{
						M: statusUpdate,
					},
				},
			},
			":ttl": {
				N: aws.String(ttl),
			},
		},
	}

	_, err := f.scylla.UpdateItem(&update)
	if err != nil {
		f.logger.Trace().
			Err(err).
			Str("execution-id", executionId).
			Msgf("[%s] Error updating status", FunctionExecutionRepositoryLogIdentifier)

		return errors.NewGenericError(err)
	}

	f.logger.Trace().
		Str("execution-id", executionId).
		Msgf("[%s] Status updated", FunctionExecutionRepositoryLogIdentifier)

	return nil
}

func NewFunctionExecutionRepository(container *di.Container) (*FunctionExecutionRepository, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError(FunctionExecutionRepositoryLogIdentifier)
	}

	if container.Context == nil {
		return nil, errors.NewResolveDependencyError(FunctionExecutionRepositoryLogIdentifier, "Context")
	}

	if container.Scylla == nil {
		return nil, errors.NewResolveDependencyError(FunctionExecutionRepositoryLogIdentifier, "Scylla")
	}

	if container.ScyllaConfig == nil {
		return nil, errors.NewResolveDependencyError(FunctionExecutionRepositoryLogIdentifier, "ScyllaConfig")
	}

	if err := container.ScyllaConfig.Validate("scylla"); err != nil {
		return nil, errors.NewInvalidConfigError(FunctionExecutionRepositoryLogIdentifier, err)
	}

	if container.Logger == nil {
		return nil, errors.NewResolveDependencyError(FunctionExecutionRepositoryLogIdentifier, "Logger")
	}

	if container.FunctionAsAServiceConfig == nil {
		return nil, errors.NewResolveDependencyError(FunctionExecutionRepositoryLogIdentifier, "FunctionAsAServiceConfig")
	}

	if err := container.FunctionAsAServiceConfig.Validate(); err != nil {
		return nil, errors.NewInvalidConfigError(FunctionExecutionRepositoryLogIdentifier, err)
	}

	return &FunctionExecutionRepository{
		context:    container.Context,
		scylla:     container.Scylla,
		logger:     container.Logger,
		table:      container.ScyllaConfig.Table,
		timeToLive: time.Duration(container.FunctionAsAServiceConfig.StorageTtl) * time.Second,
	}, nil
}
