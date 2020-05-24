package mentions

import (
	"github.com/gin-gonic/gin"
	"github.com/hayashiki/mentions/account"
	"github.com/hayashiki/mentions/config"
	"github.com/hayashiki/mentions/gh"
	"github.com/hayashiki/mentions/notifier"
	"log"
)

func main() {
	env := config.NewEnvironment("production")

	list, err := account.LoadAccountFromFile("github-config.json")
	if err != nil {
		log.Printf(err.Error())
		return
	}

	gVerifier := gh.NewGithubVerifier()
	sNotifier := notifier.NewSlackNotifier()

	githubWebhookController := gh.NewWebhookHandler(gVerifier, sNotifier, env, list)

	r := gin.New()
	//r.Use(gin.Recovery(), Log())


	g := r.Group("/webhook")
	g.POST("/github", func(c *gin.Context) { githubWebhookController.PostWebhook(c) })
}
