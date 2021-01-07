package config

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	GCPProject        string `envconfig:"GCP_PROJECT" required:"true"`
	MemcachedServer   string `envconfig:"MEMCACHED_SERVER" required:"true"`
	MemcachedUsername string `envconfig:"MEMCACHED_USERNAME" required:"true"`
	MemcachedPassword string `envconfig:"MEMCACHED_PASSWORD" required:"true"`
	Github
	Slack
}

type Github struct {
	GithubAppID                 int64  `envconfig:"GH_APP_ID" required:"true"`
	GithubAppClientID           string `envconfig:"GH_APP_CLIENT_ID" required:"true"`
	GithubAppSecretID           string `envconfig:"GH_APP_CLIENT_SECRET" required:"true"`
	GithubAppPrivateKeyFileName string `envconfig:"GH_APP_PRIVATE_KEY_FILENAME" required:"true"`
	GithubWebhookSecret         string `envconfig:"GH_WEBHOOK_SECRET" required:"true"`
	GithubSecretToken           string `envconfig:"GH_SECRET_TOKEN" required:"true"`
}

type Slack struct {
	SlackClientID        string `envconfig:"SLACK_CLIENT_ID" required:"true"`
	SlackSecretID        string `envconfig:"SLACK_SECRET_ID" required:"true"`
	SlackRedirectURL     string `envconfig:"SLACK_REDIRECT_URL" required:"true"`
	SlackUserRedirectURL string `envconfig:"SLACK_USER_REDIRECT_URL" required:"true"`
}

func NewConfigFromEnv() Config {
	env := Config{}
	err := envconfig.Process("", &env)
	if err != nil {
		log.Errorf("Can not read environment variables ", err)
		panic(err)
	}
	return env
}
