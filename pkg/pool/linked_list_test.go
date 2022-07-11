package pool

import (
	"reflect"
	"testing"
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
		check  func(*DoubleLinkedList[string]) bool
	}{
		{
			name:   "Works with no inserts",
			before: func(dll *DoubleLinkedList[string]) {},
			check: func(dll *DoubleLinkedList[string]) bool {
				if dll.Len() != 0 {
					return false
				}

				return reflect.DeepEqual(dll.toArray(), []string{})
			},
		},
		{
			name: "Works with 1 insert",
			before: func(dll *DoubleLinkedList[string]) {
				dll.Insert("One")
			},
			check: func(dll *DoubleLinkedList[string]) bool {
				if dll.Len() != 1 {
					return false
				}

				return reflect.DeepEqual(dll.toArray(), []string{"One"})
			},
		},
		{
			name: "Works with 2 inserts",
			before: func(dll *DoubleLinkedList[string]) {
				dll.Insert("One")
				dll.Insert("Two")
			},
			check: func(dll *DoubleLinkedList[string]) bool {
				if dll.Len() != 2 {
					return false
				}

				return reflect.DeepEqual(dll.toArray(), []string{"Two", "One"})
			},
		},
		{
			name: "Works with 100 inserts",
			before: func(dll *DoubleLinkedList[string]) {
				for i := 0; i < 100; i++ {
					dll.Insert("Test")
				}
			},
			check: func(dll *DoubleLinkedList[string]) bool {
				if dll.Len() != 100 {
					return false
				}

				expected := []string{}
				for i := 0; i < 100; i++ {
					expected = append(expected, "Test")
				}

				return reflect.DeepEqual(dll.toArray(), expected)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewDoubleLinkedList[string]()

			tt.before(list)

			if rv := tt.check(list); rv != true {
				t.Errorf("check DoubleLinkedList.Insert() = %v, want true", rv)
			}
		})
	}
}
