package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/hayashiki/mentions/pkg/model"
	"google.golang.org/api/iterator"
)

type GroupRepository interface {
	Put(teamID string, group *model.Group) error
	List(teamID string, cursor string, limit int) ([]*model.Group, string, error)
}

type groupRepository struct {
	dsClient *datastore.Client
}

func (r *groupRepository) key(teamID string) *datastore.Key {
	return datastore.IncompleteKey(model.GroupKind, r.parentKey(teamID))
}

func (r *groupRepository) parentKey(teamID string) *datastore.Key {
	return datastore.NameKey(model.TeamKind, teamID, nil)
}

func (r groupRepository) Put(teamID string, group *model.Group) error {
	ctx := context.Background()
	_, err := r.dsClient.Put(ctx, r.key(teamID), group)
	if err != nil {
		return err
	}
	return nil
}

func (r groupRepository) List(teamID string, cursor string, limit int) ([]*model.Group, string, error) {
	q := datastore.NewQuery(model.GroupKind).Ancestor(r.parentKey(teamID))
	if cursor != "" {
		dsCursor, err := datastore.DecodeCursor(cursor)
		if err != nil {

		}
		q = q.Start(dsCursor)
	}
	q = q.Limit(limit)

	var el []*model.Group
	ctx := context.Background()
	it := r.dsClient.Run(ctx, q)
	for {
		var e model.Group
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

func NewGroupRepository(client *datastore.Client) GroupRepository {
	return &groupRepository{dsClient: client}
}
