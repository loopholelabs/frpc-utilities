// SPDX-License-Identifier: Apache-2.0

package hashlock

import (
	"context"
	"sync"
	"time"
)

const (
	// DefaultTimeout is the default timeout for locks
	DefaultTimeout = time.Second

	// GCTime is the time between garbage collection runs
	GCTime = time.Minute
)

type HashLock[T comparable] struct {
	locks   map[T]chan struct{}
	mu      sync.Mutex
	timeout time.Duration

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New[T comparable](d time.Duration) *HashLock[T] {
	if d < 0 {
		d = DefaultTimeout
	}

	ctx, cancel := context.WithCancel(context.Background())
	h := &HashLock[T]{
		locks:   make(map[T]chan struct{}),
		timeout: d,
		ctx:     ctx,
		cancel:  cancel,
	}

	h.wg.Add(1)
	go h.gc()

	return h
}

func (l *HashLock[T]) Close() {
	l.cancel()
	l.wg.Wait()
}

func (l *HashLock[T]) Lock(key T) {
	l.get(key) <- struct{}{}
	if l.timeout > 0 {
		time.AfterFunc(l.timeout, func() {
			l.Unlock(key)
		})
	}
}

func (l *HashLock[T]) Unlock(key T) {
	select {
	case <-l.get(key):
	default:
	}
}

func (l *HashLock[T]) DeleteKey(key T) {
	l.mu.Lock()
	delete(l.locks, key)
	l.mu.Unlock()
}

func (l *HashLock[T]) get(key T) chan struct{} {
	l.mu.Lock()
	lockCh, found := l.locks[key]
	if !found {
		lockCh = make(chan struct{}, 1)
		l.locks[key] = lockCh
	}
	l.mu.Unlock()
	return lockCh
}

func (l *HashLock[T]) gc() {
	for {
		select {
		case <-l.ctx.Done():
			l.wg.Done()
			return
		case <-time.After(GCTime):
			l.mu.Lock()
			for k, v := range l.locks {
				if len(v) == 0 {
					delete(l.locks, k)
				}
			}
			l.mu.Unlock()
		}
	}
}
