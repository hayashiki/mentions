package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/hayashiki/mentions/pkg/model"
	"google.golang.org/api/iterator"
)

//go:generate mockgen -source repo_repository.go -destination mock_repo/repo_repository.go
type RepoRepository interface {
	List(ctx context.Context, cursor string, limit int) ([]*model.Repo, string, error)
	Get(ctx context.Context, id int64) (*model.Repo, error)
	Put(ctx context.Context, repo *model.Repo) error
	Delete(ctx context.Context, id int64) error
}

type repoRepository struct {
	dsClient *datastore.Client
}

func NewRepoRepository(client *datastore.Client) RepoRepository {
	return &repoRepository{dsClient: client}
}

func (r *repoRepository) key(id int64) *datastore.Key {
	return datastore.IDKey(model.RepoKind, id, nil)
}

func (r *repoRepository) List(ctx context.Context, cursor string, limit int) ([]*model.Repo, string, error) {
	q := datastore.NewQuery(model.RepoKind)
	if cursor != "" {
		dsCursor, err := datastore.DecodeCursor(cursor)
		if err != nil {

		}
		q = q.Start(dsCursor)
	}
	q = q.Limit(limit)

	var el []*model.Repo
	it := r.dsClient.Run(ctx, q)
	for {
		var e model.Repo
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

func (r *repoRepository) Get(ctx context.Context, id int64) (*model.Repo, error) {
	dst := &model.Repo{}
	err := r.dsClient.Get(ctx, r.key(id), dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (r *repoRepository) Put(ctx context.Context, repo *model.Repo) error {
	_, err := r.dsClient.Put(ctx, r.key(repo.ID), repo)
	if err != nil {
		return err
	}
	return nil
}

func (r *repoRepository) Delete(ctx context.Context, id int64) error {
	return r.dsClient.Delete(ctx, r.key(id))
}
