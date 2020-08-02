package notifier

import (
	"github.com/hayashiki/mentions/account"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlackNotifier_Notify(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://hooks.slack.com/services/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",
		httpmock.NewStringResponder(200, `{}`))

	type fields struct {
		AccountList account.List
	}

	type args struct {
		webhookURL string
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "simple",
			fields: fields{
				AccountList: account.List{},
			},
			args: args{
				webhookURL: "https://hooks.slack.com/services/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",
				message: "message <@hayashiki>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &SlackNotifier{
				AccountList: tt.fields.AccountList,
			}
			err := n.Notify(tt.args.webhookURL, tt.args.message)
			assert.NoError(t, err)
		})
	}
}
