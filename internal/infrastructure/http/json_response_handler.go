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

package http

import (
	"cloud-toolbox/internal/infrastructure/http/models"
	"encoding/json"
	"net/http"
)

type JsonResponseHandler struct {
}

func (h *JsonResponseHandler) sendJsonResponse(writer http.ResponseWriter, request *http.Request, statusCode int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")

	var msg string
	if statusCode > 399 {
		msg = "error"
	} else {
		msg = "ok"
	}

	response := models.JsonResponse{
		Request: request.Context().Value("request_id").(string),
		Message: msg,
		Data:    data,
	}

	writer.WriteHeader(statusCode)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		h.sendErrorResponse(writer, request, http.StatusInternalServerError, err.Error())
	}
}

func (h *JsonResponseHandler) sendErrorResponse(writer http.ResponseWriter, request *http.Request, statusCode int, message string) {
	writer.Header().Set("Content-Type", "application/json")

	response := models.JsonResponse{
		Request: request.Context().Value("request_id").(string),
		Message: message,
	}

	writer.WriteHeader(statusCode)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
