package notifier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//import (
//	"github.com/slack-go/slack"
//)

type SlackNotifier struct {
	WebhookURL string
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
	}
}

type PostMessageRequest struct {
	Text        string        `json:"text,omitempty"`
	LinkNames   string        `json:"link_names,omitempty"`
}

func (n *SlackNotifier) Notify(payload *PostMessageRequest) error {

	body, err := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, n.WebhookURL, bytes.NewReader(body))
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

func (n *SlackNotifier) GenerateMessage(repository, title, url, user, comment string) string {
	var text string
	text = fmt.Sprintf("%v *【%v】%v* \n", text, repository, title)
	text = fmt.Sprintf("%v%v\n", text, url)
	text = fmt.Sprintf("%v>Comment created by: %v\n", text, user)
	text = fmt.Sprintf("%v\n%v\n", text, comment)
	return text
}
