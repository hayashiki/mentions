package repository

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/hayashiki/mentions/model"
)

//go:generate mockgen -source task_repository.go -destination mock_repo/task_repository.go
type TaskRepository interface {
	Put(task *model.Task) error
	GetByID(ID int64) (*model.Task, error)
}

type taskRepository struct {
	dsClient *datastore.Client
}

func NewTaskRepository(client *datastore.Client) TaskRepository {
	return &taskRepository{dsClient: client}
}

func (r *taskRepository) GetByID(ID int64) (*model.Task, error) {
	ctx := context.Background()
	dst := &model.Task{}
	k := datastore.IDKey(model.TaskKind, ID, nil)
	err := r.dsClient.Get(ctx, k, dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func (r *taskRepository) Put(task *model.Task) error {
	ctx := context.Background()
	k := datastore.IDKey(model.TaskKind, task.ID, nil)
	_, err := r.dsClient.Put(ctx, k, task)
	if err != nil {
		return err
	}
	return nil
}
