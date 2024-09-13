// SPDX-License-Identifier: Apache-2.0

package linkedlist

import (
	"sync"

	"github.com/loopholelabs/common/pkg/pool"
)

// NewDouble creates a new double linked list that can function as a
// FIFO queue when used with the Push and Pop methods.
func NewDouble[T any, P Pointer[T]]() *Double[T, P] {
	return &Double[T, P]{
		pool: pool.NewPool[Node[T, P], *Node[T, P]](NewNode[T, P]),
	}
}

// Double is a double-linked list
type Double[T any, P Pointer[T]] struct {
	_padding0 [8]uint64 //nolint:structcheck,unused
	lock      sync.Mutex
	_padding1 [8]uint64 //nolint:structcheck,unused
	head      *Node[T, P]
	_padding2 [8]uint64 //nolint:structcheck,unused
	tail      *Node[T, P]
	_padding3 [8]uint64 //nolint:structcheck,unused
	pool      *pool.Pool[Node[T, P], *Node[T, P]]
	_padding4 [8]uint64 //nolint:structcheck,unused
	len       uint64
}

// Length returns the count of nodes stored in the double linked list
func (l *Double[T, P]) Length() (len uint64) {
	l.lock.Lock()
	len = l.len
	l.lock.Unlock()
	return
}

// PushBack adds a new node to the end of the double linked list
func (l *Double[T, P]) PushBack(val P) (node *Node[T, P]) {
	node = l.pool.Get()
	node.value = val
	l.lock.Lock()
	node.prev = l.tail
	if l.tail != nil {
		l.tail.next = node
	}
	l.tail = node
	if l.head == nil {
		l.head = node
	}
	l.len++
	l.lock.Unlock()
	return
}

// Push adds a new node at the beginning of the double linked list
func (l *Double[T, P]) Push(val P) (node *Node[T, P]) {
	node = l.pool.Get()
	node.value = val
	l.lock.Lock()
	node.next = l.head
	if l.head != nil {
		l.head.prev = node
	}
	l.head = node
	if l.tail == nil {
		l.tail = node
	}
	l.len++
	l.lock.Unlock()
	return
}

// Delete removes a node from the double linked list
func (l *Double[T, P]) Delete(node *Node[T, P]) {
	decrement := false
	l.lock.Lock()
	if node == l.head {
		l.head = node.next
		decrement = true
	}
	if node == l.tail {
		l.tail = node.prev
		decrement = true
	}
	if node.next != nil {
		node.next.prev = node.prev
		decrement = true
	}
	if node.prev != nil {
		node.prev.next = node.next
		decrement = true
	}
	if decrement {
		l.len--
	}
	l.pool.Put(node)
	l.lock.Unlock()
}

// Pop removes and returns the node from the end of the double linked list
func (l *Double[T, P]) Pop() (val P) {
	l.lock.Lock()
	if l.tail != nil {
		node := l.tail
		l.tail = l.tail.prev
		l.len--
		val = node.Value()
		l.pool.Put(node)
	}
	l.lock.Unlock()

	return
}

// PopFront removes and returns the node from the front of the double linked list
func (l *Double[T, P]) PopFront() (val P) {
	l.lock.Lock()
	if l.head != nil {
		node := l.head
		l.head = l.head.next
		l.len--
		val = node.Value()
		l.pool.Put(node)
	}
	l.lock.Unlock()

	return
}

// toArray is a helper functions to simplify testing that converts
// the double linked list to an array ordered from head to tail
func (l *Double[T, P]) toArray() (out []P) {
	l.lock.Lock()
	out = []P{}
	el := l.tail
	for el != nil {
		out = append(out, el.Value())
		el = el.prev
	}
	l.lock.Unlock()
	return
}
