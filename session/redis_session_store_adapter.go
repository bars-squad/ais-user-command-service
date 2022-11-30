package session

import (
	"context"
	"fmt"
	"time"

	rv8 "github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

const ()

// RedisSessionStoreAdapter is a concrete struct of redis session store adapter.
type RedisSessionStoreAdapter struct {
	logger    *logrus.Logger
	maxAge    time.Duration
	c         rv8.UniversalClient
	keyPrefix string
}

// NewRedisSessionStoreAdapter is a constructor.
func NewRedisSessionStoreAdapter(rdb rv8.UniversalClient, maxAge time.Duration, keyPrefix string) Session {
	return RedisSessionStoreAdapter{
		logger:    logrus.New(),
		maxAge:    maxAge,
		c:         rdb,
		keyPrefix: keyPrefix,
	}
}

// Set will store the key and value as session.
func (s RedisSessionStoreAdapter) Set(ctx context.Context, key string, value []byte) (err error) {
	// span, ctx := apm.StartSpan(ctx, "Redis Session Store: Set", "cache.session")
	// defer span.End()

	key = fmt.Sprintf("%s.%s", s.keyPrefix, key)
	_, err = s.c.SetEX(ctx, key, value, s.maxAge).Result()
	if err != nil {
		return ErrUnexpected
	}
	return
}

// Get get will get the session by the given key.
func (s RedisSessionStoreAdapter) Get(ctx context.Context, key string) (value []byte, err error) {
	// span, ctx := apm.StartSpan(ctx, "Redis Session Store: Get", "cache.session")
	// defer span.End()

	key = fmt.Sprintf("%s.%s", s.keyPrefix, key)
	value, err = s.c.Get(ctx, key).Bytes()
	if err != nil {
		if err == rv8.Nil {
			return value, ErrSessionNotFound
		}

		return value, ErrUnexpected
	}

	return
}

// Update will update the session with but never change the time to live.
func (s RedisSessionStoreAdapter) Update(ctx context.Context, key string, value []byte) (err error) {
	// span, ctx := apm.StartSpan(ctx, "Redis Session Store: Update", "cache.session")
	// defer span.End()

	watchTxID := fmt.Sprintf("watch:transaction:session:update:%s", key)

	wrappedKey := fmt.Sprintf("%s.%s", s.keyPrefix, key)

	err = s.c.Watch(ctx, func(tx *rv8.Tx) (err error) {
		duration, err := tx.TTL(ctx, wrappedKey).Result()
		if err != nil {
			s.logger.Error(err)
			return ErrUnexpected
		}

		_, err = tx.TxPipelined(ctx, func(pipe rv8.Pipeliner) (err error) {
			_, err = pipe.SetEX(ctx, wrappedKey, value, duration).Result()
			return
		})

		if err != nil {
			s.logger.Error(err)
			return ErrUnexpected
		}

		return
	}, watchTxID)

	return
}

func (s RedisSessionStoreAdapter) Delete(ctx context.Context, key string) (err error) {
	// span, ctx := apm.StartSpan(ctx, "Redis Session Store: Delete", "cache.session")
	// defer span.End()

	key = fmt.Sprintf("%s.%s", s.keyPrefix, key)
	err = s.c.Del(ctx, key).Err()
	if err != nil {
		if err == rv8.Nil {
			return ErrSessionNotFound
		}

		return ErrUnexpected
	}

	return
}
