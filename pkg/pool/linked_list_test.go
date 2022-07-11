package pool

import "testing"

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
