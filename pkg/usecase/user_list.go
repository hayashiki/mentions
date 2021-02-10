package usecase

import (
	"context"
	"fmt"
	"github.com/hayashiki/mentions/pkg/repository"
)

type userList struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserList(teamRepo repository.TeamRepository, userRepo repository.UserRepository) *userCreate {
	return &userCreate{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

type UserListInput struct {
	TeamID   string
	UserID   string
	GithubID string
	SlackID  string
}

func (uc *userList) Do(ctx context.Context, input UserCreateInput) error {

	team, err := uc.teamRepo.Get(ctx, input.TeamID)
	if err != nil {
		return fmt.Errorf("fail to get team, err=%v", err)
	}
	uc.userRepo.List(ctx, team.ID, "", 100)

	return err
}
