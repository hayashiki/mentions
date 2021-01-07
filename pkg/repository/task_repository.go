package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/hayashiki/mentions/pkg/model"
	"google.golang.org/api/iterator"
)

//go:generate mockgen -source task_repository.go -destination mock_repo/task_repository.go
type TaskRepository interface {
	List(ctx context.Context, cursor string, limit int) ([]*model.Task, string, error)
	Put(ctx context.Context, task *model.Task) error
	Get(ctx context.Context, id int64) (*model.Task, error)
	Delete(ctx context.Context, id int64) error
}

type taskRepository struct {
	dsClient *datastore.Client
}

func NewTaskRepository(client *datastore.Client) TaskRepository {
	return &taskRepository{dsClient: client}
}

func (r *taskRepository) key(id int64) *datastore.Key {
	return datastore.IDKey(model.TaskKind, id, nil)
}

func (r *taskRepository) List(ctx context.Context, cursor string, limit int) ([]*model.Task, string, error) {
	q := datastore.NewQuery(model.TaskKind)
	if cursor != "" {
		dsCursor, err := datastore.DecodeCursor(cursor)
		if err != nil {

		}
		q = q.Start(dsCursor)
	}
	q = q.Limit(limit)

	var el []*model.Task
	it := r.dsClient.Run(ctx, q)
	for {
		var e model.Task
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

func (r *taskRepository) Get(ctx context.Context, id int64) (*model.Task, error) {
	dst := &model.Task{}
	err := r.dsClient.Get(ctx, r.key(id), dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (r *taskRepository) Put(ctx context.Context, task *model.Task) error {
	_, err := r.dsClient.Put(ctx, r.key(task.ID), task)
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id int64) error {
	return r.dsClient.Delete(ctx, r.key(id))
}
