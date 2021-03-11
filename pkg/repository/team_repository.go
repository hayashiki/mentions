package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/hayashiki/mentions/pkg/model"
	"log"
)

//go:generate mockgen -source team_repository.go -destination mock_repo/team_repository.go
type TeamRepository interface {
	Put(ctx context.Context, team *model.Team) error
	Get(ctx context.Context, id int64) (*model.Team, error)
	// TODO: cursor
	List()
	GetBySlackTeamID(ctx context.Context, slackTeamID string) (*model.Team, error)
}

type teamRepository struct {
	dsClient *datastore.Client
}

func NewTeamRepository(client *datastore.Client) TeamRepository {
	return &teamRepository{dsClient: client}
}

func (r *teamRepository) generateKey() *datastore.Key {
	return datastore.IncompleteKey(model.TeamKind, nil)
}

func (r *teamRepository) key(id int64) *datastore.Key {
	return datastore.IDKey(model.TeamKind, id, nil)
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

func (r *teamRepository) Get(ctx context.Context, id int64) (*model.Team, error) {
	k := datastore.IDKey(model.TeamKind, id, nil)
	dst := &model.Team{}
	err := r.dsClient.Get(ctx, k, dst)
	return dst, err
}

func (r *teamRepository) GetBySlackTeamID(ctx context.Context, slackTeamID string) (*model.Team, error) {
	q := datastore.NewQuery(model.TeamKind).KeysOnly().Filter("slackTeamId =", slackTeamID).Limit(1)

	keys, err := r.dsClient.GetAll(ctx, q, nil)

	if err != nil {
		return nil, fmt.Errorf("not found user %w", err)
	}

	var team model.Team

	if err := r.dsClient.Get(ctx, keys[0], &team); err != nil {
		return nil, fmt.Errorf("not found team %w", err)
	}

	team.ID = keys[0].ID

	return &team, nil
}

func (r *teamRepository) Put(ctx context.Context, team *model.Team) error {
	key, err := r.dsClient.Put(ctx, r.generateKey(), team)
	team.ID = key.ID
	if err != nil {
		return err
	}
	return nil
}

func (r *teamRepository) Delete(ctx context.Context, id int64) error {
	return r.dsClient.Delete(ctx, r.key(id))
}
