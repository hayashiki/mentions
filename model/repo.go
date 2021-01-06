package model

const RepoKind = "repos"

// TODO: json, datastore
type Repo struct {
	ID       int64
	Owner    string
	Name     string
	FullName string
}
