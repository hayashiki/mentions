package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hayashiki/mentions/model"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var r = regexp.MustCompile(`@[a-zA-Z0-9_\-\.]+`)

type Notifier interface {
	Notify(webhookURL, message string) error
	ConvertComment(payload ConvertPayload, users []model.User) (convertMessage string, ok bool)
}

type SlackNotifier struct{}

func NewSlackNotifier() Notifier {
	return &SlackNotifier{}
}

type ConvertPayload struct {
	Comment  string
	RepoName string
	HTMLURL  string
	Title    string
	User     string
}

type PostMessageRequest struct {
	Text      string `json:"text,omitempty"`
	LinkNames string `json:"link_names,omitempty"`
}

func (p ConvertPayload) buildMessage() string {
	var text string
	text = fmt.Sprintf("*%v <%v|%v> * by: %v", p.RepoName, p.HTMLURL, p.Title, p.User)
	text = fmt.Sprintf("%v\n%v", text, p.Comment)
	return text
}

func (n *SlackNotifier) Notify(webhookURL, message string) error {

	pm := PostMessageRequest{
		Text:      message,
		LinkNames: "1",
	}

	body, err := json.Marshal(pm)

	if err != nil {
		fmt.Errorf("failed to marshal to byte, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader([]byte(body)))

	if err != nil {
		fmt.Errorf("failed to create request, err: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send a http request, err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read request body, err: %w", err)
		}
		return errors.New(string(b))
	}
	return nil
}

// ReplaceComment replace github account to slack
func (n *SlackNotifier) ConvertComment(payload ConvertPayload, users []model.User) (convertMessage string, ok bool) {

	ok = false
	// eg. hello @hayashiki , I hava a question
	matches := r.FindAllStringSubmatch(payload.Comment, -1)

	if len(matches) == 0 {
		return payload.Comment, ok
	}

	for _, val := range matches {
		//eg. val[0] -> @hayashiki
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
