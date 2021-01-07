package model

import "log"

const TaskKind = "tasks"

type Task struct {
	ID             int64 // github repositoryID
	Repo           Repo
	Channel        string `json:"channel" datastore:"channel"`
	Users          []*User
	UserIDs        []string
	TeamID         string `json:"teamId" datastore:"teamId"`
	Team           *Team
	InstallationID int64 `json:"id" datastore:"installationId"`
}

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
