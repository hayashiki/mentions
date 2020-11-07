package slack

import "fmt"

type ConvertPayload struct {
	Comment  string
	RepoName string
	HTMLURL  string
	Title    string
	User     string
}

func (p ConvertPayload) buildMessage() string {
	var text string
	text = fmt.Sprintf("*%v <%v|%v> * by: %v", p.RepoName, p.HTMLURL, p.Title, p.User)
	text = fmt.Sprintf("%v\n%v", text, p.Comment)
	return text
}
