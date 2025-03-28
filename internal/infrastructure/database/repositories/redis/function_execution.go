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
