package usecase

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/config"
	"github.com/hayashiki/mentions/event"
	"github.com/hayashiki/mentions/notifier"
	"github.com/hayashiki/mentions/repository"
	"log"
	"net/http"
	"strings"
)

type WebhookProcess interface {
	Do(r *http.Request) error
}

type webhookProcess struct {
	config        config.Environment
	githubService repository.Github
	notifyService notifier.Notifier
	taskRepo      repository.TaskRepository
}

func NewWebhookProcess(
	env config.Environment,
	gSvc repository.Github,
	nSvc notifier.Notifier,
	taskRepo repository.TaskRepository,
) WebhookProcess {
	return &webhookProcess{
		config:        env,
		githubService: gSvc,
		notifyService: nSvc,
		taskRepo:      taskRepo,
	}
}

func (w webhookProcess) Do(r *http.Request) error {
	// TODO: get from secret manager
	payload, err := w.githubService.Verify(r, []byte(w.config.GithubWebhookSecret))

	if err != nil {
		return err
	}

	ghEvent, err := w.githubService.ParseWebHook(r, payload)

	if err != nil {
		return err
	}

	switch ev := ghEvent.(type) {
	case *github.IssueCommentEvent:
		return w.processIssueComment(ev)
	case *github.PullRequestReviewCommentEvent:
		return w.processPullRequestComment(ev)
	default:
		return nil
	}
}

func (w *webhookProcess) processIssueComment(ghEvent *github.IssueCommentEvent) error {
	log.Printf("Called.processIssueComment")
	ev := event.NewIssueComment(ghEvent)

	task, err := w.taskRepo.GetByID(ev.Repository.ID)

	if err != nil {
		return fmt.Errorf("failed to get task %v", err)
	}
	log.Printf("task is %v", task)

	if hasReviewMagicWord(ev.Comment) {

		user, ok := task.GetUserByGithubID(ev.IssueOwner)
		if !ok {
			return fmt.Errorf("github user not found user %s", ev.IssueOwner)
		}

		payload := &repository.CreateReviewersPayload{
			Owner:       ev.Repository.Owner,
			Name:        ev.Repository.Name,
			IssueNumber: ev.IssueNumber,
			Reviewers:   user.Reviewers.String(),
		}

		_, resp, err := w.githubService.CreateReviewers(payload)
		if err != nil {
			return fmt.Errorf("failed to create reviewer resp %v, err=%v", resp, err)
		}

		comment := strings.Join(user.ReviewersWithAt(), " ") + " レビューお願いします😀"
		ev.Comment = comment

		commentPayload := &repository.EditIssueCommentPayload{
			Owner:     ev.Repository.Owner,
			Name:      ev.Repository.Name,
			CommentID: ev.CommentID,
			Comment:   ev.Comment,
		}

		_, resp, err = w.githubService.EditIssueComment(commentPayload)
		if err != nil {
			return fmt.Errorf("failed to edit issue resp %v, err=%v", resp, err)
		}
	}

	payload := notifier.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	if comment, ok := w.notifyService.ConvertComment(payload, task.Users); ok {
		if err := w.notifyService.Notify(task.WebhookURL, comment); err != nil {
			return err
		}
	}

	return nil
}

func (w *webhookProcess) processPullRequestComment(ghEvent *github.PullRequestReviewCommentEvent) error {
	ev := event.NewPullRequestCommentEvent(ghEvent)

	task, err := w.taskRepo.GetByID(ev.Repository.ID)

	if err != nil {
		return err
	}

	payload := notifier.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	if comment, ok := w.notifyService.ConvertComment(payload, task.Users); ok {
		if err := w.notifyService.Notify(task.WebhookURL, comment); err != nil {
			return fmt.Errorf("failed to send to slack err=%v", err)
		}
	}

	return nil
}

func hasReviewMagicWord(s string) bool {
	return string([]rune(s)[:2]) == "r?"
}
