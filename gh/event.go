package gh

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/account"
	"log"
)

func FindWebhookURL(event *Event, list account.List) (string, error){
	url := list.Repos[event.Repository]
	log.Printf("url is %v", url)
	return "", nil
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

func ConvertIssueCommentEvent(original *github.IssueCommentEvent) *Event {
	return &Event{
		IssueNumber: original.Issue.GetNumber(),
		Title:       original.Issue.GetTitle(),
		Action:      original.GetAction(),
		Repository:  original.Repo.GetName(),
		User:        original.Comment.User.GetLogin(),
		HTMLURL:     original.Issue.GetHTMLURL(),
		Comment:     original.Comment.GetBody(),
		Type:        IssueCommentEvent,
	}
}

func (event *Event) GenerateMessage() string {
	var text string
	text = fmt.Sprintf("%v *%v <%v|%v> * by: %v\n", text, event.Repository, event.HTMLURL, event.Title, event.User)
	text = fmt.Sprintf("%v\n%v\n", text, event.Comment)
	return text
}
