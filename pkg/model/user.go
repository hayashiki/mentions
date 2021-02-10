package model

import (
	"fmt"
	"time"
)

const UserKind = "users"

type User struct {
	ID           string    `json:"id" datastore:"id"`
	Name         string    `json:"name" datastore:"name"`
	Email        string    `json:"email" datastore:"email"`
	GoogleID     string    `json:"googleId" datastore:"googleId"`
	Avatar       string    `json:"avatar" datastore:"avatar"`
	Token        string    `json:"token" datastore:"token"`
	SlackID      string    `json:"slackId" datastore:"slackId"`
	SlackIsOwner string    `json:"slackIsOwner" datastore:"slackIsOwner"`
	SlackIsAdmin string    `json:"slackIsAdmin" datastore:"slackIsAdmin"`
	GithubID     GithubID  `json:"githubId" datastore:"githubId"`
	Reviewers    Reviewers `json:"reviewers" datastore:"reviewers"`
	GroupID      int64     `json:"groupId" datastore:"groupId"`
	TeamID       string    `json:"teamId" datastore:"teamId"`
	CreatedAt    time.Time `datastore:"createdAt,noindex"`
	UpdatedAt    time.Time `datastore:"updatedAt,noindex"`
}

type UserID int64

type GithubID string

func (g GithubID) String() string {
	return string(g)
}

func (g GithubID) WithAt() string {
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

func (u User) SetGithubID(ghID string) {
	u.GithubID = GithubID(ghID)
}

func (u User) SetSlackID(slackID string) {
	u.SlackID = slackID
}
