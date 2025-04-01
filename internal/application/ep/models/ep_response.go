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

type EpResponse struct {
	EventId string
	Status  int
	Data    map[string]string
}

func (r *EpResponse) GetStatusCode() int {
	return r.Status
}

func (r *EpResponse) GetBody() interface{} {
	if _, ok := r.Data["event_id"]; !ok {
		r.Data["event_id"] = r.EventId
	}

	return r.Data
}

func NewEpResponse(eventId string, status int, data map[string]string) *EpResponse {
	if data == nil {
		data = make(map[string]string)
	}

	return &EpResponse{
		EventId: eventId,
		Status:  status,
		Data:    data,
	}
}
