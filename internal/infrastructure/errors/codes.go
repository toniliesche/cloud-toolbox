package errors

import "net/http"

const (
	ErrorCodeGenericError = 0

	/* startup errors */
	ErrorCodeEnvironmentVariableMissing          = 10000
	ErrorCodeEnvironmentVariableMustBeInteger    = 10001
	ErrorCodeEnvironmentVariableMustBeBoolean    = 10002
	ErrorCodeConfigValueMissing                  = 10100
	ErrorCodeConfigSectionMissing                = 10101
	ErrorCodeConfigSectionInvalid                = 10102
	ErrorCodeConfigValueMustBeGreaterThanZero    = 10103
	ErrorCodeConfigValueMustBeGreaterThanOrEqual = 10104
	ErrorCodeConfigValueMustBeLessThanOrEqual    = 10105
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
