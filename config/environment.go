package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	GithubWebhookSecret string `envconfig:"GITHUB_WEBHOOK_SECRET" required:"true"`
}

func getProductionEnv() Environment {
	env := Environment{}
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatalln("[ERROR] Can not read environment variables ", err)
	}
	return env
}

func NewEnvironment(params ...string) Environment {
	if len(params) == 0 {
		log.Println("[INFO] using production environment")
		return getProductionEnv()
	}

	switch requiredEnv := params[0]; requiredEnv {
	case "production":
		log.Println("[INFO] using production environment")
		return getProductionEnv()
	case "development":
		log.Println("[INFO] using unit test environment")
		return getProductionEnv()
	default:
		log.Println("[INFO] using production environment")
		return getProductionEnv()
	}
}
