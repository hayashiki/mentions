package main

import (
	"fmt"
	"github.com/hayashiki/mentions/pkg/config"
	"github.com/hayashiki/mentions/pkg/github"
	"github.com/hayashiki/mentions/pkg/handler"
	"github.com/hayashiki/mentions/pkg/repository"
	"log"
	"net/http"
	"os"
)

func main() {
	conf := config.NewConfigFromEnv()
	ghSvc := github.NewClient(github.GetClient(conf.GithubSecretToken))
	taskRepo := repository.NewTaskRepository(repository.GetClient(conf.GCPProject))
	teamRepo := repository.NewTeamRepository(repository.GetClient(conf.GCPProject))
	userRepo := repository.NewUserRepository(repository.GetClient(conf.GCPProject))
	repoRepo := repository.NewRepoRepository(repository.GetClient(conf.GCPProject))
	installationRepo := repository.NewInstallationRepository(repository.GetClient(conf.GCPProject))
	ghAppCli := github.NewClient(github.GetAppClient(conf.GithubAppID, int64(13344101), conf.GithubAppPrivateKeyFileName))

	h := handler.NewApp(conf, ghSvc, userRepo, teamRepo, repoRepo, taskRepo, installationRepo, ghAppCli)
	r := h.Handler()
	port := os.Getenv("PORT")
	if os.Getenv("PORT") == "" {
		port = "8000"
	}

	log.Printf("connect to port:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
