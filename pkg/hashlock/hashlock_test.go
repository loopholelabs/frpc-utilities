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
