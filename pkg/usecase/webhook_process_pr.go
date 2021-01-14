package usecase

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/pkg/event"
	"github.com/hayashiki/mentions/pkg/memcache"
	"github.com/hayashiki/mentions/pkg/slack"
	log "github.com/sirupsen/logrus"
)

func (w *webhookProcess) processPullRequestComment(ctx context.Context, ghEvent *github.PullRequestReviewCommentEvent) error {
	ev := event.NewPullRequestCommentEvent(ghEvent)

	conf := memcache.NewClient(w.config.MemcachedServer, w.config.MemcachedUsername, w.config.MemcachedPassword)
	mem, quit := memcache.NewCommentCache(conf)
	defer quit()

	task, err := w.taskRepo.Get(ctx, ev.Repository.ID)
	slackSvc := slack.NewClient(slack.New(task.Team.Token))
	users, _, err := w.userRepo.List(ctx, task.Team, "", 100)
	if err != nil {
		return err
	}
	task.Users = users

	if err != nil {
		return err
	}

	payload := slack.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	if comment, ok := slackSvc.ConvertComment(payload, task.Users); ok {
		slackMessageCache, err := mem.Get(ev.IssueCacheKey())

		var postResp slack.MessageResponse

		// ヒットした場合 == スレッド表示したい
		var ts string
		if slackMessageCache != nil {
			ts = postResp.Timestamp
		}

		slackSvc := slack.NewClient(slack.New(task.Team.Token))
		users, _, err := w.userRepo.List(ctx, task.Team, "", 100)
		if err != nil {
			return err
		}
		task.Users = users

		log.Printf("task is %v", task.Channel)
		resp, err := slackSvc.PostMessage(task.Channel, ts, comment)
		if err != nil {
			return err
		}

		// セットしなおし不要
		// スレッドキャッシュがない場合 つまり最初の投稿の場合にキャッシュする
		if ts == "" {
			err = mem.Set(ev.IssueCacheKey(), resp)
			log.Printf("memcached, %v", err)
		}

		err = mem.Set(ev.CommentCacheKey(), resp)
		log.Printf("memcached, %v", err)
	}

	return nil
}
