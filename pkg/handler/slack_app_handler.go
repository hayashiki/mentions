package handler

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hayashiki/go-pkg/slack/auth"
	"github.com/hayashiki/mentions/pkg/usecase"
	"net/http"
	"time"
)

func (a *App) slackTeamAuthorize(w http.ResponseWriter, r *http.Request) error {
	state := uuid.New().String()
	authz := auth.NewAuth(a.config.SlackClientID, a.config.SlackRedirectURL, state, TeamScopes)

	http.SetCookie(w, &http.Cookie{
		Name:    "state",
		Value:   state,
		Expires: time.Now().Add(300 * time.Second),
		Secure:  true,
	})

	authz.Redirect(w, r)

	return nil
}

// slackTeamCallback
func (a *App) slackTeamCallback(w http.ResponseWriter, r *http.Request) error {
	resp := auth.ParseRequest(r)
	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return fmt.Errorf("invalid cookie, err=%v", err)
	}
	if resp.State != state.Value {
		return fmt.Errorf("invalid state, err=%v", err)
	}
	token := auth.NewToken(a.config.SlackClientID, a.config.SlackSecretID, a.config.SlackRedirectURL)
	authz := usecase.NewTeamAuthProcess(token, a.teamRepo)
	input := &usecase.AuthInput{Code: resp.Code}
	if err := authz.Do(r.Context(), input); err != nil {
		return fmt.Errorf("invalid state, err=%v", err)
	}
	// TODO: リダイレクトする
	fmt.Fprintf(w, "team is: %s", "ok")
	return nil
}

func (a *App) slackUserAuthorize(w http.ResponseWriter, r *http.Request) error {
	state := uuid.New().String()
	authz := auth.NewUserAuth(a.config.SlackClientID, a.config.SlackUserRedirectURL, state, UserScopes)

	http.SetCookie(w, &http.Cookie{
		Name:    "state",
		Value:   state,
		Expires: time.Now().Add(300 * time.Second),
		Secure:  true,
	})

	authz.Redirect(w, r)

	return nil
}

func (a *App) slackUserCallback(w http.ResponseWriter, r *http.Request) error {
	resp := auth.ParseRequest(r)

	state, err := r.Cookie("state")
	if err != nil {
		return fmt.Errorf("invalid cookie, err=%v", err)
	}
	if resp.State != state.Value {
		return fmt.Errorf("invalid state, err=%v", err)

	}

	token := auth.NewToken(a.config.SlackClientID, a.config.SlackSecretID, a.config.SlackUserRedirectURL)
	userAuthenticateProcess := usecase.NewUserAuthenticateProcess(token, a.userRepo, a.teamRepo)
	input := usecase.UserAuthInput{Code: resp.Code}
	res, err := userAuthenticateProcess.Do(r.Context(), input)
	if err != nil {
		return fmt.Errorf("err=%v", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: res,
		Path:  "/",
	})

	// TODO
	urlStr := "/"
	http.Redirect(w, r, urlStr, 302)
	return nil
}
