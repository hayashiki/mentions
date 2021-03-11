package repository

import (
	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"

	"context"
	"fmt"
	"github.com/hayashiki/mentions/pkg/model"
)

//go:generate mockgen -source user_repository.go -destination mock_repo/user_repository.go
type UserRepository interface {
	List(ctx context.Context, teamID int64, cursor string, limit int) ([]*model.User, string, error)
	Put(ctx context.Context, teamID int64, user *model.User) error
	Get(ctx context.Context, teamID int64, id int64) (*model.User, error)
	FindByGithubID(ctx context.Context, githubID string) (*model.User, error)
	FindBySlackID(ctx context.Context, slackID string) (*model.User, error)
	Delete(ctx context.Context, teamID int64, id int64) error
}

type userRepository struct {
	dsClient *datastore.Client
}

func NewUserRepository(client *datastore.Client) UserRepository {
	return &userRepository{dsClient: client}
}

func (r *userRepository) generateKey(teamID int64) *datastore.Key {
	return datastore.IncompleteKey(model.UserKind, r.parentKey(teamID))
}

func (r *userRepository) key(id int64, teamID int64) *datastore.Key {
	return datastore.IDKey(model.UserKind, id, r.parentKey(teamID))
}

func (r *userRepository) parentKey(teamID int64) *datastore.Key {
	return datastore.IDKey(model.TeamKind, teamID, nil)
}

func (r *userRepository) List(ctx context.Context, teamID int64, cursor string, limit int) ([]*model.User, string, error) {
	q := datastore.NewQuery(model.UserKind).Ancestor(r.parentKey(teamID))
	if cursor != "" {
		dsCursor, err := datastore.DecodeCursor(cursor)
		if err != nil {

		}
		q = q.Start(dsCursor)
	}
	q = q.Limit(limit)

	var el []*model.User
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

type Reviewer struct {
	SlackID string
}

func (r *userRepository) Get(ctx context.Context, teamID int64, id int64) (*model.User, error) {
	dst := &model.User{}
	err := r.dsClient.Get(ctx, r.key(id, teamID), dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (r *userRepository) FindByGithubID(ctx context.Context, githubID string) (*model.User, error) {
	q := datastore.NewQuery(model.UserKind).KeysOnly().Filter("githubId =", githubID).Limit(1)

	keys, err := r.dsClient.GetAll(context.Background(), q, nil)

	if err != nil {
		return nil, fmt.Errorf("not found user keys %w", err)
	}

	var user model.User

	if err := r.dsClient.Get(ctx, keys[0], &user); err != nil {
		return nil, fmt.Errorf("not found user %w", err)
	}

	user.ID = keys[0].ID

	return &user, nil
}

func (r *userRepository) FindBySlackID(ctx context.Context, slackID string) (*model.User, error) {
	q := datastore.NewQuery(model.UserKind).KeysOnly().Filter("slackId =", slackID).Limit(1)

	keys, err := r.dsClient.GetAll(ctx, q, nil)

	if err != nil {
		return nil, fmt.Errorf("not found user %w", err)
	}

	var user model.User

	if err := r.dsClient.Get(ctx, keys[0], &user); err != nil {
		return nil, fmt.Errorf("not found user %w", err)
	}

	user.ID = keys[0].ID

	return &user, nil
}

func (r *userRepository) Put(ctx context.Context, teamID int64, user *model.User) error {
	_, err := r.dsClient.Put(ctx, r.key(user.ID, teamID), user)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, teamID int64, id int64) error {
	return r.dsClient.Delete(ctx, r.key(id, teamID))
}
