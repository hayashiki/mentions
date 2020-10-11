package main

import (
	"fmt"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/hayashiki/mentions/config"
	"github.com/hayashiki/mentions/handler"
	"github.com/hayashiki/mentions/infrastructure"
	"github.com/hayashiki/mentions/notifier"
	"github.com/hayashiki/mentions/repository"
	"github.com/hayashiki/mentions/usecase"
	"github.com/hayashiki/mentions/utils"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(utils.LoadEnvVariable())

	if err != nil {
		panic(err)
	}

	env := config.NewEnvironment("production")

	ghSvc := repository.NewClient(infrastructure.NewClient(env.GithubSecretToken))
	slackSvc := notifier.NewSlackNotifier()
	taskRepo := repository.NewTaskRepository(infrastructure.GetDSClient(env.GCPProject))

	uc := usecase.NewWebhookProcess(env, ghSvc, slackSvc, taskRepo)

	ghHandler := handler.NewWebhookHandler(uc)

	r := chi.NewRouter()
	r.Use(
		chiMiddleware.Logger,
		chiMiddleware.Recoverer,
	)
	r.Post("/webhook/github", ghHandler.PostWebhook)

	port := os.Getenv("PORT")
	if os.Getenv("PORT") == "" {
		port = "8000"
	}

	log.Printf("connect to port:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
