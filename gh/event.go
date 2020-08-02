package gh

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/account"
	"log"
	"os"
	"strings"
)

//func FindWebhookURL(event *Event, list account.List) (string, error){
//	url := list.Repos[event.Repository.FullName()]
//	log.Printf("url is %v", url)
//	return "", nil
//}

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
	Owner string
	Name  string
	FullName string
}

//func(r Repository) FullName() string{
//	return r.Owner + "/" + r.Name
//}

type EventType int

const (
	IssueCommentEvent EventType = iota
	PullRequestCommentEvent
	PingEvent
)

func ConvertIssueCommentEvent(original *github.IssueCommentEvent) *Event {

	log.Printf("Event Issue debug %s", original.Issue.GetUser().GetName())
	log.Printf("Event Issue debug %s", original.Issue.GetUser().GetLogin())
	return &Event{
		IssueNumber: original.Issue.GetNumber(),
		IssueOwner: original.Issue.GetUser().GetLogin(),
		Title:       original.Issue.GetTitle(),
		Action:      original.GetAction(),
		Repository: Repository{
			Owner: original.Repo.GetOwner().GetLogin(),
			Name: original.Repo.GetName(),
			FullName: original.Repo.GetFullName(),
		},
		CommentID: original.Comment.GetID(),
		User:        original.Comment.User.GetLogin(),
		HTMLURL:     original.Issue.GetHTMLURL(),
		Comment:     original.Comment.GetBody(),
		Type:        IssueCommentEvent,
	}
}

func ConvertPullRequestCommentEvent(original *github.PullRequestReviewCommentEvent) *Event {

	log.Printf("Event PR debug %s", original.Repo.GetOwner().GetLogin())

	return &Event{
		IssueNumber: original.PullRequest.GetNumber(),
		IssueOwner: original.PullRequest.GetUser().GetLogin(),
		Title:       original.PullRequest.GetTitle(),
		Action:      original.GetAction(),
		Repository: Repository{
			Owner: original.Repo.GetOwner().GetLogin(),
			Name: original.Repo.GetName(),
			FullName: original.Repo.GetFullName(),
		},
		CommentID: original.Comment.GetID(),
		User:        original.Comment.User.GetLogin(),
		HTMLURL:     original.PullRequest.GetHTMLURL(),
		Comment:     original.Comment.GetBody(),
		Type:        PullRequestCommentEvent,
	}
}

func (event *Event) GenerateMessage() string {
	var text string
	text = fmt.Sprintf("%v *%v <%v|%v> * by: %v\n", text, event.Repository.FullName, event.HTMLURL, event.Title, event.User)
	text = fmt.Sprintf("%v\n%v\n", text, event.Comment)
	return text
}

func (event *Event) CreateReviewers(reviewers account.Reviewers) (error){

	ctx := context.Background()
	var secretGithub = os.Getenv("GITHUB_SECRET_TOKEN")
	client := initGithubClient(ctx, secretGithub)

	issueSvc := client.Issues

	log.Printf("User is %s", event)
	log.Printf("User is %s", fmt.Sprintf("@", event.IssueOwner))
	revis := reviewers[fmt.Sprintf("@%s", event.IssueOwner)]

	log.Printf("revis is %v", revis)

	comment := new(github.IssueComment)
	comment.Body = github.String(strings.Join(revis, " ") + " レビューお願いします😀")

	var gReviewers []string
	for _, reviewer := range revis {
		r := strings.Replace(reviewer, "@", "", 1)
		if r != event.User {
			gReviewers = append(gReviewers, r)
		}
	}
	Reviewers := github.ReviewersRequest{Reviewers: gReviewers}

	resp, hogee, err := issueSvc.EditComment(ctx, event.Repository.Owner, event.Repository.Name, event.CommentID, comment)
	log.Printf("resp is %v", resp)
	log.Printf("hogee is %v", hogee)
	log.Printf("resp is %v", err)

	aa, cc, dd := client.PullRequests.RequestReviewers(ctx, event.Repository.Owner, event.Repository.Name, event.IssueNumber, Reviewers)
	log.Printf("resp is %v", aa)
	log.Printf("hogee is %v", cc)
	log.Printf("resp is %v", dd)

	return nil
}

//func formatReviewer()
