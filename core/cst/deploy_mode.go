package cst

type DeployMode = string

const (
	DeployModeProduction  DeployMode = "production"
	DeployModeStaging     DeployMode = "staging"
	DeployModeDevelopment DeployMode = "development"
)
