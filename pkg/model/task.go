package model

import "log"

const TaskKind = "tasks"

type Task struct {
	ID             int64 `json:"id" datastore:"id"`
	Repo           Repo
	RepoID         int64  `json:"repoId" datastore:"repoId"`
	Channel        string `json:"channel" datastore:"channel"`
	Users          []*User
	UserIDs        []string
	TeamID         string `json:"teamId" datastore:"teamId"`
	Team           *Team
	InstallationID int64 `json:"installationId" datastore:"installationId"`
}

func (Task) IsNode() {}

func (t *Task) GetUserByGithubID(githubID string) (*User, bool) {

	found := false
	var user *User
	for _, u := range t.Users {
		log.Printf("u is %v", u)
		log.Printf(githubID)
		if u.GithubID.String() == githubID {
			user = u
			found = true
			break
		}
	}

	return user, found
}

func (t *Task) SetUsers(users []*User) {
	t.Users = users
}
