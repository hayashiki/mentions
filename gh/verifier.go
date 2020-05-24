package gh

import (
	"errors"
	"github.com/google/go-github/github"
	"net/http"
)

//go:generate mockgen -source verifier.go -destination mocks/verifier.go
type Verifier interface {
	Verify(r *http.Request, secret []byte) ([]byte, error)
}

type GithubVerifier struct{}

func NewGithubVerifier() *GithubVerifier {
	return &GithubVerifier{}
}

func (v *GithubVerifier) Verify(r *http.Request, secret []byte) ([]byte, error) {
	payload, err := github.ValidatePayload(r, secret)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

//gomock使う程でもないので、VerifierをみたすMockGithubVerifierを用意
type MockGithubVerifier struct {
	Valid bool
}

func (v *MockGithubVerifier) Verify(r *http.Request, secret []byte) ([]byte, error) {
	if v.Valid {
		return nil, nil
	}
	return nil, errors.New("Invalid Secrets")
}
