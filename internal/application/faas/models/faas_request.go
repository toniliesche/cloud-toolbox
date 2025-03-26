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

import (
	faasinterfaces "cloud-toolbox/internal/application/faas/models/interfaces"
	"cloud-toolbox/internal/domain/errors"
)

type FaasRequest[K faasinterfaces.FaasRecord] struct {
	FaasRequestSource
	Async   bool `json:"async"`
	Records []K  `json:"records"`
}

func (r *FaasRequest[K]) Validate() error {
	if r.Records == nil {
		return errors.NewMissingRequestFieldError("records")
	}

	if len(r.Records) == 0 {
		return errors.NewEmptyRequestFieldError("records")
	}

	if r.Source == "" {
		return errors.NewMissingRequestFieldError("source")
	}

	for _, record := range r.Records {
		if err := record.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (r *FaasRequest[K]) ToAbstractRequest() FaasRequest[faasinterfaces.FaasRecord] {
	request := FaasRequest[faasinterfaces.FaasRecord]{
		FaasRequestSource: r.FaasRequestSource,
		Async:             r.Async,
		Records:           make([]faasinterfaces.FaasRecord, 0),
	}

	for _, record := range r.Records {
		request.Records = append(request.Records, record)
	}

	return request
}

type FaasRequestSource struct {
	Source string `json:"source"`
}

func (s FaasRequestSource) Validate() error {
	if s.Source == "" {
		return errors.NewMissingRequestFieldError("source")
	}

	return nil
}
