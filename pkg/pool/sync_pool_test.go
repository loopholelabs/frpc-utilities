package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type demoData struct {
	Test string
}

func (d *demoData) Reset() {
	d.Test = ""
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
