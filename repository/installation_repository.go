package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/mentions/model"
)

//go:generate mockgen -source repo_repository.go -destination mock_repo/repo_repository.go
type InstallationRepository interface {
	Put(installation *model.Installation) error
	GetByID(ID int64) (*model.Installation, error)
}

type installationRepository struct {
	dsClient *datastore.Client
}

func NewInstallationRepository(client *datastore.Client) InstallationRepository {
	return &installationRepository{dsClient: client}
}

func (r *installationRepository) GetByID(ID int64) (*model.Installation, error) {
	ctx := context.Background()
	dst := &model.Installation{}
	k := datastore.IDKey(model.InstallationKind, ID, nil)
	err := r.dsClient.Get(ctx, k, dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (r *installationRepository) Put(installation *model.Installation) error {
	ctx := context.Background()
	k := datastore.IDKey(model.InstallationKind, installation.ID, nil)
	_, err := r.dsClient.Put(ctx, k, installation)
	if err != nil {
		return err
	}
	return nil
}
