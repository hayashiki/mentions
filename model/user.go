package model

import (
	"fmt"
	"time"
)

const UserKind = "users"

type User struct {
	ID        string
	Workspace string
	SlackID   string `json:"slack_id" datastore:"slack_id"`
	GithubID  GithubID `json:"github_id" datastore:"github_id"`
	Reviewers Reviewers
	CreatedAt time.Time `datastore:"created_at,noindex"`
	UpdatedAt time.Time `datastore:"updated_at,noindex"`
}

type UserID int64

type GithubID string

func (g GithubID) String() string{
	return string(g)
}

func (g GithubID) WithAt() string{
	return fmt.Sprintf("@%s", g.String())
}

type Reviewers []GithubID

func (rs Reviewers) String() []string {
	var reviewers []string
	for _, r := range rs {
		reviewers = append(reviewers, r.String())
	}
	return reviewers
}

type Reviewer struct {
	SlackID string
}

func (u User) GithubWithAt() string {
	return u.GithubID.WithAt()
}

func (u User) ReviewersWithAt() []string {
	var reviewers []string

	for _, r := range u.Reviewers {
		reviewers = append(reviewers, r.WithAt())
	}

	return reviewers
}

func (u User) SlackWithBracketAt() string {
	return fmt.Sprintf("<@%s>", u.SlackID)
}
