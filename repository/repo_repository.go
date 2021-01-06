package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/mentions/model"
)

//go:generate mockgen -source repo_repository.go -destination mock_repo/repo_repository.go
type RepoRepository interface {
	Put(repo *model.Repo) error
	GetByID(ID int64) (*model.Repo, error)
	Delete(ID int64) error
}

type repoRepository struct {
	dsClient *datastore.Client
}

func NewRepoRepository(client *datastore.Client) RepoRepository {
	return &repoRepository{dsClient: client}
}

func (r *repoRepository) Delete(ID int64) error {
	ctx := context.Background()
	k := datastore.IDKey(model.RepoKind, ID, nil)
	return r.dsClient.Delete(ctx, k)
}

func (r *repoRepository) GetByID(ID int64) (*model.Repo, error) {
	ctx := context.Background()
	dst := &model.Repo{}
	k := datastore.IDKey(model.RepoKind, ID, nil)
	err := r.dsClient.Get(ctx, k, dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (r *repoRepository) Put(repo *model.Repo) error {
	ctx := context.Background()
	k := datastore.IDKey(model.RepoKind, repo.ID, nil)
	_, err := r.dsClient.Put(ctx, k, repo)
	if err != nil {
		return err
	}
	return nil
}
