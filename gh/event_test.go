package gh

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/hayashiki/mentions/account"
	"github.com/hayashiki/mentions/config"
	mock_gh "github.com/hayashiki/mentions/gh/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

//func TestWebhook_ParsePayloadHandler(t *testing.T) {
//	g := gin.New()
//	//buf := new(strings.Builder)
//
//	req := httptest.NewRequest(http.MethodPost, "/", nil)
//	rec := httptest.NewRecorder()
//
//	hook := NewWebhookHandler()
//
//	g.POST("/", hook.githubEventHandler)
//	g.ServeHTTP(rec, req)
//
//	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
//}

func TestFindWebhookURL(t *testing.T) {
	var list account.List
	var event *Event
	event.Repository = "hogep"
	list.Repos = map[string]string{"hayashiki/gcp-functions-fw":"https://hooks.slack.com/services/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"}

	url, err := FindWebhookURL(event, list)
	assert.Equal(t, url, "")
	assert.Equal(t, err, "")
}

func TestWebhook_ParsePayloadHandler2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//g := gin.New()
	gin.SetMode(gin.TestMode)
	requestJSON, _ := ioutil.ReadFile(filepath.Join("testdata", "githubIssueCommentEvent.json"))

	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestJSON))
	rec := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "issue_comment")
	ctx, _ := gin.CreateTestContext(rec)
	ctx.Request = req

	mockVefifier := mock_gh.NewMockVerifier(ctrl)
	mockVefifier.EXPECT().Verify(req, []byte("dummy")).Return([]byte(string(requestJSON)), nil)

	var env config.Environment
	env.GithubWebhookSecret = "dummy"
	hook := NewWebhookHandler(mockVefifier, env)

	hook.PostWebhook(ctx)
	assert.Equal(t, http.StatusOK, rec.Code)
}
