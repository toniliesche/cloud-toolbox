package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
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

	result, _ := json.Marshal(&faasResult{Status: "success"})

	fmt.Print(string(result))
}
