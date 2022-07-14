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

package list

import "sync"

// NewDoubleLinkedList creates a new double linked list
func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return new(DoubleLinkedList[T])
}

// Node is a container for data in the double linked list
type Node[T any] struct {
	prev  *Node[T]
	next  *Node[T]
	value T
}

// Value returns the data stored in the node container
func (n *Node[T]) Value() T {
	return n.value
}

// DoubleLinkedList is a double linked list
type DoubleLinkedList[T any] struct {
	lock sync.Mutex
	head *Node[T]
	len  uint64
}

// Len returns the count of nodes stored in the double linked list
func (l *DoubleLinkedList[T]) Len() (len uint64) {
	l.lock.Lock()

	len = l.len

	l.lock.Unlock()

	return
}

// Insert adds a new node at the beginning of the double linked list
func (l *DoubleLinkedList[T]) Insert(key T) (node *Node[T]) {
	node = &Node[T]{
		value: key,
		next:  l.head,
	}

	l.lock.Lock()

	if l.head != nil {
		l.head.prev = node
	}
	l.head = node

	l.len++

	l.lock.Unlock()

	return
}

// Insert removes a node from the double linked list
func (l *DoubleLinkedList[T]) Delete(node *Node[T]) {
	l.lock.Lock()

	decrement := false
	if node == l.head {
		l.head = node.next

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

	node.prev = nil
	node.next = nil

	l.lock.Unlock()
}

// PopFirst removes and returns the first node from the double linked list
func (l *DoubleLinkedList[T]) PopFirst() (node *Node[T]) {
	l.lock.Lock()

	if l.head != nil {
		node = l.head

		l.head = l.head.next

		l.len--
	}

	l.lock.Unlock()

	return
}

// toArray is a helper functions to simplify testing
func (l *DoubleLinkedList[T]) toArray() (out []T) {
	l.lock.Lock()

	out = []T{}
	el := l.head
	for el != nil {
		out = append(out, el.value)

		el = el.next
	}

	l.lock.Unlock()

	return
}
