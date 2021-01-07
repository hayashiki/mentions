package usecase

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/pkg/event"
	ghSvc "github.com/hayashiki/mentions/pkg/github"
	"github.com/hayashiki/mentions/pkg/mem"
	"github.com/hayashiki/mentions/pkg/model"
	"github.com/hayashiki/mentions/pkg/slack"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (w *webhookProcess) editIssue(task *model.Task, ev *event.Event) error {
	log.Debugf("editIssue by: %v", ev.IssueOwner)
	user, found := task.GetUserByGithubID(ev.IssueOwner)
	if !found {
		log.Errorf("github user not found user %s", ev.IssueOwner)
		return fmt.Errorf("github user not found user %s", ev.IssueOwner)
	}
	log.Debugf("user is %+v", user)

	payload := &ghSvc.CreateReviewersPayload{
		Owner:       ev.Repository.Owner,
		Name:        ev.Repository.Name,
		IssueNumber: ev.IssueNumber,
		Reviewers:   user.Reviewers.String(),
	}

	_, resp, err := w.ghSvc.CreateReviewers(payload)
	if err != nil {
		return fmt.Errorf("failed to create reviewer resp %v, err=%v", resp, err)
	}

	comment := strings.Join(user.ReviewersWithAt(), " ") + " レビューお願いします😀"
	ev.Comment = comment

	commentPayload := &ghSvc.EditIssueCommentPayload{
		Owner:     ev.Repository.Owner,
		Name:      ev.Repository.Name,
		CommentID: ev.CommentID,
		Comment:   ev.Comment,
	}

	log.Printf("ev.InstallationID is %v", ev.InstallationID)
	if ev.InstallationID != 0 {
		ghAppCli := ghSvc.NewClient(ghSvc.GetAppClient(
			w.config.GithubAppID,
			ev.InstallationID,
			w.config.GithubAppPrivateKeyFileName,
		))
		_, resp, err = ghAppCli.EditIssueComment(commentPayload)
		if err != nil {
			return fmt.Errorf("failed to edit issue resp %v, err=%v", resp, err)
		}
	} else {
		_, resp, err = w.ghSvc.EditIssueComment(commentPayload)
		log.Printf("ev.InstallationID b")
		if err != nil {
			return fmt.Errorf("failed to edit issue resp %v, err=%v", resp, err)
		}
	}
	return nil
}

func (w *webhookProcess) processEditIssueComment(ctx context.Context, ghEvent *github.IssueCommentEvent) error {
	log.Debug("called processEditIssueComment")
	conf := mem.NewConfig(w.config.MemcachedServer, w.config.MemcachedUsername, w.config.MemcachedPassword)
	mem, quit := mem.NewCommentCache(conf)
	defer quit()

	ev := event.NewIssueComment(ghEvent)
	issueCommentKey := strconv.Itoa(int(ev.CommentID))
	slackMessageCache, err := mem.Get(issueCommentKey)
	task, err := w.taskRepo.Get(ctx, ev.Repository.ID)
	if err != nil {
		return err
	}
	slackSvc := slack.NewClient(slack.New(task.Team.Token))
	users, _, err := w.userRepo.List(ctx, task.Team, "", 100)
	if err != nil {
		return err
	}
	task.Users = users

	if err != nil {
		return fmt.Errorf("failed to get task %v", err)
	}

	payload := slack.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	if comment, ok := slackSvc.ConvertComment(payload, task.Users); ok {
		// キャッシュしなおしてもいいかも
		if slackMessageCache != nil {
			if _, err := slackSvc.UpdateMessage(slackMessageCache.Channel, slackMessageCache.Timestamp, comment); err != nil {
				log.Printf("n debug err %v", err)
				return err
			}
		} else {
			// r?のケース
			issueNumberKey := strconv.Itoa(int(ev.IssueNumber))
			slackMessageCache, err = mem.Get(issueNumberKey)
			var ts string
			// r?のケースで手前に一度でもMentions経由投稿があった場合
			if slackMessageCache != nil {
				ts = slackMessageCache.Timestamp
			}
			log.Printf("task is %v", task.Channel)
			if _, err := slackSvc.PostMessage(task.Channel, ts, comment); err != nil {
				log.Printf("n debug err %v", err)
				return err
			}
		}
	}

	return nil
}

func (w *webhookProcess) processIssueComment(ctx context.Context, ghEvent *github.IssueCommentEvent) error {
	log.Printf("Called.processCreateIssueComment")

	conf := mem.NewConfig(w.config.MemcachedServer, w.config.MemcachedUsername, w.config.MemcachedPassword)
	mem, quit := mem.NewCommentCache(conf)
	defer quit()

	//ghEvent.Installation.IDをつかってteamsを判定する
	ev := event.NewIssueComment(ghEvent)

	// 複数になる
	log.Debugf("ev.Repository.ID %v", ev.Repository.ID)
	task, err := w.taskRepo.Get(ctx, ev.Repository.ID)

	if err != nil {
		return fmt.Errorf("failed to get task %v", err)
	}

	var slackSvc slack.Client
	if task.Team == nil {
		log.Errorf("task.Team is not exists task: %v", task)
		return fmt.Errorf("task.Team is not exists %v", task)
	}
	slackSvc = slack.NewClient(slack.New(task.Team.Token))

	users, _, err := w.userRepo.List(ctx, task.Team, "", 100)
	if err != nil {
		return err
	}
	task.Users = users

	if hasReviewMagicWord(ev.Comment) {
		if err := w.editIssue(task, ev); err != nil {
			return fmt.Errorf("failed to edit github issue %v", err)
		}
		log.Infof("gor r? word issue: %v", ev.IssueNumber)
		// Edit処理でWebhookがとぶのでそれで再度editする
		return nil
	}

	payload := slack.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	if comment, ok := slackSvc.ConvertComment(payload, task.Users); ok {
		log.Debug("Convert Comment")
		issueNumberKey := strconv.Itoa(int(ev.IssueNumber))
		slackMessageCache, err := mem.Get(issueNumberKey)

		// ヒットした場合 == スレッド表示したい
		var ts string
		if slackMessageCache != nil {
			ts = slackMessageCache.Timestamp
		}
		resp, err := slackSvc.PostMessage(task.Channel, ts, comment)
		if err != nil {
			log.Printf("err is %v", err)
			return err
		}
		issueCommentKey := strconv.Itoa(int(ev.CommentID))

		log.WithFields(log.Fields{
			"ts": ts,
			"issueNumberKey": issueNumberKey,
			"issueCommentKey": issueCommentKey,
		})
		// 最初の投稿の場合にキャッシュする
		if ts == "" {
			err = mem.Set(issueNumberKey, resp)
		}
		if err := mem.Set(issueCommentKey, resp); err != nil {
			return err
		}
	}
	return nil
}
