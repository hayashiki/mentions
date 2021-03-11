package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hayashiki/mentions/pkg/slack"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (a *App) currentUser(ctx context.Context) *User {
	log.Printf("UserIDKey %v", UserIDKey)

	if u, ok := ctx.Value(UserIDKey).(*User); ok {
		log.Printf("currentUser %v", u)
		return u
	}
	return nil
}

func (a *App) ListSlackAPIUser(w http.ResponseWriter, r *http.Request) error {

	u := a.currentUser(r.Context())
	team, err := a.teamRepo.GetBySlackTeamID(r.Context(), u.TeamID)
	if err != nil {
		return err
	}
	slackSvc := slack.NewClient(slack.New(team.Token))

	users, err := slackSvc.GetUsers()
	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		return fmt.Errorf("fail to render users, err=%v", err)
	}
	//jsonResponse(w, http.StatusOK, users)
	return nil
}

// TODO: model.userで特定カラムを非公開にすればよいだけ？
type listResp struct {
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	SlackID  string `json:"slackId"`
	GithubID string `json:"githubId"`
}

func (a *App) ListUser(w http.ResponseWriter, r *http.Request) error {
	u := a.currentUser(r.Context())

	team, err := a.teamRepo.GetBySlackTeamID(r.Context(), u.TeamID)
	if err != nil {
		return fmt.Errorf("fail to get team, err=%v", err)
	}

	//TODO: cursorは一旦考えない
	users, _, err := a.userRepo.List(r.Context(), team.ID, "", 100)
	if err != nil {
		return fmt.Errorf("fail to get users, err=%v", err)
	}

	var listResps []*listResp

	for _, u := range users {
		uu := &listResp{
			//ID:       u.SlackID,
			SlackID:  u.SlackID,
			Name:     u.Name,
			Avatar:   u.Avatar,
			GithubID: u.GithubID.String(),
		}
		listResps = append(listResps, uu)
	}

	if err := json.NewEncoder(w).Encode(listResps); err != nil {
		return fmt.Errorf("fail to render users, err=%v", err)
	}

	return nil
}

func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")
	req := struct {
		SlackID  string `json:"slackId"`
		GithubID string `json:"githubId"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	defer r.Body.Close()

	u := a.currentUser(r.Context())

	team, err := a.teamRepo.GetBySlackTeamID(r.Context(), u.TeamID)

	user, err := a.userRepo.FindBySlackID(r.Context(), id)
	if err != nil {
		return err
	}
	user.SetSlackID(req.SlackID)
	user.SetGithubID(req.GithubID)

	a.userRepo.Put(r.Context(), team.ID, user)

	return nil
}
