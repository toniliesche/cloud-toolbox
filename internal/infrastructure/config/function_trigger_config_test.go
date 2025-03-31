package config_test

import (
	"cloud-toolbox/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateFunctionTriggerConfigFailsOnMissingSystemConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.SystemConfig = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing system section error") {
		return
	}

	if !assert.Equal(t, "config section `system` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnInvalidSystemConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.SystemConfig = &config.SystemConfig{}

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid system section error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `system`", "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnMissingRabbitMQConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.RabbitMQ = nil

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing rabbitmq section error") {
		return
	}

	if !assert.Equal(t, "config section `rabbitmq` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnInvalidRabbitMQConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.RabbitMQ.Consumer.Host = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid rabbitmq section error") {
		return
	}

	if !assert.Contains(t, err.Error(), "Error in config section `rabbitmq`", "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnMissingTriggerSource(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.TriggerSource = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing trigger source error") {
		return
	}

	if !assert.Equal(t, "config value `trigger_source` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnInvalidTriggerSource(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.TriggerSource = "invalid-source"

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid trigger source error") {
		return
	}

	if !assert.Equal(t, "config value `publisher_source` must be one of [rabbitmq], but is `invalid-source`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnMissingTriggerName(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.TriggerName = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing trigger name error") {
		return
	}

	if !assert.Equal(t, "config value `trigger_name` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnMissingFaasHost(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.FaasHost = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing faas host error") {
		return
	}

	if !assert.Equal(t, "config value `faas_host` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnMissingFaasPort(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.FaasPort = 0

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing faas port error") {
		return
	}

	if !assert.Equal(t, "config value `faas_port` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnInvalidFaasPort(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.FaasPort = -1

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid faas port error") {
		return
	}

	if !assert.Equal(t, "config value `faas_port` must be greater than `0`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnMissingFaasFunctionName(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.FaasFunctionName = ""

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch missing faas function name error") {
		return
	}

	if !assert.Equal(t, "config value `faas_function_name` must exist", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnInvalidFaasTimeout(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.FaasTimeout = -1

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid faas timeout error") {
		return
	}

	if !assert.Equal(t, "config value `faas_timeout` must be greater than or equal to `10`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnInvalidFaasTimeoutTooHigh(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.FaasTimeout = 1000

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid faas timeout error") {
		return
	}

	if !assert.Equal(t, "config value `faas_timeout` must be less than or equal to `900`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnInvalidBatchSize(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.BatchSize = -1

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid batch size error") {
		return
	}

	if !assert.Equal(t, "config value `batch_size` must be greater than `0`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigFailsOnInvalidBatchTimeout(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")
	cfg.BatchTimeout = -1

	err := cfg.Validate()
	if !assert.Error(t, err, "did not catch invalid batch timeout error") {
		return
	}

	if !assert.Equal(t, "config value `batch_timeout` must be greater than `0`", err.Error(), "unexpected error message") {
		return
	}
}

func TestValidateFunctionTriggerConfigSucceedsOnValidConfig(t *testing.T) {
	t.Parallel()
	cfg := getValidFunctionTriggerConfig("rabbitmq")

	if err := cfg.Validate(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func getValidFunctionTriggerConfig(source string) *config.FunctionTriggerConfig {
	var rabbitMQConfig *config.RabbitMQConfig

	switch source {
	case "rabbitmq":
		rabbitMQConfig = getValidRabbitMQConfig(config.RabbitMQModeConsumer)
	}

	return &config.FunctionTriggerConfig{
		ApplicationConfig: config.ApplicationConfig{
			SystemConfig: getValidSystemConfig(),
		},
		RabbitMQ:         rabbitMQConfig,
		TriggerName:      "faas",
		TriggerSource:    source,
		BatchTimeout:     10,
		BatchSize:        10,
		FaasHost:         "localhost",
		FaasPort:         8080,
		FaasTimeout:      10,
		FaasSsl:          false,
		FaasFunctionName: "function",
	}
}
