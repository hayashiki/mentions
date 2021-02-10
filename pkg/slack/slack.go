package slack

import (
	"github.com/hayashiki/mentions/pkg/model"
	"github.com/slack-go/slack"
	"regexp"
	"strings"
)

var r = regexp.MustCompile(`@[a-zA-Z0-9_\-\.]+`)

type Client interface {
	GetUsers() ([]*model.User, error)
	GetUser(id string) (*model.User, error)
	PostMessage(channel, ts, message string) (*MessageResponse, error)
	UpdateMessage(channel, ts, message string) (*MessageResponse, error)
	ConvertComment(payload ConvertPayload, users []*model.User) (convertMessage string, ok bool)
}

type client struct {
	bot *slack.Client
}

func NewClient(cli *slack.Client) Client {
	return &client{bot: cli}
}

func New(token string) *slack.Client {
	return slack.New(token)
}

// GetUsers
func (c *client) GetUsers() ([]*model.User, error) {
	slackUsers, err := c.bot.GetUsers()
	if err != nil {
		return nil, err
	}

	//users := make([]*model.User, len(slackUsers))
	//for i, user := range slackUsers {
	//	if user.IsBot { continue }
	//	if user.Deleted { continue }
	//	if user.IsInvitedUser { continue }
	//	users[i] = &model.User{
	//		ID:     user.ID,
	//		Name:   user.Name,
	//		Avatar: user.Profile.Image192,
	//		TeamID: user.TeamID,
	//	}
	//}

	var users []*model.User
	for _, user := range slackUsers {

		if user.IsBot {
			continue
		}
		if user.Deleted {
			continue
		}
		if user.IsInvitedUser {
			continue
		}

		name := user.Profile.DisplayName
		if name == "" {
			name = user.Name
		}

		//IsRestricted、IsAdmin、IsOwnerがほしい

		users = append(users, &model.User{
			ID:     user.ID,
			Name:   name,
			Avatar: user.Profile.Image192,
			TeamID: user.TeamID,
		})
	}

	return users, nil
}

func (c *client) GetUser(id string) (*model.User, error) {
	user, err := c.bot.GetUserInfo(id)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:     user.ID,
		Name:   user.Profile.DisplayName,
		Avatar: user.Profile.Image192,
	}, nil
}

// ReplaceComment replace github account to bot
func (c *client) ConvertComment(payload ConvertPayload, users []*model.User) (convertMessage string, ok bool) {

	// okというかfoundかな
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
	Channel   string
	Timestamp string
}

func (c *client) PostMessage(channel, ts, message string) (*MessageResponse, error) {

	opts := []slack.MsgOption{
		slack.MsgOptionText(message, false),
		slack.MsgOptionTS(ts),
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
