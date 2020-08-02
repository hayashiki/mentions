package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hayashiki/mentions/account"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var r = regexp.MustCompile(`@[a-zA-Z0-9_\-\.]+`)

type SlackNotifier struct {
	AccountList account.List
}

func NewSlackNotifier(list account.List) *SlackNotifier {
	return &SlackNotifier{
		AccountList: list,
	}
}

type PostMessageRequest struct {
	Text      string `json:"text,omitempty"`
	LinkNames string `json:"link_names,omitempty"`
}

func (n *SlackNotifier) Notify(webhookURL, message string) error {

	message, ok := n.toMentionCommentBody(message)
	if !ok {
		return nil
	}

	pm := PostMessageRequest{
		Text: message,
		LinkNames: "1",
	}

	body, err := json.Marshal(pm)
	log.Printf("hoge %s", webhookURL)
	req, _ := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(b))
	}
	return nil
}


func (n *SlackNotifier) toMention(slackName string) string {
	name := n.AccountList.Accounts[slackName]

	if name == "" {
		name = slackName
	}

	return fmt.Sprintf("<@%s>", name)
}

func (n *SlackNotifier) ConvertMessage(message string) (ok bool, convertMessage string) {
	return false, message
}

// ReplaceComment replace github account to slack
func (n *SlackNotifier) toMentionCommentBody(comment string) (string, bool) {
	return comment, true
	matches := r.FindAllStringSubmatch(comment, -1)
	if len(matches) == 0 {
		return "", false
	}
	for _, val := range matches {
		slackName, _ := n.AccountList.Accounts[val[0]]
		log.Printf("slackName %v", slackName)
		comment = strings.Replace(comment, val[0], slackName, -1)
	}
	return comment, true
}
