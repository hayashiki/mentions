package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/mentions/pkg/model"
	"log"
)

//go:generate mockgen -source team_repository.go -destination mock_repo/team_repository.go
type TeamRepository interface {
	Put(ctx context.Context, team *model.Team) error
	Get(ctx context.Context, id string) (*model.Team, error)
	// TODO: cursor
	List()
}

type teamRepository struct {
	dsClient *datastore.Client
}

func NewTeamRepository(client *datastore.Client) TeamRepository {
	return &teamRepository{dsClient: client}
}

func (r *teamRepository) key(id string) *datastore.Key {
	return datastore.NameKey(model.TeamKind, id, nil)
}

func (r *teamRepository) List() {
	ctx := context.Background()
	var teams []*model.Team
	q := datastore.NewQuery(model.TeamKind)
	r.dsClient.GetAll(ctx, q, &teams)

	for _, team := range teams {
		log.Printf("team is %v", team)
	}
}

func (r *teamRepository) Get(ctx context.Context, id string) (*model.Team, error) {
	k := datastore.NameKey(model.TeamKind, id, nil)
	dst := &model.Team{}
	err := r.dsClient.Get(ctx, k, dst)
	return dst, err
}

func (r *teamRepository) Put(ctx context.Context, team *model.Team) error {
	_, err := r.dsClient.Put(ctx, r.key(team.ID), team)
	if err != nil {
		return err
	}
	return nil
}

func (r *teamRepository) Delete(ctx context.Context, id string) error {
	return r.dsClient.Delete(ctx, r.key(id))
}
