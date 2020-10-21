package usecase

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/appcache"
	"github.com/hayashiki/mentions/config"
	"github.com/hayashiki/mentions/event"
	"github.com/hayashiki/mentions/notifier"
	"github.com/hayashiki/mentions/repository"
	"github.com/slack-go/slack"
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
	notifyService notifier.Notifier
	taskRepo      repository.TaskRepository
	cache         appcache.InMemoryCache
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

	var postResp notifier.BotPostResp
	w.cache.Get(key, &postResp)

	log.Printf("commentID %v", postResp)

	task, err := w.taskRepo.GetByID(ev.Repository.ID)

	bot := slack.New(task.Slack.BotToken)
	n := notifier.NewSlackNotifier(bot)

	if err != nil {
		return fmt.Errorf("failed to get task %v", err)
	}

	payload := notifier.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	if comment, ok := w.notifyService.ConvertComment(payload, task.Users); ok {
		// キャッシュしなおしてもいいかも
		if _, err := n.UpdateSilently(postResp.Channel, postResp.Timestamp, comment); err != nil {
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

		if 1 ==0 {
			if err := w.notifyService.Notify(task.WebhookURL, comment); err != nil {
				return err
			}
		}

		bot := slack.New(task.Slack.BotToken)
		n := notifier.NewSlackNotifier(bot)

		resp, err := n.BotNotify(task.Slack.Channel, comment)

		if err != nil {
			return err
		}

		key := strconv.Itoa(int(ev.CommentID))
		log.Printf("commentID %d", ev.CommentID)

		w.cache.Add(key, resp, 10 * time.Minute)

		var postResp notifier.BotPostResp
		w.cache.Get(key, &postResp)

		log.Printf("postResp %v", postResp)
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

		bot := slack.New(task.Slack.BotToken)
		n := notifier.NewSlackNotifier(bot)

		resp, err := n.BotNotify(task.Slack.Channel, comment)
		if err != nil {
			return err
		}

		key := strconv.Itoa(int(ev.CommentID))
		log.Printf("commentID %d", ev.CommentID)

		w.cache.Add(key, resp, 10 * time.Minute)
	}

	return nil
}

func hasReviewMagicWord(s string) bool {
	return string([]rune(s)[:2]) == "r?"
}
