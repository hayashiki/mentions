package gh

import (
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/account"
	"github.com/hayashiki/mentions/config"
	notifier2 "github.com/hayashiki/mentions/notifier/interface/notifier"
	"log"
	"net/http"
)

type WebhookHandler struct {
	env      config.Environment
	Verifier Verifier
	Notifier notifier2.Notifier
	List     account.List
}

func NewWebhookHandler(
	verifier Verifier,
	notifier notifier2.Notifier,
	env config.Environment,
	list account.List) WebhookHandler {
	return WebhookHandler{
		env:      env,
		Notifier: notifier,
		Verifier: verifier,
		List: list,
	}
}

func (h WebhookHandler) PostWebhook(c *gin.Context) {
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

	//url, err := FindWebhookURL(event, h.List)

	err = h.Notifier.Notify(event, "url")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
	return
}

func FindWebhookURL(event *Event, list account.List) (string, error){
	url := list.Repos[event.Repository]
	log.Printf("url is %v", url)
	return "", nil
}

func (h WebhookHandler) ConvertGithubEvent(original interface{}) (*Event, bool) {
	switch event := original.(type) {
	case *github.IssueCommentEvent:
		return convertIssueCommentEvent(event), true
	case github.PullRequestEvent:
		return nil, false
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

// Event is the internal structure used for an event
type Event struct {
	Action      string
	IssueNumber int
	Repository  string
	User        string
	Title       string
	Comment     string
	HTMLURL     string
	Type        EventType

	OriginalEvent interface{}
}

type EventType int

const (
	IssueCommentEvent EventType = iota
	PullRequestCommentEvent
	PingEvent
)

func convertIssueCommentEvent(original *github.IssueCommentEvent) *Event {
	return &Event{
		IssueNumber: original.Issue.GetNumber(),
		Action:      original.GetAction(),
		Repository:  original.Repo.GetName(),
		User:        original.Comment.User.GetLogin(),
		Comment:     original.Comment.GetBody(),
		Type:        IssueCommentEvent,
	}
}
