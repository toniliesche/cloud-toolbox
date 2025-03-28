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

package models

import "time"

type FunctionExecution struct {
	Id            string         `json:"id"`
	Output        string         `json:"output"`
	Status        string         `json:"status"`
	StatusUpdates []StatusUpdate `json:"status_updates"`
	Error         string         `json:"error"`
	Created       time.Time      `json:"created"`
	Updated       time.Time      `json:"updated"`
}

func FunctionExecutionFromRecord(record *ExecutionRecord) *FunctionExecution {
	timeObj := time.Now()
	status := record.FunctionExecution.GetStatus()

	statusUpdates := make([]StatusUpdate, 0)
	statusUpdate := &StatusUpdate{
		Status: status,
		Time:   timeObj,
	}

	statusUpdates = append(statusUpdates, *statusUpdate)

	return &FunctionExecution{
		Id:            record.ExecutionId,
		Output:        "",
		Status:        status,
		StatusUpdates: statusUpdates,
		Error:         "",
		Created:       timeObj,
		Updated:       timeObj,
	}
}
