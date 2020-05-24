package account

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadAccountFromFile(t *testing.T) {
	type args struct {
		filename string
	}

	var tests = []struct {
		name string
		args args
		want List
	}{
		{
			name: "blank filename",
			args: args{
				filename: "",
			},
			want: List{},
		},
		{
			name: "with filename",
			args: args{
				filename: "testdata/github-config.json",
			},
			want: List{
				Accounts: map[string]string{"@hayashiki": "<@U021G3CDSJP>"},
				Repos:    map[string]string{"hayashiki/gcp-functions-fw": "https://hooks.slack.com/services/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"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadAccountFromFile(tt.args.filename)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
