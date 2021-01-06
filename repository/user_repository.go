package repository

import (
	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"

	"context"
	"fmt"
	"github.com/hayashiki/mentions/model"
)

//go:generate mockgen -source user_repository.go -destination mock_repo/user_repository.go
type UserRepository interface {
	List(team *model.Team, cursor string, limit int) ([]*model.User, string, error)
	Put(team *model.Team, user *model.User) error
	FindByGithubID(githubID string) (*model.User, error)
	FindBySlackID(githubID string) (*model.User, error)
}

type userRepository struct {
	dsClient *datastore.Client
}

func (r *userRepository) List(team *model.Team, cursor string, limit int) ([]*model.User, string, error) {
	q := datastore.NewQuery(model.UserKind).Ancestor(r.parentKey(team))
	if cursor != "" {
		dsCursor, err := datastore.DecodeCursor(cursor)
		if err != nil {

		}
		q = q.Start(dsCursor)
	}
	q = q.Limit(limit)

	var el []*model.User
	ctx := context.Background()
	it := r.dsClient.Run(ctx, q)
	for {
		var e model.User
		if _, err := it.Next(&e); err == iterator.Done {
			break
		} else if err != nil {
			fmt.Errorf("fail err=%v", err)
		}
		el = append(el, &e)
	}

	nextCursor, err := it.Cursor()
	if err != nil {
		return nil, "", err
	}

	return el, nextCursor.String(), nil
}

func (r *userRepository) key(id string, t *model.Team) *datastore.Key {
	return datastore.NameKey(model.UserKind, id, r.parentKey(t))
}

func (r *userRepository) parentKey(t *model.Team) *datastore.Key {
	return datastore.NameKey(model.TeamKind, t.ID, nil)
}

func NewUserRepository(client *datastore.Client) UserRepository {
	return &userRepository{dsClient: client}
}

type Reviewer struct {
	SlackID string
}

func (r *userRepository) Put(team *model.Team, user *model.User) error {
	ctx := context.Background()
	_, err := r.dsClient.Put(ctx, r.key(user.ID, team), user)
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
