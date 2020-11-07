package slack

import (
	"github.com/hayashiki/mentions/model"
	"github.com/slack-go/slack"
	"regexp"
	"strings"
)

var r = regexp.MustCompile(`@[a-zA-Z0-9_\-\.]+`)

type Client interface {
	PostMessage(channel, message string) (*MessageResponse, error)
	UpdateMessage(channel, ts, message string) (*MessageResponse, error)
	ConvertComment(payload ConvertPayload, users []model.User) (convertMessage string, ok bool)
}

type client struct{
	bot *slack.Client
}

func NewClient(cli *slack.Client) Client {
	return &client{bot: cli}
}

func New(token string) *slack.Client {
	return slack.New(token)
}

type PostMessageRequest struct {
	Text      string `json:"text,omitempty"`
	LinkNames string `json:"link_names,omitempty"`
}

// ReplaceComment replace github account to bot
func (c *client) ConvertComment(payload ConvertPayload, users []model.User) (convertMessage string, ok bool) {

	ok = false
	// eg. hello @hayashiki , I hava a question
	matches := r.FindAllStringSubmatch(payload.Comment, -1)

	if len(matches) == 0 {
		return payload.Comment, ok
	}

	for _, val := range matches {
		//eg. val[0] is @hayashiki
		for _, user := range users {
			if user.GithubWithAt() == val[0] {
				payload.Comment = strings.Replace(payload.Comment, val[0], user.SlackWithBracketAt(), -1)
				ok = true
			}
		}
	}
	msg := payload.buildMessage()

	return msg, ok
}

type MessageResponse struct {
	Channel string
	Timestamp string
}

func (c *client) PostMessage(channel, message string) (*MessageResponse, error) {

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
	}
	ch, ts, err := c.bot.PostMessage(channel, opts...)

	return &MessageResponse{
		Channel:   ch,
		Timestamp: ts,
	}, err
}

func (c *client) UpdateMessage(channel, ts, message string) (*MessageResponse, error) {
	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
	}
	ch, ts, _, err := c.bot.UpdateMessage(channel, ts, opts...)
	return &MessageResponse{
		Channel:   ch,
		Timestamp: ts,
	}, err
}
