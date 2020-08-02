package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/account"
	"github.com/hayashiki/mentions/config"
	"github.com/hayashiki/mentions/gh"
	"github.com/hayashiki/mentions/notifier/interface/notifier"
	"log"
	"net/http"
)

type WebhookHandler struct {
	env      config.Environment
	Verifier gh.Verifier
	Notifier notifier.Notifier
	List     account.List
}

func NewWebhookHandler(
	verifier gh.Verifier,
	notifier notifier.Notifier,
	env config.Environment,
	list account.List) WebhookHandler {
	return WebhookHandler{
		env:      env,
		Notifier: notifier,
		Verifier: verifier,
		List: list,
	}
}

func (h *WebhookHandler) PostWebhook(c *gin.Context) {
	secretKey := []byte(h.env.GithubWebhookSecret)
	payload, err := h.Verifier.Verify(c.Request, secretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ghEvent, err := github.ParseWebHook(github.WebHookType(c.Request), payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	event, ok := h.ConvertGithubEvent(ghEvent)
	if !ok {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if event.Action != "created" {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if string([]rune(event.Comment)[:2]) == "r?" {
		event.CreateReviewers(h.List.Reviewers)
	}

	//TODO handle not found
	log.Printf("repo full %v", event.Repository.FullName)
	url := h.List.Repos[event.Repository.FullName]

	log.Printf("url %v", url)

	err = h.Notifier.Notify(url, event.GenerateMessage())
	if err != nil {
		log.Printf("call error: %v", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
	return
}

func (h WebhookHandler) ConvertGithubEvent(original interface{}) (*gh.Event, bool) {
	switch event := original.(type) {
	case *github.IssueCommentEvent:
		return gh.ConvertIssueCommentEvent(event), true
	case *github.PullRequestReviewCommentEvent:
		return gh.ConvertPullRequestCommentEvent(event), true
	default:
		return nil, false
	}

	//switch event := event.(type) {
	//case *github.IssueCommentEvent:
	//	log.Printf("IssueComment event %v", *event.Action)
	//	h.HandleGithubCommentEvent(c.Writer, event)
	//default:
	//	log.Printf("no event %v", event)
	//	return
	//}
}
