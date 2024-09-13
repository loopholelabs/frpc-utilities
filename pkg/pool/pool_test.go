// SPDX-License-Identifier: Apache-2.0

package pool

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type demoData struct {
	Test string
}

func (d *demoData) Reset() {
	d.Test = ""
}

type testData struct {
	Data []byte
}

func (t *testData) Reset() {
	t.Data = t.Data[:0]
}

func TestPoolIntegration(t *testing.T) {
	tests := []struct {
		name   string
		before func(*Pool[demoData, *demoData]) *demoData
		check  func(*Pool[demoData, *demoData], *demoData)
	}{
		{
			name: "Can allocate with New()",
			before: func(p *Pool[demoData, *demoData]) *demoData {
				return p.New()
			},
			check: func(p *Pool[demoData, *demoData], dd *demoData) {
				assert.NotNil(t, dd)
			},
		},
		{
			name: "Can allocate with Get() if nothing has been Put() before",
			before: func(p *Pool[demoData, *demoData]) *demoData {
				return p.Get()
			},
			check: func(p *Pool[demoData, *demoData], dd *demoData) {
				assert.NotNil(t, dd)
			},
		},
		{
			name: "Can Put() non-nil and Get() an allocated object",
			before: func(p *Pool[demoData, *demoData]) *demoData {
				d := p.Get()

				d.Test = "Testing"

				p.Put(d)

				return nil
			},
			check: func(p *Pool[demoData, *demoData], dd *demoData) {
				assert.Equal(t, "", p.Get().Test)
			},
		},
		{
			name: "Can Put() nil and Get() an allocated object",
			before: func(p *Pool[demoData, *demoData]) *demoData {
				p.Put(nil)

				return nil
			},
			check: func(p *Pool[demoData, *demoData], dd *demoData) {
				assert.Equal(t, "", p.Get().Test)
			},
		},
		{
			name: "Can Put() multiple times and Get() an allocated object",
			before: func(p *Pool[demoData, *demoData]) *demoData {
				d := p.Get()

				d.Test = "Testing"

				p.Put(d)

				d2 := p.Get()

				d2.Test = "Testing 2"

				p.Put(d2)

				return nil
			},
			check: func(p *Pool[demoData, *demoData], dd *demoData) {
				assert.Equal(t, "", p.Get().Test)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := NewPool(func() *demoData {
				return new(demoData)
			})

			tt.check(pool, tt.before(pool))
		})
	}
}

func TestPoolAllocation(t *testing.T) {
	newTestData := func() *testData {
		t := new(testData)
		t.Data = make([]byte, 0, 32)
		return t
	}

	d := newTestData()
	assert.Equal(t, 0, len(d.Data))
	assert.Equal(t, 32, cap(d.Data))

	dataPool := NewPool[testData, *testData](newTestData)

	randomData := make([]byte, 64)
	_, err := rand.Read(randomData)
	assert.NoError(t, err)

	d = dataPool.Get()
	assert.Equal(t, 0, len(d.Data))
	assert.Equal(t, 32, cap(d.Data))

	for {
		d.Data = append(d.Data, randomData...)
		assert.Equal(t, 64, len(d.Data))
		assert.Equal(t, 64, cap(d.Data))

		dataPool.Put(d)
		d = dataPool.Get()
		if cap(d.Data) == 64 {
			assert.Equal(t, 0, len(d.Data))
			assert.Equal(t, 64, cap(d.Data))
			break
		}
	}

	d = dataPool.Get()
	assert.Equal(t, 0, len(d.Data))
	assert.Equal(t, 32, cap(d.Data))

	dataPool.Put(nil)
	d2 := dataPool.Get()
	assert.Equal(t, 0, len(d2.Data))
	assert.Equal(t, 32, cap(d2.Data))

	allocs := testing.AllocsPerRun(100, func() {
		d = dataPool.Get()
		dataPool.Put(d)
	})
	assert.Equal(t, float64(0), allocs)
}
