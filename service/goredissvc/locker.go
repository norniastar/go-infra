package goredissvc

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/norniastar/go-infra/contract"
	"math/rand"
	"reflect"
	"time"
)

const (
	maxRetryDelayMilliSec = 250
	minRetryDelayMilliSec = 50
)

type syncLock struct {
	ctx     context.Context
	redSync *redislock.Client
}

func (s *syncLock) Lock(key string, expires time.Duration) (func() error, error) {
	locker, err := s.redSync.Obtain(s.ctx, key, expires, nil)
	if err != nil {
		return nil, err
	}
	refreshCtx, cancel := context.WithCancel(s.ctx)

	go func() {
		ticker := time.NewTicker(time.Duration(rand.Intn(maxRetryDelayMilliSec-minRetryDelayMilliSec)+minRetryDelayMilliSec) * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-refreshCtx.Done():
				return
			case <-ticker.C:
				if err := locker.Refresh(refreshCtx, expires, nil); err != nil {
					return
				}
			}
		}
	}()
	return func() error {
		defer cancel()
		err = locker.Release(s.ctx)
		return err
	}, nil
}

func (s *syncLock) WithContext(ctx context.Context) reflect.Value {
	return reflect.ValueOf(&syncLock{
		redSync: s.redSync,
		ctx:     ctx,
	})
}

func NewLock(redis contract.IRedis) *syncLock {
	//redisSync := redislock.New(redis.(*redisAdapter).getClient())
	//return &syncLock{
	//	ctx:     context.Background(),
	//	redSync: redisSync,
	//}
	return nil
}
