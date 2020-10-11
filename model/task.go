package model

const TaskKind = "tasks"

type Task struct {
	ID         int64 // github repositoryID
	Workspace  string
	Repo       Repo
	WebhookURL string
	Users      []User
}

type Repo struct {
	ID    string
	Owner string
	Name  string
}

func (t *Task) GetUserByGithubID(githubID string) (*User, bool) {

	found := false

	var user *User

	for _, u := range t.Users {
		if u.GithubID.String() == githubID {
			user = &u
			found = true
		}
	}

	return user, found
}
