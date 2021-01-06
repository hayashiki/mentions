package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/mentions/pkg/model"
)

//go:generate mockgen -source installation_repository.go -destination mock_repo/installation_repository.go
type InstallationRepository interface {
	Put(ctx context.Context, installation *model.Installation) error
	Get(ctx context.Context, id int64) (*model.Installation, error)
}

type installationRepository struct {
	dsClient *datastore.Client
}

func NewInstallationRepository(client *datastore.Client) InstallationRepository {
	return &installationRepository{dsClient: client}
}

func (r *installationRepository) key(id int64) *datastore.Key {
	return datastore.IDKey(model.InstallationKind, id, nil)
}

func (r *installationRepository) Get(ctx context.Context, id int64) (*model.Installation, error) {
	dst := &model.Installation{}
	err := r.dsClient.Get(ctx, r.key(id), dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (r *installationRepository) Put(ctx context.Context, installation *model.Installation) error {
	_, err := r.dsClient.Put(ctx, r.key(installation.ID), installation)
	if err != nil {
		return err
	}
	return nil
}

func (r *installationRepository) Delete(ctx context.Context, id int64) error {
	return r.dsClient.Delete(ctx, r.key(id))
}
