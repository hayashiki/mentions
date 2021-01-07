package usecase

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/hayashiki/mentions/config"
	"github.com/hayashiki/mentions/event"
	ghSvc "github.com/hayashiki/mentions/github"
	"github.com/hayashiki/mentions/mem"
	"github.com/hayashiki/mentions/pkg/model"
	"github.com/hayashiki/mentions/pkg/repository"
	"github.com/hayashiki/mentions/slack"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type WebhookProcess interface {
	Do(r *http.Request) error
}

type webhookProcess struct {
	config        config.Config
	githubService ghSvc.Github
	ghAppSvc      ghSvc.Github
	userRepo      repository.UserRepository
	taskRepo      repository.TaskRepository
	repoRepo      repository.RepoRepository
	teamRepo      repository.TeamRepository
}

func NewWebhookProcess(
	env config.Config,
	ghSvc ghSvc.Github,
	ghAppSvc ghSvc.Github,
	userRepo repository.UserRepository,
	taskRepo repository.TaskRepository,
	repoRepo repository.RepoRepository,
) WebhookProcess {
	return &webhookProcess{
		config:        env,
		githubService: ghSvc,
		ghAppSvc:      ghAppSvc,
		userRepo:      userRepo,
		taskRepo:      taskRepo,
		repoRepo:      repoRepo,
	}
}

func (w webhookProcess) Do(r *http.Request) error {
	// TODO: get from secret manager
	//payload, err := ghSvc.Verify(r, []byte(w.config.GithubWebhookSecret))
	// TODO: webhookProcessのinitでnewする
	payload, err := ghSvc.NewGithubValidator().Validate(r, []byte(w.config.GithubWebhookSecret))
	if err != nil {
		return err
	}

	ghEvent, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return err
	}

	switch ev := ghEvent.(type) {
	case *github.IssueCommentEvent:
		switch ev.GetAction() {
		case "created":
			return w.processIssueComment(r.Context(), ev)
		case "edited":
			return w.processEditIssueComment(r.Context(), ev)
		}
	case *github.PullRequestReviewCommentEvent:
		return w.processPullRequestComment(r.Context(), ev)
	case *github.InstallationRepositoriesEvent:
		log.Printf("action is %v", ev.GetAction())
		switch ev.GetAction() {
		case "added":
			return w.processInstallationRepositoriesEvent(r.Context(), ev)
		case "removed":
			return w.processInstallationReposRemovedEvent(r.Context(),ev)
		}
	case *github.InstallationEvent:
		log.Printf("ev2 is %+v", ev)
	default:
		return nil
	}
	return nil
}

func (w *webhookProcess) processEditIssueComment(ctx context.Context, ghEvent *github.IssueCommentEvent) error {
	log.Printf("Called.processEditIssueComment")

	conf := mem.NewConfig("memcached-16535.c1.asia-northeast1-1.gce.cloud.redislabs.com:16535", "mc-KpxsD", "FRvZcLiVqPSFcMA98tgendBx1O4EXBRJ")
	mem, quit := mem.NewCommentCache(conf)
	defer quit()

	ev := event.NewIssueComment(ghEvent)
	issueCommentKey := strconv.Itoa(int(ev.CommentID))
	slackMessageCache, err := mem.Get(issueCommentKey)
	task, err := w.taskRepo.Get(ctx, ev.Repository.ID)
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

	conf := mem.NewConfig("memcached-16535.c1.asia-northeast1-1.gce.cloud.redislabs.com:16535", "mc-KpxsD", "FRvZcLiVqPSFcMA98tgendBx1O4EXBRJ")
	mem, quit := mem.NewCommentCache(conf)
	defer quit()

	//ghEvent.Installation.IDをつかってteamsを判定する

	ev := event.NewIssueComment(ghEvent)

	// 複数になる
	log.Printf("ev.Repository.ID %v", ev.Repository.ID)
	task, err := w.taskRepo.Get(ctx, ev.Repository.ID)

	if err != nil {
		return fmt.Errorf("failed to get task %v", err)
	}

	var slackSvc slack.Client
	if task.Team == nil {
		log.Fatalf("task.Team is not exsits")
		return fmt.Errorf("task.Team is not exsits %v", task)
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
		log.Printf("return")
		// このeditでまたwebhookがとぶのでそれでeditする
		return nil
	}

	payload := slack.ConvertPayload{
		Comment:  ev.Comment,
		RepoName: ev.Repository.FullName,
		HTMLURL:  ev.HTMLURL,
		Title:    ev.Title,
		User:     ev.User,
	}

	log.Printf("users  %v", task.Users)
	if comment, ok := slackSvc.ConvertComment(payload, task.Users); ok {

		log.Printf("ev.Cnvert")
		issueNumberKey := strconv.Itoa(int(ev.IssueNumber))
		slackMessageCache, err := mem.Get(issueNumberKey)
		log.Printf("is  err %v", err)

		// ヒットした場合 == スレッド表示したい
		var ts string
		log.Printf("slackMessageCache exists check start %v", slackMessageCache)

		if slackMessageCache != nil {
			log.Printf("slackMessageCachee %v", slackMessageCache)
			ts = slackMessageCache.Timestamp
		}

		log.Printf("ts  is %v", ts)

		log.Printf("task is %v", task.Channel)
		resp, err := slackSvc.PostMessage(task.Channel, ts, comment)

		log.Printf("resp  is %v", resp)

		if err != nil {
			log.Printf("err is %v", err)
			return err
		}
		issueCommentKey := strconv.Itoa(int(ev.CommentID))

		log.Printf("[create] issueNumberKey is %v", issueNumberKey)
		log.Printf("[create] issueCommentKey is %v", issueCommentKey)
		log.Printf("cached, %v", resp)

		// セットしなおし不要
		// スレッドキャッシュがない場合 つまり最初の投稿の場合にキャッシュする
		if ts == "" {
			err = mem.Set(issueNumberKey, resp)
			log.Printf("memcached, %v", err)
		}

		err = mem.Set(issueCommentKey, resp)
		log.Printf("memcached, %v", err)
	}
	return nil
}

