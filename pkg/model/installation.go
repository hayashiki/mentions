package model

import "time"

const InstallationKind = "installations"

type Installation struct {
	ID   int64  `json:"id" datastore:"id"`
	Name string `json:"name" datastore:"name"`
	// TODO AddedBy的なのいる？
	CreatedAt time.Time `datastore:"createdAt,noindex"`
	UpdatedAt time.Time `datastore:"updatedAt,noindex"`
}
