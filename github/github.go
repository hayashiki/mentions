package github

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
)

type Github interface {
	Verify(r *http.Request, secret []byte) ([]byte, error)
	ParseWebHook(r *http.Request, payload []byte) (interface{}, error)
	CreateReviewers(payload *CreateReviewersPayload) (*github.PullRequest, *github.Response, error)
	EditIssueComment(payload *EditIssueCommentPayload) (*github.IssueComment, *github.Response, error)
}

type client struct {
	ghClient *github.Client
}

func NewClient(ghClient *github.Client) Github {
	return &client{
		ghClient: ghClient,
	}
}

func GetClient(token string) *github.Client {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

type CreateReviewersPayload struct {
	Owner       string
	Name        string
	IssueNumber int
	Reviewers   []string
}

type EditIssueCommentPayload struct {
	Owner     string
	Name      string
	CommentID int64
	Comment   string
}

func (c *client) CreateReviewers(payload *CreateReviewersPayload) (*github.PullRequest, *github.Response, error) {
	ctx := context.Background()
	return c.ghClient.PullRequests.RequestReviewers(ctx, payload.Owner, payload.Name, payload.IssueNumber, github.ReviewersRequest{Reviewers: payload.Reviewers})
}

func (c *client) EditIssueComment(payload *EditIssueCommentPayload) (*github.IssueComment, *github.Response, error) {
	ctx := context.Background()

	comment := new(github.IssueComment)
	comment.Body = github.String(payload.Comment)

	return c.ghClient.Issues.EditComment(ctx, payload.Owner, payload.Name, payload.CommentID, comment)
}

func (c *client) Verify(r *http.Request, secret []byte) ([]byte, error) {
	payload, err := github.ValidatePayload(r, secret)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *client) ParseWebHook(r *http.Request, payload []byte) (interface{}, error) {
	return github.ParseWebHook(github.WebHookType(r), payload)
}

////go:generate mockgen -source verifier.go -destination mocks/verifier.go
//
////gomock使う程でもないので、VerifierをみたすMockGithubVerifierを用意
//type MockGithubVerifier struct {
//	Valid bool
//}
//
//func (v *MockGithubVerifier) Verify(r *http.Request, secret []byte) ([]byte, error) {
//	if v.Valid {
//		return nil, nil
//	}
//	return nil, errors.New("Invalid Secrets")
//}
