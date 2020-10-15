package model

import (
	"reflect"
	"testing"
	"time"
)

func TestGithubID_String(t *testing.T) {
	tests := []struct {
		name string
		g    GithubID
		want string
	}{
		{
			name: "success",
			g:    "github",
			want: "github",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubID_WithAt(t *testing.T) {
	tests := []struct {
		name string
		g    GithubID
		want string
	}{
		{
			name: "@ are added at the beginning",
			g:    "github",
			want: "@github",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.WithAt(); got != tt.want {
				t.Errorf("WithAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReviewers_String(t *testing.T) {
	tests := []struct {
		name string
		rs   Reviewers
		want []string
	}{
		{
			name: "@ are added at the beginning",
			rs:   Reviewers{"github"},
			want: []string{"@github"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rs.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_GetUserByGithubID(t1 *testing.T) {
	type fields struct {
		ID         int64
		Workspace  string
		Repo       Repo
		WebhookURL string
		Users      []User
	}
	type args struct {
		githubID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *User
		want1  bool
	}{
		{
			name: "success",
			fields: fields{
				ID:         1,
				Workspace:  "workspace",
				Repo:       Repo{ID: "example", Owner: "repo", Name: "example"},
				WebhookURL: "",
				Users: []User{
					User{
						ID:       "user",
						GithubID: "hayashiki",
					},
				},
			},
			args: args{
				githubID: "hayashiki",
			},
			want: &User{
				ID:       "user",
				GithubID: "hayashiki",
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Task{
				ID:         tt.fields.ID,
				Workspace:  tt.fields.Workspace,
				Repo:       tt.fields.Repo,
				WebhookURL: tt.fields.WebhookURL,
				Users:      tt.fields.Users,
			}
			got, got1 := t.GetUserByGithubID(tt.args.githubID)
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetUserByGithubID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t1.Errorf("GetUserByGithubID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUser_GithubWithAt(t *testing.T) {
	type fields struct {
		ID        string
		Workspace string
		SlackID   string
		GithubID  GithubID
		Reviewers Reviewers
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				ID:        tt.fields.ID,
				Workspace: tt.fields.Workspace,
				SlackID:   tt.fields.SlackID,
				GithubID:  tt.fields.GithubID,
				Reviewers: tt.fields.Reviewers,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if got := u.GithubWithAt(); got != tt.want {
				t.Errorf("GithubWithAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_ReviewersWithAt(t *testing.T) {
	type fields struct {
		ID        string
		Workspace string
		SlackID   string
		GithubID  GithubID
		Reviewers Reviewers
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				ID:        tt.fields.ID,
				Workspace: tt.fields.Workspace,
				SlackID:   tt.fields.SlackID,
				GithubID:  tt.fields.GithubID,
				Reviewers: tt.fields.Reviewers,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if got := u.ReviewersWithAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReviewersWithAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_SlackWithBracketAt(t *testing.T) {
	type fields struct {
		ID        string
		Workspace string
		SlackID   string
		GithubID  GithubID
		Reviewers Reviewers
		CreatedAt time.Time
		UpdatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				ID:        tt.fields.ID,
				Workspace: tt.fields.Workspace,
				SlackID:   tt.fields.SlackID,
				GithubID:  tt.fields.GithubID,
				Reviewers: tt.fields.Reviewers,
				CreatedAt: tt.fields.CreatedAt,
				UpdatedAt: tt.fields.UpdatedAt,
			}
			if got := u.SlackWithBracketAt(); got != tt.want {
				t.Errorf("SlackWithBracketAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
