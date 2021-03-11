package model

import "time"

const RepoKind = "repos"

// TODO: json, datastore
type Repo struct {
	ID        int64     `json:"id" datastore:"id"`
	Owner     string    `json:"owner" datastore:"owner"`
	Name      string    `json:"name" datastore:"name"`
	FullName  string    `json:"fullname" datastore:"fullname"`
	CreatedAt time.Time `datastore:"createdAt,noindex"`
	UpdatedAt time.Time `datastore:"updatedAt,noindex"`
}

func (Repo) IsNode() {}

