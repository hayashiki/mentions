package usecase

import (
	"context"
	"fmt"
	"github.com/hayashiki/go-pkg/slack/auth"
	"github.com/hayashiki/mentions/pkg/model"
	"github.com/hayashiki/mentions/pkg/repository"
	"log"
	"time"
)

type Auth interface {
	Do(ctx context.Context, i *AuthInput) error
}

type authUsecase struct {
	token    auth.Token
	teamRepo repository.TeamRepository
}

func NewTeamAuthProcess(token auth.Token, teamRepo repository.TeamRepository) Auth {
	return &authUsecase{
		token:    token,
		teamRepo: teamRepo,
	}
}

type AuthInput struct {
	Code string
}

func (uc *authUsecase) Do(ctx context.Context, i *AuthInput) error {
	authResp, err := uc.token.GetAccessToken(i.Code)

	log.Printf("authResp is %+v", authResp)
	if err != nil {
		return fmt.Errorf("fail to get token %v", err)
	}

	team := &model.Team{
		ID:        authResp.Team.ID,
		Name:      authResp.Team.Name,
		Token:     authResp.AccessToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return uc.teamRepo.Put(ctx, team)
}
