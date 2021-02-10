package usecase

import (
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/pkg/config"
	ghSvc "github.com/hayashiki/mentions/pkg/github"
	"github.com/hayashiki/mentions/pkg/repository"
	"log"
	"net/http"
)

type WebhookProcess interface {
	Do(r *http.Request) error
}

type webhookProcess struct {
	config         config.Config
	ghSvc          ghSvc.Github
	ghAppSvc       ghSvc.Github
	ghReqValidator ghSvc.RequestValidator
	userRepo       repository.UserRepository
	taskRepo       repository.TaskRepository
	repoRepo       repository.RepoRepository
	teamRepo       repository.TeamRepository
}

func NewWebhookProcess(
	env config.Config,
	ghSvc ghSvc.Github,
	ghAppSvc ghSvc.Github,
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
	repoRepo repository.RepoRepository,
) WebhookProcess {
	return &webhookProcess{
		config:   env,
		ghSvc:    ghSvc,
		ghAppSvc: ghAppSvc,
		userRepo: userRepo,
		taskRepo: taskRepo,
		repoRepo: repoRepo,
	}
}

func (w webhookProcess) Do(r *http.Request) error {
	// TODO: get from secret manager
	//payload, err := ghSvc.Verify(r, []byte(w.config.GithubWebhookSecret))
	// TODO: webhookProcessのinitでnewする
	payload, err := ghSvc.NewGithubValidator().Validate(r, []byte(w.config.GithubWebhookSecret))
	if err != nil {
		return err
	}

	ghEvent, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return err
	}

	switch ev := ghEvent.(type) {
	case *github.IssueCommentEvent:
		switch ev.GetAction() {
		case "created":
			return w.processIssueComment(r.Context(), ev)
		case "edited":
			return w.processEditIssueComment(r.Context(), ev)
		}
	case *github.PullRequestReviewCommentEvent:
		return w.processPullRequestComment(r.Context(), ev)
	case *github.InstallationRepositoriesEvent:
		log.Printf("action is %v", ev.GetAction())
		switch ev.GetAction() {
		case "added":
			return w.processInstallationReposAddedEvent(r.Context(), ev)
		case "removed":
			return w.processInstallationReposRemovedEvent(r.Context(), ev)
		}
	case *github.InstallationEvent:
		log.Printf("ev2 is %+v", ev)
	default:
		return nil
	}
	return nil
}

func hasReviewMagicWord(s string) bool {
	return string([]rune(s)[:2]) == "r?"
}
