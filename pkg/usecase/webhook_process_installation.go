package usecase

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/pkg/event"
	"github.com/hayashiki/mentions/pkg/model"
	log "github.com/sirupsen/logrus"
)

func (w *webhookProcess) processInstallationReposAddedEvent(ctx context.Context, ghEvent *github.InstallationRepositoriesEvent) error {
	log.Printf("added repo event called")
	repos := event.NewInstallationRepositoriesEvent(ghEvent)
	// TODO addedByがほしい
	for _, repo := range repos {
		err := w.repoRepo.Put(ctx, &model.Repo{
			ID:       repo.ID,
			Owner:    repo.Owner,
			Name:     repo.Name,
			FullName: repo.FullName,
		})
		if err != nil {
			log.Printf("error is err: %v", err)
			return err
		}
	}
	return nil
}

func (w *webhookProcess) processInstallationReposRemovedEvent(ctx context.Context, ghEvent *github.InstallationRepositoriesEvent) error {
	log.Printf("delete repo event called")
	repos := event.NewDeleteRepos(ghEvent)
	// TODO addedByがほしい
	for _, repo := range repos {
		log.Printf("repo.ID: %v", repo.ID)
		err := w.repoRepo.Delete(ctx, repo.ID)
		if err != nil {
			log.Printf("error is err: %v", err)
			return err
		}
	}
	return nil
}
