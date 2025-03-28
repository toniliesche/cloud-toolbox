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

package connectors

import (
	"cloud-toolbox/internal/infrastructure/di"
	"cloud-toolbox/internal/infrastructure/errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewScylla(container *di.Container) (*dynamodb.DynamoDB, errors.ApplicationError) {
	if container == nil {
		return nil, errors.NewContainerMissingError("Scylla")
	}

	if container.ScyllaConfig == nil {
		return nil, errors.NewResolveDependencyError("Scylla", "ScyllaConfig")
	}

	var proto string
	if container.ScyllaConfig.SSL == true {
		proto = "https"
	} else {
		proto = "http"
	}

	scyllaCreds := credentials.NewStaticCredentials(container.ScyllaConfig.Username, container.ScyllaConfig.Password, "")

	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(container.ScyllaConfig.Region),
			Endpoint:    aws.String(fmt.Sprintf("%s://%s:%d", proto, container.ScyllaConfig.Host, container.ScyllaConfig.Port)),
			Credentials: scyllaCreds,
		},
	)

	if err != nil {
		return nil, errors.NewGenericError(err)
	}

	return dynamodb.New(sess), nil
}
