package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type demoData struct {
	Test string
}

func TestPoolIntegration(t *testing.T) {
	tests := []struct {
		name   string
		before func(*Pool[demoData, *demoData]) *demoData
		check  func(*Pool[demoData, *demoData], *demoData)
	}{
		{
			name: "Can allocate the contained value",
			before: func(p *Pool[demoData, *demoData]) *demoData {
				return p.New()
			},
			check: func(p *Pool[demoData, *demoData], dd *demoData) {
				assert.NotNil(t, dd)
			},
		},
		{
			name: "Can put and get a non-nil value",
			before: func(p *Pool[demoData, *demoData]) *demoData {
				v := p.New()
				v.Test = "Testing"

				p.Put(v)

				return v
			},
			check: func(p *Pool[demoData, *demoData], dd *demoData) {
				assert.Equal(t, "Testing", p.Get().Test)
			},
		},
		{
			name: "Can put and get a nil value",
			before: func(p *Pool[demoData, *demoData]) *demoData {
				p.Put(nil)

				return nil
			},
			check: func(p *Pool[demoData, *demoData], dd *demoData) {
				assert.Nil(t, p.Get())
			},
		},
		{
			name: "Can put, change and get a value",
			before: func(p *Pool[demoData, *demoData]) *demoData {
				v := p.New()
				v.Test = "Testing"

				p.Put(v)

				assert.Equal(t, "Testing", p.Get().Test)

				v2 := p.New()
				v2.Test = "Testing 2"

				p.Put(v2)

				return v2
			},
			check: func(p *Pool[demoData, *demoData], dd *demoData) {
				assert.Equal(t, "Testing 2", p.Get().Test)
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
