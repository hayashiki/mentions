package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/mentions/model"
	"log"
)

type TeamRepository interface {
	Put(team *model.Team) error
	GetByID(id string) (*model.Team, error)
	GetAll()
}

type teamRepository struct {
	dsClient *datastore.Client
}

func (r *teamRepository) GetByID(id string) (*model.Team, error) {
	ctx := context.Background()
	k := datastore.NameKey(model.TeamKind, id, nil)
	//dst := &model.Team{}
	log.Printf("k is %v", k)
	dst := &model.Team{}
	err := r.dsClient.Get(ctx, k, dst)

	log.Printf("dst is %s", dst.Name)

	return dst, err
}

func (r *teamRepository) GetAll() {
	ctx := context.Background()
	//k := datastore.NameKey(model.TeamKind, id, nil)
	var teams []*model.Team
	q := datastore.NewQuery(model.TeamKind)
	r.dsClient.GetAll(ctx, q, &teams)

	for _, team := range teams {
		log.Printf("team is %v", team)
	}
}

func (r *teamRepository) Put(team *model.Team) error {
	ctx := context.Background()

	//k := datastore.NameKey(model.TeamKind, team.ID.String(), nil)
	_, err := r.dsClient.Put(ctx, r.key(team.ID), team)
	if err != nil {
		return err
	}
	return nil
}

func NewTeamRepository(client *datastore.Client) TeamRepository {
	return &teamRepository{dsClient: client}
}

func (r *teamRepository) key(id string) *datastore.Key {
	return datastore.NameKey(model.TeamKind, id, nil)
}
