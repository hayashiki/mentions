package appcache

import (
	"time"
	"errors"
)

const (
	DefaultExpiryTime = time.Duration(0)
	ForeverNeverExpiry = time.Duration(-1)
)

var (
	ErrCacheMiss    = errors.New("cache: miss")
	ErrCASConflict  = errors.New("cache: compare-and-swap conflict")
	ErrNoStats      = errors.New("cache: no statistics available")
	ErrNotStored    = errors.New("cache: item not stored")
	ErrServerError  = errors.New("cache: server error")
	ErrInvalidValue = errors.New("cache: invalid value")
)

type Getter interface {
	Get(key string, prtValue interface{}) error
}

type Cache interface {
	Getter
	Set(key string, value interface{}, expires time.Duration) error
	SetFields(key string, value map[string]interface{}, expires time.Duration)
	GetMulti(keys ...string) (Getter, error)
	Delete(key string) error
	Add(key string, value interface{}, expires time.Duration) error
	Replace(key string, value interface{}, expires time.Duration) error
	Flush() error
	Keys() ([]string, error)
}


