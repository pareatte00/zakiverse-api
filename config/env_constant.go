package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/zakiverse/zakiverse-api/util/validator"
)

type ConfigConstant struct {
	CredentialPrefix string                    `mapstructure:"credential_prefix"`
	Application      ConfigConstantApplication `mapstructure:"application" validate:"required"`
	Database         ConfigConstantDatabase    `mapstructure:"database" validate:"required"`
	Outbound         ConfigConstantOutbound    `mapstructure:"outbound" validate:"required"`
}

type ConfigConstantApplication struct {
	Name               string   `mapstructure:"name" validate:"required"`
	DeployPort         int      `mapstructure:"deploy_port" validate:"required"`
	DeployMode         string   `mapstructure:"deploy_mode" validate:"required,oneof=development staging production"`
	Timezone           string   `mapstructure:"timezone" validate:"required"`
	BaseUrl            string   `mapstructure:"base_url" validate:"required"`
	CorsAllowOrigins   []string `mapstructure:"cors_allow_origins" validate:"required"`
	RateLimitPerSecond int      `mapstructure:"rate_limit_per_second" validate:"required"`
	Timeout            int      `mapstructure:"timeout" validate:"required"`
}

type ConfigConstantDatabase struct {
	Host                  string `mapstructure:"host" validate:"required"`
	Port                  int    `mapstructure:"port" validate:"required"`
	Name                  string `mapstructure:"name" validate:"required"`
	MaxOpenConnection     int    `mapstructure:"max_open_connection" validate:"required"`
	MaxIdleConnection     int    `mapstructure:"max_idle_connection" validate:"required"`
	MaxConnectionLifetime int    `mapstructure:"max_connection_lifetime" validate:"required"`
	MaxConnectionIdleTime int    `mapstructure:"max_connection_idle_time" validate:"required"`
	CreateBatchSize       int    `mapstructure:"create_batch_size" validate:"required"` // only for postgresql
}

type ConfigConstantOutbound struct {
	Jikan ConfigConstantOutboundJikan `mapstructure:"jikan" validate:"required"`
}

type ConfigConstantOutboundJikan struct {
	BaseUrl string `mapstructure:"base_url" validate:"required"`
}

func initConstant() ConfigConstant {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config_file")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to load config.json file: %v", err)
	}

	var c ConfigConstant
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("Failed to unmarshal configuration: %v", err)
	}

	isValidatePass, errorFieldList, validatorError := validator.Validate(c)
	if validatorError != nil {
		log.Fatalf("Failed to validate configuration: %v", validatorError)
	}
	if !isValidatePass {
		log.Fatalf("Failed to validate configuration: %v", errorFieldList)
	}

	return c
}
