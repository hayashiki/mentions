package memcache

import (
	"github.com/hayashiki/mentions/pkg/slack"
	"time"
)

const (
	prefixKey = "mentionsCommentCache"
)

type commentCache struct {
	client *client
}

func NewCommentCache(client *client) (*commentCache, Quit) {
	return &commentCache{
		client: client,
	}, client.quit
}

var (
	commentExp = uint32((time.Hour * 24 * 14).Seconds())
)

func (c *commentCache) Get(key string) (*slack.MessageResponse, error) {
	var comment slack.MessageResponse
	err := c.client.GetInterface(prefixKey+key, &comment)
	return &comment, err
}

func (c *commentCache) Set(key string, comment *slack.MessageResponse) error {
	return c.client.SetInterface(prefixKey+key, comment, commentExp)
}
