package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/hayashiki/mentions/model"
	"time"
)

//go:generate mockgen -source user_repository.go -destination mock_repo/user_repository.go
type UserRepository interface {
	//List()
	//Get()
	Put(user *model.User) error
	FindByGithubID(githubID string) (*model.User, error)
	FindBySlackID(githubID string) (*model.User, error)
}

type userRepository struct {
	dsClient *datastore.Client
}

func NewUserRepository(client *datastore.Client) UserRepository {
	return &userRepository{dsClient: client}
}

type user struct {
	ID        *datastore.Key `json:"id" datastore:"__key__"`
	SlackID   string         `json:"slack_id" datastore:"slack_id"`
	GithubID  string         `json:"slack_id" datastore:"github_id"`
	Reviewers []Reviewer
	CreatedAt time.Time `datastore:"RegisteredAt,noindex"`
}

type Reviewer struct {
	SlackID string
}

func (r *userRepository) Put(user *model.User) error {
	ctx := context.Background()
	k := datastore.NameKey(model.UserKind, user.ID, nil)
	_, err := r.dsClient.Put(ctx, k, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByGithubID(githubID string) (*model.User, error) {
	q := datastore.NewQuery(model.UserKind).KeysOnly().Filter("GithubID =", githubID).Limit(1)

	keys, err := r.dsClient.GetAll(context.Background(), q, nil)

	if err != nil {
		return nil, fmt.Errorf("not found user keys %w", err)
	}

	var user model.User

	if err := r.dsClient.Get(context.Background(), keys[0], &user); err != nil {
		return nil, fmt.Errorf("not found user %w", err)
	}

	user.ID = keys[0].Name

	return &user, nil
	//if err = Transaction.Get(keys[0], &user); err != nil {
	//
	//}
}

func (r *userRepository) FindBySlackID(slackID string) (*model.User, error) {
	q := datastore.NewQuery(model.UserKind).KeysOnly().Filter("GithubID =", slackID).Limit(1)

	keys, err := r.dsClient.GetAll(context.Background(), q, nil)

	if err != nil {
		return nil, fmt.Errorf("not found user %w", err)
	}

	var user model.User

	if err := r.dsClient.Get(context.Background(), keys[0], &user); err != nil {
		return nil, fmt.Errorf("not found user %w", err)
	}

	user.ID = keys[0].Name

	return &user, nil
}

//func (r *userRepository) Get(ctx context.Context, id string) (*user.User, error) {
//	dst := &model.User{}
//	err := r.dsClient.Get(ctx, repo.key(id), dst)
//	if err != nil {
//		return nil, err
//	}
//	return dst, nil
//}
