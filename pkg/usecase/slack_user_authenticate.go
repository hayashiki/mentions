package usecase

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hayashiki/go-pkg/slack/auth"
	"github.com/hayashiki/mentions/pkg/model"
	"github.com/hayashiki/mentions/pkg/repository"
	"github.com/hayashiki/mentions/pkg/slack"
	"log"
	"time"
)

var CookieNameToken = "token"

type UserAuth interface {
	Do(ctx context.Context, i UserAuthInput) (token string, err error)
}

type userAuth struct {
	token    auth.Token
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserAuthenticateProcess(token auth.Token, userRepo repository.UserRepository, teamRepo repository.TeamRepository) UserAuth {
	return &userAuth{
		token:    token,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

type UserAuthInput struct {
	Code string
}

func (uc *userAuth) Do(ctx context.Context, i UserAuthInput) (token string, err error) {
	authResp, err := uc.token.GetAccessToken(i.Code)

	log.Printf("authResp is %+v", authResp)
	if err != nil {
		return "", fmt.Errorf("fail to get token %v", err)
	}

	log.Printf("authResp.Team.ID is %+v", authResp.Team.ID)
	team, err := uc.teamRepo.Get(ctx, authResp.Team.ID)
	if err != nil {
		log.Printf("err %v", err)
		//return err
	}
	log.Printf("team %v", team)

	slackSvc := slack.NewClient(slack.New(team.Token))
	slackUser, err := slackSvc.GetUser(authResp.AuthedUser.ID)
	if err != nil {
		return "", err
	}

	user := &model.User{
		Token:     authResp.AuthedUser.AccessToken,
		ID:        authResp.AuthedUser.ID,
		SlackID:   authResp.AuthedUser.ID,
		Name:      slackUser.Name,
		Avatar:    slackUser.Avatar,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// TODO: ErrNoSuchEntity設定する。もどる可能性があるので。
	err = uc.userRepo.Put(ctx, team, user)
	if err != nil {
		return "", err
	}

	// HS512でよいか確認
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), jwt.MapClaims{
		"custom.id":        user.SlackID,
		"custom.namespace": team.ID,
		"exp":              time.Now().Add(time.Hour * 24).Unix(),
	})

	log.Printf("jwt stamp %v %v", user.SlackID, team.ID)

	// TODO: hogeではダメよね
	signedToken, err := jwtToken.SignedString([]byte("hoge"))

	log.Printf("signedToken %v", signedToken)
	log.Printf("err %v", err)

	return signedToken, err
}
