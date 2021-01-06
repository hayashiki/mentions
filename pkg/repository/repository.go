package repository

import (
	"cloud.google.com/go/datastore"
	"context"
)

func GetClient(projectID string) *datastore.Client {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	return client
}
