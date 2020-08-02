package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hayashiki/mentions/account"
	"github.com/hayashiki/mentions/config"
	"github.com/hayashiki/mentions/gh"
	"github.com/hayashiki/mentions/handler"
	"github.com/hayashiki/mentions/notifier"
	"github.com/hayashiki/mentions/utils"
	"log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(utils.LoadEnvVariable())

	env := config.NewEnvironment("production")

	list, err := account.LoadAccountFromFile("github-config.json")
	if err != nil {
		log.Printf(err.Error())
		return
	}

	gVerifier := gh.NewGithubVerifier()
	sNotifier := notifier.NewSlackNotifier(list)
	githubWebhookController := handler.NewWebhookHandler(gVerifier, sNotifier, env, list)

	r := gin.Default()
	r.Use(gin.Recovery())

	g := r.Group("/webhook")
	g.POST("/github", func(c *gin.Context) { githubWebhookController.PostWebhook(c) })
	r.Run()
}
