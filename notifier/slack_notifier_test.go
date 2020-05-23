package notifier

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlackNotifier_Notify(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://hooks.slack.com/services/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",
		httpmock.NewStringResponder(200, ""))

	type fields struct {
		WebhookURL string
	}

	type args struct {
		payload *PostMessageRequest
	}
	tests := []struct{
		name string
		fields fields
		args args
	}{
		{
			name: "simple",
			fields: fields{
				WebhookURL: "https://hooks.slack.com/services/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",
			},
			args: args{
				&PostMessageRequest{
					Text: "text",
					LinkNames: "1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &SlackNotifier{
				WebhookURL: tt.fields.WebhookURL,
			}
			err := n.Notify(tt.args.payload)
			assert.NoError(t, err)
		})
	}
}
