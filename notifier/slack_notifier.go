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

//import (
//	"github.com/slack-go/slack"
//)

type SlackNotifier struct {
	WebhookURL  string
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

	log.Printf("calll slack %v", account.List{})
	log.Printf("calll slack %v", n.AccountList)

	pm := PostMessageRequest{
		Text: message,
		LinkNames: "1",
	}

	body, err := json.Marshal(pm)
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

//func (n *SlackNotifier) generateMessage(event *gh.Event) string {
//	var text string
//	text = fmt.Sprintf("%v *【%v】%v* \n", text, event.Repository, event.Title)
//	text = fmt.Sprintf("%v%v\n", text, event.HTMLURL)
//	text = fmt.Sprintf("%v>Comment created by: %v\n", text, event.User)
//	text = fmt.Sprintf("%v\n%v\n", text, event.Comment)
//	return text
//}

func (n *SlackNotifier) toMention(slackName string) string {
	name := n.AccountList.Accounts[slackName]

	if name == "" {
		name = slackName
	}

	return fmt.Sprintf("<@%s>", name)
}

// ReplaceComment replace github account to slack
func (n *SlackNotifier) toMentionCommentBody(comment string) (string, bool) {
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
