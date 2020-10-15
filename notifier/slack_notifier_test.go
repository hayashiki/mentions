package notifier

import (
	"github.com/hayashiki/mentions/model"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestSlackNotifier_Notify(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://hooks.slack.com/services/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",
		httpmock.NewStringResponder(200, `{}`))

	type args struct {
		webhookURL string
		message    string
	}

	type want struct {
		comment string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "simple",
			args: args{
				webhookURL: "https://hooks.slack.com/services/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",
				message:    "message <@hayashiki>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewSlackNotifier()
			err := n.Notify(tt.args.webhookURL, tt.args.message)
			assert.NoError(t, err)
		})
	}
}

func TestSlackNotifier_ConvertComment(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://hooks.slack.com/services/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",
		httpmock.NewStringResponder(200, `{}`))

	type args struct {
		payload ConvertPayload
		users   []model.User
	}
	type want struct {
		comment   string
		converted bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "convert success",
			args: args{
				payload: ConvertPayload{
					Comment:  "hi @hayashiki, please fix this bug.",
					RepoName: "hayashiki/sample",
					HTMLURL:  "https://example.com",
					Title:    "Cant open top page",
					User:     "hayashiki",
				},
				users: []model.User{
					model.User{
						ID:        "hayashiki",
						Workspace: "hayashiki",
						SlackID:   "UD7AKTEFK",
						GithubID:  "hayashiki",
						Reviewers: nil,
						CreatedAt: time.Now(),
					},
				},
			},
			want: want{
				comment:   "*hayashiki/sample <https://example.com|Cant open top page> * by: hayashiki\nhi <@UD7AKTEFK>, please fix this bug.",
				converted: true,
			},
		},
		{
			name: "not convert",
			args: args{
				payload: ConvertPayload{
					Comment:  "hi, please fix this bug.",
					RepoName: "hayashiki/sample",
					HTMLURL:  "https://example.com",
					Title:    "Cant open top page",
					User:     "hayashiki",
				},
				users: []model.User{
					model.User{
						ID:        "hayashiki",
						Workspace: "hayashiki",
						SlackID:   "UD7AKTEFK",
						GithubID:  "hayashiki",
						Reviewers: nil,
						CreatedAt: time.Now(),
					},
				},
			},
			want: want{
				//comment: "*hayashiki/sample <https://example.com|Cant open top page> * by: hayashiki\nhi <@UD7AKTEFK>, please fix this bug.",
				converted: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewSlackNotifier()

			got, ok := n.ConvertComment(tt.args.payload, tt.args.users)

			if ok != tt.want.converted {
				t.Errorf("ok returned: %v", ok)
			}

			log.Printf("got %s", got)
		})
	}
}
