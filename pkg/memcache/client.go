package memcache

import (
	"encoding/json"
	"fmt"
	"github.com/memcachier/mc"
	log "github.com/sirupsen/logrus"
	"sync"
)

var lock sync.Mutex

type client struct {
	memcached *mc.Client
	quit Quit
}

type Quit func()

type Item struct {
	Key string
	Value string
	Expiration uint32
}

func NewClient(server, username, password string) *client {
	mc := mc.NewMC(server, username, password)
	return &client{memcached: mc, quit: mc.Quit}
}


func (c *client) Get(key string) (val string, err error) {
	value, _, _, err := c.memcached.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			log.Printf("not found")
			return "", nil
		}
		return "", fmt.Errorf("failed to get cache")
	}
	return value, err
}

func (c *client) Set(key string, val string, exp uint32) (err error) {

	_, err = c.memcached.Set(key, val, 0, exp, 0)
	return err
}

func (c *client) SetInterface(key string, val interface{}, exp uint32) (err error) {

	b, err := json.Marshal(val)
	if err != nil {
		return err
	}

	_, err = c.memcached.Set(key, string(b), 0, exp, 0)
	return err
}

func (c *client) GetInterface(key string, i interface{}) (err error) {
	val, _, _, err := c.memcached.Get(key)
	if err != nil {
		if err == mc.ErrNotFound {
			log.Debug("not found")
			return nil
		}
		log.Debugf("get cache err %v", err)
		return err
	}

	return json.Unmarshal([]byte(val), i)
}
