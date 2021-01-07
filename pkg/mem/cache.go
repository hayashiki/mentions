package mem

import (
	"encoding/json"
	"github.com/hayashiki/mentions/pkg/slack"
	"github.com/memcachier/mc"
	"github.com/pkg/errors"
	"log"
	"time"
)

const (
	prefixKey = "mentionsCommentCache"
)

type commentCache struct {
	memcached *mc.Client
}

type Quit func()

type Comment struct {
	Body string
}

type slackMessageCache slack.MessageResponse

func NewCommentCache(config *MemcachedConfig) (*commentCache, Quit) {
	memcached := mc.NewMC(config.Server, config.Username, config.Password)

	return &commentCache{memcached: memcached}, memcached.Quit
}

func (c *commentCache) Get(key string) (*slack.MessageResponse, error) {
	value, _, _, err := c.memcached.Get(prefixKey + key)

	if err != nil {
		if err == mc.ErrNotFound {
			log.Printf("not found")
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	var comment slack.MessageResponse
	err = json.Unmarshal([]byte(value), &comment)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Printf("pass??")
	return &comment, nil
}

func (c *commentCache) Set(key string, comment *slack.MessageResponse) error {
	bytes, err := json.Marshal(comment)

	if err != nil {
		return errors.WithStack(err)
	}

	expiration := time.Hour * 24 // 1 day
	_, err = c.memcached.Set(prefixKey+key, string(bytes), 0, uint32(expiration.Seconds()), 0)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
