package model

import "time"

const GroupKind = "groups"

type Group struct {
	ID        int64     `json:"id" datastore:"id"`
	Name      string    `json:"name" datastore:"name"`
	UserIDs   []string  `json:"userIds" datastore:"userIds"`
	CreatedAt time.Time `datastore:"createdAt,noindex"`
	UpdatedAt time.Time `datastore:"updatedAt,noindex"`
}
