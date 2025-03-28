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

package redis

import (
	"cloud-toolbox/internal/application/faas/models"
	"encoding/json"
	"time"
)

type FunctionExecution struct {
	Id            string                `json:"id"`
	Output        interface{}           `json:"output"`
	Status        string                `json:"status"`
	StatusUpdates []models.StatusUpdate `json:"status_updates"`
	Error         string                `json:"error"`
	Created       time.Time             `json:"created"`
	Updated       time.Time             `json:"updated"`
}

func (e *FunctionExecution) ToModel() *models.FunctionExecution {
	var output string
	var ok bool

	if output, ok = e.Output.(string); !ok {
		jsonString, _ := json.Marshal(e.Output)
		output = string(jsonString)
	}

	return &models.FunctionExecution{
		Id:            e.Id,
		Output:        output,
		Status:        e.Status,
		StatusUpdates: e.StatusUpdates,
		Error:         e.Error,
		Created:       e.Created,
		Updated:       e.Updated,
	}
}
