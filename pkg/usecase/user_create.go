package usecase

import (
	"context"
	"fmt"
	"github.com/hayashiki/mentions/pkg/model"
	"github.com/hayashiki/mentions/pkg/repository"
	"github.com/hayashiki/mentions/pkg/slack"
	"log"
	"time"
)

type userCreate struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserCreate(teamRepo repository.TeamRepository, userRepo repository.UserRepository) *userCreate {
	return &userCreate{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

type UserCreateInput struct {
	TeamID   string
	UserID   string
	GithubID string
	SlackID  string
}

type UserCreateOutput struct {
	User *model.User
}

func (uc *userCreate) Do(ctx context.Context, input UserCreateInput) (*UserCreateOutput, error) {

	team, err := uc.teamRepo.Get(ctx, input.TeamID)
	if err != nil {
		return nil, fmt.Errorf("fail to get team, err=%v", err)
	}

	log.Printf("input %+v", input)

	// slack_user_authenticateと同じ
	slackSvc := slack.NewClient(slack.New(team.Token))
	slackUser, err := slackSvc.GetUser(input.SlackID)

	user := &model.User{
		ID:        input.SlackID,
		SlackID:   input.SlackID,
		GithubID:  model.GithubID(input.GithubID),
		Name:      slackUser.Name,
		Avatar:    slackUser.Avatar,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log.Printf("user %+v", user)

	err = uc.userRepo.Put(ctx, team.ID, user)
	if err != nil {
		return nil, err
	}

	return &UserCreateOutput{User: user}, nil
}
