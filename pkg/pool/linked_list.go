package pool

import "sync"

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

type Node[T any] struct {
	prev *Node[T]
	next *Node[T]
	key  T
}

type LinkedList[T any] struct {
	lock sync.Mutex
	head *Node[T]
	// tail *Node[T]
	len uint64
}

func (l *LinkedList[T]) Len() uint64 {
	l.lock.Lock()
	defer l.lock.Unlock()

	return l.len
}

func (l *LinkedList[T]) Insert(key T) *Node[T] {
	l.lock.Lock()
	defer l.lock.Unlock()

	newNode := &Node[T]{
		key:  key,
		next: l.head,
	}

	if l.head != nil {
		l.head.prev = newNode
	}
	l.head = newNode

	l.len++

	return newNode
}

// TODO: Make private and only use in list
func (l *LinkedList[T]) ToArray() []T {
	l.lock.Lock()
	defer l.lock.Unlock()

	out := []T{}

	el := l.head
	for el != nil {
		out = append(out, el.key)

		el = el.next
	}

	return out
}
