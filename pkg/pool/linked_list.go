package pool

import "sync"

func NewDoubleLinkedList[T any]() *DoubleLinkedList[T] {
	return &DoubleLinkedList[T]{}
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

func (l *DoubleLinkedList[T]) Len() uint64 {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.len
}

func (l *DoubleLinkedList[T]) Insert(key T) *Node[T] {
	l.lock.Lock()
	defer l.lock.Unlock()

	newNode := &Node[T]{
		value: key,
		next:  l.head,
	}

	if l.head != nil {
		l.head.prev = newNode
	}
	l.head = newNode

	l.len++

	return newNode
}

func (l *DoubleLinkedList[T]) Delete(node *Node[T]) {
	l.lock.Lock()
	defer l.lock.Unlock()

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
}

func (l *DoubleLinkedList[T]) Shift() *Node[T] {
	l.lock.Lock()
	defer l.lock.Unlock()

	var head Node[T]
	if l.head != nil {
		head = *l.head

		l.head = l.head.next

		l.len--
	}

	return &head
}

func (l *DoubleLinkedList[T]) toArray() []T {
	l.lock.Lock()
	defer l.lock.Unlock()

	out := []T{}

	el := l.head
	for el != nil {
		out = append(out, el.value)

		el = el.next
	}

	return out
}
