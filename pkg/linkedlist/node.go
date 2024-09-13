// SPDX-License-Identifier: Apache-2.0

package linkedlist

type Pointer[T any] interface {
	*T
}

// Node is a container for data in the double linked list
type Node[T any, P Pointer[T]] struct {
	_padding0 [8]uint64 //nolint:structcheck,unused
	prev      *Node[T, P]
	_padding1 [8]uint64 //nolint:structcheck,unused
	next      *Node[T, P]
	_padding2 [8]uint64 //nolint:structcheck,unused
	value     P
}

// NewNode returns a pointer to a typed Node
func NewNode[T any, P Pointer[T]]() *Node[T, P] {
	return new(Node[T, P])
}

// Value returns the data stored in the node container
func (n *Node[T, P]) Value() P {
	return n.value
}

func (n *Node[T, P]) Reset() {
	n.prev = nil
	n.next = nil
	n.value = nil
}
