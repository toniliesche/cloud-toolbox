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

package main

import (
	"cloud-toolbox/internal/infrastructure/rabbitmq/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type faasResult struct {
	Status        string   `json:"status"`
	FailedRecords []string `json:"failed_records,omitempty"`
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	log.SetOutput(os.Stderr)
	log.Println("hello world started")
	log.Printf("input: %s\n", string(input))

	records := make([]*model.RabbitMQRecord, 0)
	if err := json.Unmarshal(input, &records); err != nil {
		log.Println("Error unmarshalling input:", err)
		return
	}

	failed := make([]string, 0)
	for _, record := range records {
		failed = append(failed, record.GetRecordIdentifier())
	}

	faasResult := &faasResult{
		FailedRecords: failed,
		Status:        "success",
	}

	result, _ := json.Marshal(faasResult)

	fmt.Print(string(result))

	time.Sleep(5 * time.Second)
}
