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

type Lock struct {
	ch chan struct{}
	mu sync.RWMutex
}

type HashLock[T comparable] struct {
	locks   map[T]*Lock
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
		locks:   make(map[T]*Lock),
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
	lock := l.get(key)
	lock.mu.RLock()
	lock.ch <- struct{}{}
	lock.mu.RUnlock()
	if l.timeout > 0 {
		time.AfterFunc(l.timeout, func() {
			l.Unlock(key)
		})
	}
}

func (l *HashLock[T]) Unlock(key T) {
	lock := l.get(key)
	lock.mu.RLock()
	select {
	case <-lock.ch:
	default:
	}
	lock.mu.RUnlock()
}

func (l *HashLock[T]) get(key T) *Lock {
	l.mu.Lock()
	lock, found := l.locks[key]
	if !found {
		lock = &Lock{ch: make(chan struct{}, 1)}
		l.locks[key] = lock
	}
	l.mu.Unlock()
	return lock
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
				if v.mu.TryLock() {
					if len(v.ch) == 0 {
						delete(l.locks, k)
					}
					v.mu.Unlock()
				}
			}
			l.mu.Unlock()
		}
	}
}
