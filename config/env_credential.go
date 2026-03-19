package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type ConfigCredential struct {
	DatabaseUsername string `required:"true" split_words:"true"`
	DatabasePassword string `required:"true" split_words:"true"`
	SystemServiceKey   string `required:"true" split_words:"true"`
	JwtSecret          string `required:"true" split_words:"true"`
	DiscordClientId    string `required:"true" split_words:"true"`
	DiscordClientSecret string `required:"true" split_words:"true"`
}

func initCredential(prefix string) ConfigCredential {
	var c ConfigCredential
	err := envconfig.Process(prefix, &c)
	if err != nil {
		log.Fatalf("Failed to load environment data: %v", err)
	}

	return c
}
