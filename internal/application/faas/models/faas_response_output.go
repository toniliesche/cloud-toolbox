package models

type FaasResponseOutput struct {
	FailedRecords []string `json:"failed_records,omitempty"`
}
