package config

import "github.com/zakiverse/zakiverse-api/core/cst"

func InitConfig() (ConfigConstant, ConfigCredential) {
	conf := initConstant()
	if conf.Application.DeployMode == cst.DeployModeDevelopment {
		loadEnvFile()
	}
	credential := initCredential(conf.CredentialPrefix)

	setGlobalTimezone(conf.Application.Timezone)

	if conf.Application.DeployMode == cst.DeployModeProduction || conf.Application.DeployMode == cst.DeployModeStaging {
		setGinReleaseMode()
	}

	return conf, credential
}
