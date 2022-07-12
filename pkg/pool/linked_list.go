package pool

import "sync"

func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return new(DoubleLinkedList[T])
}

type Node[T any] struct {
	prev  *Node[T]
	next  *Node[T]
	value T
}

func (n *Node[T]) Value() T {
	return n.value
}

type DoubleLinkedList[T any] struct {
	lock sync.Mutex
	head *Node[T]
	len  uint64
}

func (l *DoubleLinkedList[T]) Len() (len uint64) {
	l.lock.Lock()

	len = l.len

	l.lock.Unlock()

	return
}

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

func (l *DoubleLinkedList[T]) Delete(node *Node[T]) {
	l.lock.Lock()

	if node == l.head {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	l.len--

	l.lock.Unlock()
}

func (l *DoubleLinkedList[T]) Shift() (node *Node[T]) {
	l.lock.Lock()

	if l.head != nil {
		node = l.head

		l.head = l.head.next

		l.len--
	}

	l.lock.Unlock()

	return
}

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
