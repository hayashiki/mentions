package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	GithubWebhookSecret string `envconfig:"GH_WEBHOOK_SECRET" required:"true"`
	GithubSecretToken   string `envconfig:"GH_SECRET_TOKEN" required:"true"`
	GCPProject          string `envconfig:"PROJECT" required:"true"`
}

func NewMustEnvironment() Environment {
	env := Environment{}
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalln("[ERROR] Can not read environment variables ", err)
		panic(err)
	}
	return env
}