func (w *webhookProcess) processPullRequestComment(ctx context.Context, ghEvent *github.PullRequestReviewCommentEvent) error {
	ev := event.NewPullRequestCommentEvent(ghEvent)

	conf := mem.NewConfig("memcached-16535.c1.asia-northeast1-1.gce.cloud.redislabs.com:16535", "mc-KpxsD", "FRvZcLiVqPSFcMA98tgendBx1O4EXBRJ")
	mem, quit := mem.NewCommentCache(conf)
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
		issueNumberKey := strconv.Itoa(ev.IssueNumber)
		issueCommentKey := strconv.Itoa(int(ev.CommentID))
		log.Printf("commentID %d", ev.CommentID)

		slackMessageCache, err := mem.Get(issueNumberKey)

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
			err = mem.Set(issueNumberKey, resp)
			log.Printf("memcached, %v", err)
		}

		err = mem.Set(issueCommentKey, resp)
		log.Printf("memcached, %v", err)
	}

	return nil
}

func (w *webhookProcess) editIssue(task *model.Task, ev *event.Event) error {
	log.Printf("editIssue %v", ev.IssueOwner)
	user, ok := task.GetUserByGithubID(ev.IssueOwner)
	if !ok {
		return fmt.Errorf("github user not found user %s", ev.IssueOwner)
	}
	log.Printf("user is %+v", user)

	payload := &ghSvc.CreateReviewersPayload{
		Owner:       ev.Repository.Owner,
		Name:        ev.Repository.Name,
		IssueNumber: ev.IssueNumber,
		Reviewers:   user.Reviewers.String(),
	}

	log.Printf("reviews is %+v", user.Reviewers.String())

	_, resp, err := w.githubService.CreateReviewers(payload)
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
		log.Printf("ev.InstallationID a")

		appID, err := strconv.Atoi(w.config.GithubAppID)
		if err != nil {
		}

		ghAppCli := ghSvc.NewClient(ghSvc.GetAppClient(
			int64(appID),
			int64(ev.InstallationID),
			w.config.GithubAppPrivateKeyFileName,
		))

		_, resp, err = ghAppCli.EditIssueComment(commentPayload)
		if err != nil {
			return fmt.Errorf("failed to edit issue resp %v, err=%v", resp, err)
		}
	} else {
		_, resp, err = w.githubService.EditIssueComment(commentPayload)
		log.Printf("ev.InstallationID b")
		if err != nil {
			return fmt.Errorf("failed to edit issue resp %v, err=%v", resp, err)
		}
	}
	return nil
}

func hasReviewMagicWord(s string) bool {
	return string([]rune(s)[:2]) == "r?"
}

func (w *webhookProcess) processInstallationRepositoriesEvent(ctx context.Context, ghEvent *github.InstallationRepositoriesEvent) error {
	log.Printf("added repo event called")
	repos := event.NewInstallationRepositoriesEvent(ghEvent)
	// TODO addedByがほしい
	for _, repo := range repos {
		err := w.repoRepo.Put(ctx, &model.Repo{
			ID: repo.ID,
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
