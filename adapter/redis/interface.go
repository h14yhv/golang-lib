package redis

import (
	"time"

	"github.com/h14yhv/golang-lib/clock"
)

type (
	Service interface {
		// Common
		Ping() error
		Delete(keys ...string) error
		Expire(key string, ttl clock.Duration) error
		ExpireAt(key string, tm time.Time) error
		// String
		Set(key, value string, ttl clock.Duration) error
		SetObject(key string, value interface{}, ttl clock.Duration) error
		Get(key string) (string, error)
		GetObject(key string, pointer interface{}) error
	}
)
