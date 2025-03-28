package models

import (
	"encoding/json"
	"time"
)

type StatusUpdate struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

func (u *StatusUpdate) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}
