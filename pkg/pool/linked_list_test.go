/*
	Copyright 2022 Loophole Labs

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		   http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package pool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	tests := []struct {
		name   string
		expect uint64
		before func(*DoubleLinkedList[string])
	}{
		{
			name:   "Works with no elements",
			expect: 0,
			before: func(dll *DoubleLinkedList[string]) {},
		},
		{
			name:   "Works with one element",
			expect: 1,
			before: func(dll *DoubleLinkedList[string]) {
				dll.Insert("Test")
			},
		},
		{
			name:   "Works with 100 elements",
			expect: 100,
			before: func(dll *DoubleLinkedList[string]) {
				for i := 0; i < 100; i++ {
					dll.Insert("Test")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDoubleLinkedList[string]()

			tt.before(list)

			if rv := list.Len(); rv != tt.expect {
				t.Errorf("DoubleLinkedList.Len() = %v, want %v", rv, tt.expect)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		name   string
		before func(*DoubleLinkedList[string])
		check  func(*DoubleLinkedList[string])
	}{
		{
			name:   "Works with no inserts",
			before: func(dll *DoubleLinkedList[string]) {},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(0))

				assert.Equal(t, dll.toArray(), []string{})
			},
		},
		{
			name: "Works with 1 insert",
			before: func(dll *DoubleLinkedList[string]) {
				dll.Insert("One")
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(1))

				assert.Equal(t, dll.toArray(), []string{"One"})
			},
		},
		{
			name: "Works with 2 inserts",
			before: func(dll *DoubleLinkedList[string]) {
				dll.Insert("One")
				dll.Insert("Two")
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(2))

				assert.Equal(t, dll.toArray(), []string{"Two", "One"})
			},
		},
		{
			name: "Works with 100 inserts",
			before: func(dll *DoubleLinkedList[string]) {
				for i := 0; i < 100; i++ {
					dll.Insert("Test")
				}
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(100))

				expected := []string{}
				for i := 0; i < 100; i++ {
					expected = append(expected, "Test")
				}

				assert.Equal(t, dll.toArray(), expected)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDoubleLinkedList[string]()

			tt.before(list)

			tt.check(list)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name   string
		before func(*DoubleLinkedList[string]) []*Node[string]
		apply  func(*DoubleLinkedList[string], []*Node[string])
		check  func(*DoubleLinkedList[string])
	}{
		{
			name: "Works with no inserts",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := []*Node[string]{}

				return nodes
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(0))

				assert.Equal(t, dll.toArray(), []string{})
			},
		},
		{
			name: "Works with 1 insert",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := append([]*Node[string]{}, dll.Insert("One"))

				return nodes
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(0))

				assert.Equal(t, dll.toArray(), []string{})
			},
		},
		{
			name: "Works with 2 inserts",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := append([]*Node[string]{}, dll.Insert("One"), dll.Insert("Two"))

				return nodes
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(0))

				assert.Equal(t, dll.toArray(), []string{})
			},
		},
		{
			name: "Works with 100 inserts",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := []*Node[string]{}

				for i := 0; i < 100; i++ {
					nodes = append(nodes, dll.Insert("Test"))
				}

				return nodes
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(0))

				assert.Equal(t, dll.toArray(), []string{})
			},
		},
		{
			name: "Can delete the head",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := append([]*Node[string]{}, dll.Insert("One"), dll.Insert("Two"))

				return nodes[0:1]
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(1))

				assert.Equal(t, dll.toArray(), []string{"Two"})
			},
		},
		{
			name: "Can delete the tail",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := append([]*Node[string]{}, dll.Insert("One"), dll.Insert("Two"))

				return nodes[1:]
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(1))

				assert.Equal(t, dll.toArray(), []string{"One"})
			},
		},
		{
			name: "Can delete the same node multiple times",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := append([]*Node[string]{}, dll.Insert("One"), dll.Insert("Two"))

				return append([]*Node[string]{nodes[1]}, nodes[1], nodes[1])
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(1))

				assert.Equal(t, dll.toArray(), []string{"One"})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDoubleLinkedList[string]()

			tt.apply(list, tt.before(list))

			tt.check(list)
		})
	}
}

func TestPopFirst(t *testing.T) {
	tests := []struct {
		name    string
		before  func(*DoubleLinkedList[string])
		shiftBy int
		check   func(*DoubleLinkedList[string])
	}{
		{
			name:    "Works with no inserts",
			before:  func(dll *DoubleLinkedList[string]) {},
			shiftBy: 0,
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(0))

				assert.Equal(t, dll.toArray(), []string{})
			},
		},
		{
			name:    "Works with no inserts and 5 shifts",
			before:  func(dll *DoubleLinkedList[string]) {},
			shiftBy: 5,
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(0))

				assert.Equal(t, dll.toArray(), []string{})
			},
		},
		{
			name: "Works with 1 inserts and 1 shift",
			before: func(dll *DoubleLinkedList[string]) {
				dll.Insert("One")
			},
			shiftBy: 1,
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(0))

				assert.Equal(t, dll.toArray(), []string{})
			},
		},
		{
			name: "Works with 3 inserts and 1 shift",
			before: func(dll *DoubleLinkedList[string]) {
				dll.Insert("One")
				dll.Insert("Two")
				dll.Insert("Three")
			},
			shiftBy: 1,
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(2))

				assert.Equal(t, dll.toArray(), []string{"Two", "One"})
			},
		},
		{
			name: "Works with 3 inserts and 3 shifts",
			before: func(dll *DoubleLinkedList[string]) {
				dll.Insert("One")
				dll.Insert("Two")
				dll.Insert("Three")
			},
			shiftBy: 3,
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(0))

				assert.Equal(t, dll.toArray(), []string{})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDoubleLinkedList[string]()

			tt.before(list)

			for i := 0; i < tt.shiftBy; i++ {
				list.PopFirst()
			}

			tt.check(list)
		})
	}
}

func TestDoubleLinkedListIntegration(t *testing.T) {
	tests := []struct {
		name   string
		before func(*DoubleLinkedList[string]) []*Node[string]
		apply  func(*DoubleLinkedList[string], []*Node[string])
		check  func(*DoubleLinkedList[string])
	}{
		{
			name: "Can delete something at the end and then insert again",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := append([]*Node[string]{}, dll.Insert("One"), dll.Insert("Two"))

				return nodes[1:]
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}

				dll.Insert("New")
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(2))

				assert.Equal(t, dll.toArray(), []string{"New", "One"})
			},
		},
		{
			name: "Can delete everything then insert again",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := append([]*Node[string]{}, dll.Insert("One"), dll.Insert("Two"))

				return nodes
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}

				dll.Insert("New")
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(1))

				assert.Equal(t, dll.toArray(), []string{"New"})
			},
		},
		{
			name: "Can shift then insert",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := append([]*Node[string]{}, dll.Insert("One"), dll.Insert("Two"), dll.Insert("Three"))

				return nodes
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				dll.PopFirst()
				dll.PopFirst()

				dll.Insert("New")
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(2))

				assert.Equal(t, dll.toArray(), []string{"New", "One"})
			},
		},
		{
			name: "Can insert one, then delete, then insert",
			before: func(dll *DoubleLinkedList[string]) []*Node[string] {
				nodes := append([]*Node[string]{}, dll.Insert("One"))

				return nodes
			},
			apply: func(dll *DoubleLinkedList[string], n []*Node[string]) {
				for _, node := range n {
					dll.Delete(node)
				}

				dll.Insert("New")
			},
			check: func(dll *DoubleLinkedList[string]) {
				assert.Equal(t, dll.Len(), uint64(1))

				assert.Equal(t, dll.toArray(), []string{"New"})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDoubleLinkedList[string]()

			tt.apply(list, tt.before(list))

			tt.check(list)
		})
	}
}
