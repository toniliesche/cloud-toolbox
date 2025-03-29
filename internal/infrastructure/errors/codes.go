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

package errors

import "net/http"

const (
	ErrorCodeGenericError = 0

	/* startup errors */
	ErrorCodeEnvironmentVariableMissing          = 10000
	ErrorCodeEnvironmentVariableMustBeInteger    = 10001
	ErrorCodeEnvironmentVariableMustBeBoolean    = 10002
	ErrorCodeEnvironmentVariableMalformed        = 10003
	ErrorCodeConfigValueMissing                  = 10100
	ErrorCodeConfigSectionMissing                = 10101
	ErrorCodeConfigSectionInvalid                = 10102
	ErrorCodeConfigValueMustBeGreaterThanZero    = 10103
	ErrorCodeConfigValueMustBeGreaterThanOrEqual = 10104
	ErrorCodeConfigValueMustBeLessThanOrEqual    = 10105
	ErrorCodeConfigValueInvalid                  = 10106
	ErrorCodeConfigValueEmptyList                = 10107
	ErrorCodeContainerConfigMissing              = 10200
	ErrorCodeContainerMissingDependency          = 10201
	ErrorCodeContainerInvalidConfig              = 10202
	ErrorCodeContainerApplicationSetup           = 10203
	ErrorCodeContainerMissing                    = 10204
	ErrorCodeServerNoRoutes                      = 10300
	ErrorCodeServerRouteFailed                   = 10301
	ErrorCodeConfigFileReadingFailedError        = 10400
	ErrorCodeConfigFileParsingFailedError        = 10401
	ErrorCodeConfigCreationFailedError           = 10402
	ErrorCodeConfigValidationFailedError         = 10403

	/* runtime errors */
	ErrorCodeSystemFileSystemAccess = 11000

	/* request parsing errors */
	ErrorCodeRequestParsingFailed              = 12000
	ErrorCodeRequestParsingMissingPayloadField = 12001
	ErrorCodeRequestParsingEmptyPayloadField   = 12002
	ErrorCodeRequestParsingUnknownPayloadType  = 12003

	ErrorCodeDatabaseItemNotFound = 13000

	/* function as a service errors */
	ErrorCodeFaasExecutionNotFound = 20000
	ErrorCodeFaasLimitExceeded     = 20001
	ErrorCodeFaasExecutionFailed   = 20002
	ErrorCodeFaasExecutionTimedOut = 20003

	/* function as a service execution errors */
	ErrorCodeFaasCommandCouldNotBeParsed = 20100
	ErrorCodeFaasCommandCouldNotBeKilled = 20101
)

func MapToStatusCode(errorCode int) int {
	if errorCode < 10000 {
		return http.StatusInternalServerError
	}

	if errorCode < 12000 {
		return http.StatusInternalServerError
	}

	if errorCode < 13000 {
		return http.StatusBadRequest
	}

	if errorCode < 14000 {
		return http.StatusNotFound
	}

	return mapFaasErrorToStatusCode(errorCode)
}

func mapFaasErrorToStatusCode(errorCode int) int {
	switch errorCode {
	case ErrorCodeFaasExecutionNotFound:
		return http.StatusNotFound
	case ErrorCodeFaasLimitExceeded:
		return http.StatusTooManyRequests
	case ErrorCodeFaasExecutionFailed:
		return http.StatusInternalServerError
	case ErrorCodeFaasExecutionTimedOut:
		return http.StatusRequestTimeout
	default:
		return http.StatusInternalServerError
	}
}
