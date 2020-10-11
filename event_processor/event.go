package event_processor

import (
	"github.com/google/go-github/github"
	"log"
)

// Event is the internal structure used for an event
type Event struct {
	Action      string
	IssueNumber int
	IssueOwner  string
	CommentID   int64
	Repository  Repository
	User        string
	Title       string
	Comment     string
	HTMLURL     string
	Type        EventType

	OriginalEvent interface{}
}

// Repository is
type Repository struct {
	ID       int64
	Owner    string
	Name     string
	FullName string
}

type EventType int

const (
	IssueCommentEvent EventType = iota
	PullRequestCommentEvent
	PingEvent
)

func NewIssueComment(original *github.IssueCommentEvent) *Event {

	return &Event{
		IssueNumber: original.Issue.GetNumber(),
		IssueOwner:  original.Issue.GetUser().GetLogin(),
		Title:       original.Issue.GetTitle(),
		Action:      original.GetAction(),
		Repository: Repository{
			ID:       original.Repo.GetID(),
			Owner:    original.Repo.GetOwner().GetLogin(),
			Name:     original.Repo.GetName(),
			FullName: original.Repo.GetFullName(),
		},
		CommentID: original.Comment.GetID(),
		User:      original.Comment.User.GetLogin(),
		HTMLURL:   original.Issue.GetHTMLURL(),
		Comment:   original.Comment.GetBody(),
		Type:      IssueCommentEvent,
	}
}

func NewPullRequestCommentEvent(original *github.PullRequestReviewCommentEvent) *Event {

	log.Printf("Event PR debug %s", original.Repo.GetOwner().GetLogin())

	return &Event{
		IssueNumber: original.PullRequest.GetNumber(),
		IssueOwner:  original.PullRequest.GetUser().GetLogin(),
		Title:       original.PullRequest.GetTitle(),
		Action:      original.GetAction(),
		Repository: Repository{
			ID:       original.Repo.GetID(),
			Owner:    original.Repo.GetOwner().GetLogin(),
			Name:     original.Repo.GetName(),
			FullName: original.Repo.GetFullName(),
		},
		CommentID: original.Comment.GetID(),
		User:      original.Comment.User.GetLogin(),
		HTMLURL:   original.PullRequest.GetHTMLURL(),
		Comment:   original.Comment.GetBody(),
		Type:      PullRequestCommentEvent,
	}
}
