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
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"time"
)

func NewScyllaItemFromFunctionExecution(function *models.FunctionExecution) map[string]*dynamodb.AttributeValue {
	created := function.Created.Format(time.RFC3339)
	updated := function.Updated.Format(time.RFC3339)

	statusUpdates := make([]*dynamodb.AttributeValue, 0)

	for _, update := range function.StatusUpdates {
		statusUpdates = append(statusUpdates, &dynamodb.AttributeValue{
			M: map[string]*dynamodb.AttributeValue{
				"time": {
					S: aws.String(update.Time.Format(time.RFC3339)),
				},
				"status": {
					S: aws.String(update.Status),
				},
			},
		})
	}

	item := map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(function.Id),
		},
		"created": {
			S: aws.String(created),
		},
		"updated": {
			S: aws.String(updated),
		},
	}

	if function.Output != "" {
		item["output"] = &dynamodb.AttributeValue{
			S: aws.String(function.Output),
		}
	}

	if function.Error != "" {
		item["error"] = &dynamodb.AttributeValue{
			S: aws.String(function.Error),
		}
	}

	if function.Status != "" {
		item["status"] = &dynamodb.AttributeValue{
			S: aws.String(function.Status),
		}
	}

	item["status_updates"] = &dynamodb.AttributeValue{
		L: statusUpdates,
	}

	return item
}

func NewFunctionExecutionFromScyllaItem(item map[string]*dynamodb.AttributeValue) *models.FunctionExecution {
	statusUpdates := make([]models.StatusUpdate, 0)

	id := *item["id"].S

	if _, ok := item["status_updates"]; ok {
		for _, update := range item["status_updates"].L {
			var timeObj time.Time
			if _, ok := update.M["time"]; ok {
				timeObj, _ = time.Parse(time.RFC3339, *update.M["time"].S)
			}

			var status string
			if _, ok := update.M["status"]; ok {
				status = *update.M["status"].S
			}

			statusUpdates = append(statusUpdates, models.StatusUpdate{
				Status: status,
				Time:   timeObj,
			})
		}
	}

	var output string
	if _, ok := item["output"]; ok {
		output = *item["output"].S
	}

	var status string
	if _, ok := item["status"]; ok {
		status = *item["status"].S
	}

	var errorMessage string
	if _, ok := item["error"]; ok {
		errorMessage = *item["error"].S
	}

	var created time.Time
	if _, ok := item["created"]; ok {
		created, _ = time.Parse(time.RFC3339, *item["created"].S)
	}

	var updated time.Time
	if _, ok := item["updated"]; ok {
		updated, _ = time.Parse(time.RFC3339, *item["updated"].S)
	}

	if _, ok := item["ttl"]; ok {
		ttl := *item["ttl"].N

		fmt.Printf("ttl : %s\n", ttl)
	} else {
		fmt.Printf("No ttl found\n")
	}

	return &models.FunctionExecution{
		Id:            id,
		Output:        output,
		Status:        status,
		StatusUpdates: statusUpdates,
		Error:         errorMessage,
		Created:       created,
		Updated:       updated,
	}
}
