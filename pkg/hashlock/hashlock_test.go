// SPDX-License-Identifier: Apache-2.0

package hashlock

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const iterations = 10

func TestLockUnlock(t *testing.T) {
	h := New[string](0)
	t.Cleanup(func() { h.Close() })

	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(i int) {
			time.Sleep(time.Duration(iterations-i) * time.Millisecond)
			h.Lock(t.Name())
			time.Sleep(time.Duration(iterations-i) * time.Millisecond)
			h.Unlock(t.Name())
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestDoubleUnlock(t *testing.T) {
	h := New[string](0)
	t.Cleanup(func() { h.Close() })
	h.Lock(t.Name())
	h.Unlock(t.Name())
	h.Unlock(t.Name())
}

func TestUnlockWithoutLock(t *testing.T) {
	h := New[string](0)
	t.Cleanup(func() { h.Close() })
	h.Unlock(t.Name())
}

func TestTimeout(t *testing.T) {
	h := New[string](DefaultTimeout)
	t.Cleanup(func() { h.Close() })
	startTime := time.Now()
	h.Lock(t.Name())
	h.Lock(t.Name())
	endTime := time.Now()
	require.GreaterOrEqual(t, endTime.Sub(startTime), DefaultTimeout)
}

func TestRelock(t *testing.T) {
	h := New[string](DefaultTimeout)
	t.Cleanup(func() { h.Close() })
	h.Lock(t.Name())
	h.Unlock(t.Name())
	startTime := time.Now()
	h.Lock(t.Name())
	endTime := time.Now()
	require.Less(t, endTime.Sub(startTime), time.Millisecond)
}

func (l *HashLock[T]) LockHacked(key T) {
	lock := l.get(key)
	time.Sleep(3 * time.Second)
	lock.mu.RLock()
	lock.ch <- struct{}{}
	lock.mu.RUnlock()
	if l.timeout > 0 {
		time.AfterFunc(l.timeout, func() {
			l.Unlock(key)
		})
	}
}

func TestDoubleLockWhenGCDuringLock(t *testing.T) {
	GCTime = time.Nanosecond
	_DefaultTimeout := 5 * time.Second

	h := New[string](_DefaultTimeout)
	t.Cleanup(func() { h.Close() })

	h.LockHacked(t.Name())

	// Verify the first lock timeout before we're able to lock again.
	start := time.Now()
	h.Lock(t.Name())
	require.GreaterOrEqual(t, time.Now().Sub(start), DefaultTimeout)
}
