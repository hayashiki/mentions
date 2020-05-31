package gh

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func initGithubClient(ctx context.Context, secretGithub string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: secretGithub,
		},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return client
}
