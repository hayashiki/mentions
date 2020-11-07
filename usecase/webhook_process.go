package usecase

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/appcache"
	"github.com/hayashiki/mentions/config"
	"github.com/hayashiki/mentions/event"
	"github.com/hayashiki/mentions/model"
	"github.com/hayashiki/mentions/repository"
	"github.com/hayashiki/mentions/slack"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WebhookProcess interface {
	Do(r *http.Request) error
}

type webhookProcess struct {
	config        config.Environment
	githubService repository.Github
	taskRepo      repository.TaskRepository
	cache         appcache.InMemoryCache
}

func NewWebhookProcess(
	env config.Environment,
	gSvc repository.Github,
	taskRepo repository.TaskRepository,
) WebhookProcess {
	return &webhookProcess{
		config:        env,
		githubService: gSvc,
		taskRepo:      taskRepo,
		cache: appcache.NewInMemoryCache(10 * time.Minute),
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
		switch ev.GetAction() {
		case "created":
			return w.processIssueComment(ev)
		case "edited":
			return w.processEditIssueComment(ev)
		}
	case *github.PullRequestReviewCommentEvent:
		return w.processPullRequestComment(ev)
	default:
		return nil
	}
	return nil
}

func (w *webhookProcess) processEditIssueComment(ghEvent *github.IssueCommentEvent) error {
	log.Printf("Called.processEditIssueComment")

	ev := event.NewIssueComment(ghEvent)
	key := strconv.Itoa(int(ev.CommentID))

	var postResp slack.MessageResponse
	w.cache.Get(key, &postResp)

	task, err := w.taskRepo.GetByID(ev.Repository.ID)
	slackSvc := slack.NewClient(slack.New(task.Slack.BotToken))

	if err != nil {
		return fmt.Errorf("failed to get task %v", err)
	}

	payload := slack.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	if comment, ok := slackSvc.ConvertComment(payload, task.Users); ok {
		// キャッシュしなおしてもいいかも
		log.Printf("n debug")
		if _, err := slackSvc.UpdateMessage(postResp.Channel, postResp.Timestamp, comment); err != nil {
			log.Printf("n debug err %v", err)
			return err
		}
	}

	return nil
}


func (w *webhookProcess) processIssueComment(ghEvent *github.IssueCommentEvent) error {
	log.Printf("Called.processCreateIssueComment")
	ev := event.NewIssueComment(ghEvent)

	task, err := w.taskRepo.GetByID(ev.Repository.ID)
	if err != nil {
		return fmt.Errorf("failed to get task %v", err)
	}

	slackSvc := slack.NewClient(slack.New(task.Slack.BotToken))

	if hasReviewMagicWord(ev.Comment) {
		if err := w.editIssue(task, ev); err != nil {
			return fmt.Errorf("failed to edit github issue %v", err)
		}
	}

	payload := slack.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	if comment, ok := slackSvc.ConvertComment(payload, task.Users); ok {

		resp, err := slackSvc.PostMessage(task.Slack.Channel, comment)
		if err != nil {
			log.Printf("err is %v", err)
			return err
		}
		key := strconv.Itoa(int(ev.CommentID))
		w.cache.Add(key, resp, 10 * time.Minute)
		var postResp slack.MessageResponse
		w.cache.Get(key, &postResp)
	}
	return nil
}

func (w *webhookProcess) processPullRequestComment(ghEvent *github.PullRequestReviewCommentEvent) error {
	ev := event.NewPullRequestCommentEvent(ghEvent)

	task, err := w.taskRepo.GetByID(ev.Repository.ID)
	slackSvc := slack.NewClient(slack.New(task.Slack.BotToken))

	if err != nil {
		return err
	}

	payload := slack.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	if comment, ok := slackSvc.ConvertComment(payload, task.Users); ok {

		bot := slack.New(task.Slack.BotToken)
		n := slack.NewClient(bot)

		resp, err := n.PostMessage(task.Slack.Channel, comment)
		if err != nil {
			return err
		}

		key := strconv.Itoa(int(ev.CommentID))
		log.Printf("commentID %d", ev.CommentID)

		w.cache.Add(key, resp, 10 * time.Minute)
	}

	return nil
}

func (w *webhookProcess) editIssue(task *model.Task, ev *event.Event) error {
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
	return nil
}

func hasReviewMagicWord(s string) bool {
	return string([]rune(s)[:2]) == "r?"
}

