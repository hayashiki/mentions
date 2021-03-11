package model

import (
	"time"
)

const TeamKind = "teams"

type Team struct {
	ID          int64  `json:"id" datastore:"id"`
	SlackTeamID int64  `json:"slackTeamId" datastore:"slackTeamId"`
	Name        string `json:"name" datastore:"name"`
	Token       string `json:"token" datastore:"token"`
	Tasks       []Task
	// TODO: Authorized github organizations がほしい
	CreatedAt time.Time `datastore:"createdAt,noindex"`
	UpdatedAt time.Time `datastore:"updatedAt,noindex"`
}

func (Team) IsNode() {}
