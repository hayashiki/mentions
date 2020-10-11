package infrastructure

import (
	"cloud.google.com/go/datastore"
	"context"
)

func GetDSClient(projectID string) *datastore.Client {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	return client
}
