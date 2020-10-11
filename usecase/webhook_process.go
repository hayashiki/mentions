package usecase

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/config"
	"github.com/hayashiki/mentions/event_processor"
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

	switch event := ghEvent.(type) {
	case *github.IssueCommentEvent:
		return w.processIssueComment(event)
	case *github.PullRequestReviewCommentEvent:
		return w.processPullRequestComment(event)
	default:
		return nil
	}
}

func (w *webhookProcess) processIssueComment(ghEvent *github.IssueCommentEvent) error {
	log.Printf("Called.processIssueComment")
	event := event_processor.NewIssueComment(ghEvent)

	log.Printf("RepositoryID is %v", event.Repository.ID)
	task, err := w.taskRepo.GetByID(event.Repository.ID)

	if err != nil {
		return fmt.Errorf("failed to get task %v", err)
	}
	log.Printf("task is %v", task)

	if hasReviewMagicWord(event.Comment) {

		user, ok := task.GetUserByGithubID(event.IssueOwner)
		if !ok {
			return fmt.Errorf("github user not found user %s", event.IssueOwner)
		}

		payload := &repository.CreateReviewersPayload{
			Owner:       event.Repository.Owner,
			Name:        event.Repository.Name,
			IssueNumber: event.IssueNumber,
			Reviewers:   user.Reviewers.String(),
		}

		_, resp, err := w.githubService.CreateReviewers(payload)
		if err != nil {
			return fmt.Errorf("failed to create reviewer resp %v, err=%v", resp, err)
		}

		comment := strings.Join(user.ReviewersWithAt(), " ") + " レビューお願いします😀"
		event.Comment = comment

		commentPayload := &repository.EditIssueCommentPayload{
			Owner:     event.Repository.Owner,
			Name:      event.Repository.Name,
			CommentID: event.CommentID,
			Comment:   event.Comment,
		}

		_, resp, err = w.githubService.EditIssueComment(commentPayload)
		if err != nil {
			return fmt.Errorf("failed to edit issue resp %v, err=%v", resp, err)
		}
	}

	payload := notifier.ConvertPayload{
		Comment:  event.Comment,
		RepoName: event.Repository.FullName,
		HTMLURL:  event.HTMLURL,
		Title:    event.Title,
		User:     event.User,
	}

	if comment, ok := w.notifyService.ConvertComment(payload, task.Users); ok {
		w.notifyService.Notify(task.WebhookURL, comment)
	}

	return nil
}

func (w *webhookProcess) processPullRequestComment(ghEvent *github.PullRequestReviewCommentEvent) error {
	event := event_processor.NewPullRequestCommentEvent(ghEvent)

	task, err := w.taskRepo.GetByID(event.Repository.ID)

	if err != nil {
		return err
	}

	payload := notifier.ConvertPayload{
		Comment:  event.Comment,
		RepoName: event.Repository.FullName,
		HTMLURL:  event.HTMLURL,
		Title:    event.Title,
		User:     event.User,
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
