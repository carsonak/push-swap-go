package stack

import "iter"

// node represents a single element in the doubly linked list.
type node[T any] struct {
	value T
	next  *node[T]
	prev  *node[T]
}

// Stack represents a Last-In-First-Out (LIFO) data structure.
type Stack[T any] struct {
	head *node[T]
	tail *node[T]
	size int
	free *node[T] // Pointer to our pool of available, pre-allocated nodes
}

// New creates a `Stack`. You can optionally pass initialising elements.
// The first element passed will be at the top of the stack,
// and the last element will be at the bottom.
func New[T any](items ...T) *Stack[T] {
	s := &Stack[T]{}

	for _, item := range items {
		n := &node[T]{value: item}

		if s.size == 0 {
			s.head = n
			s.tail = n
		} else {
			s.tail.next = n
			n.prev = s.tail
			s.tail = n
		}

		s.size++
	}

	return s
}

// NewWithCapacity pre-allocates a contiguous block of nodes to guarantee
// no memory allocation on push for at least `capacity` items.
func NewWithCapacity[T any](capacity int) *Stack[T] {
	s := &Stack[T]{}

	if capacity > 0 {
		pool := make([]node[T], capacity)

		for i := 0; i < capacity-1; i++ {
			pool[i].next = &pool[i+1]
		}

		s.free = &pool[0]
	}
	return s
}

// Swap switches the positions of the first 2 elements in the stack.
func (s *Stack[T]) Swap() {
	if s.size < 2 {
		return
	}

	first := s.head
	second := s.head.next

	first.next = second.next
	if second.next != nil {
		second.next.prev = first
	} else {
		s.tail = first // Stack only had 2 items, so 'first' is now the tail
	}

	second.prev = nil
	second.next = first
	first.prev = second
	s.head = second
}

func (s *Stack[T]) Push(v T) {
	var n *node[T]

	if s.free == nil {
		n = &node[T]{value: v}
	} else {
		// We have pre-allocated nodes in the free pool, so use those.
		n = s.free
		s.free = s.free.next
		n.value = v
		n.next = nil
		n.prev = nil
	}

	if s.size == 0 {
		s.head = n
		s.tail = n
	} else {
		n.next = s.head
		s.head.prev = n
		s.head = n
	}

	s.size++
}

// Pop removes and returns the element at the top of the stack.
func (s *Stack[T]) Pop() (T, bool) {
	if s.size == 0 {
		var zero T

		return zero, false
	}

	n := s.head
	s.head = n.next

	if s.head == nil {
		s.tail = nil
	} else {
		s.head.prev = nil
	}

	s.size--

	val := n.value
	n.value = *new(T) // Zero out memory to drop references.
	// Return the node to the free pool so future Pushes don't allocate.
	n.prev = nil
	n.next = s.free
	s.free = n
	return val, true
}

// Index returns the value at the specified index where 0 is the top of the stack.
func (s *Stack[T]) Index(index int) (T, bool) {
	if index < 0 || index >= s.size {
		var zero T

		return zero, false
	}

	curr := s.head
	for range index {
		curr = curr.next
	}

	return curr.value, true
}

// Rotate shifts all elements up. The top element becomes the bottom element.
func (s *Stack[T]) Rotate() {
	if s.size < 2 {
		return
	}

	oldHead := s.head

	// Detach head
	s.head = oldHead.next
	s.head.prev = nil

	// Attach oldHead to tail
	oldHead.next = nil
	oldHead.prev = s.tail

	s.tail.next = oldHead
	s.tail = oldHead
}

// ReverseRotate shifts all elements down. The bottom element becomes the top element.
func (s *Stack[T]) ReverseRotate() {
	if s.size < 2 {
		return
	}

	oldTail := s.tail

	// Detach tail
	s.tail = oldTail.prev
	s.tail.next = nil

	// Attach oldTail to head
	oldTail.prev = nil
	oldTail.next = s.head

	s.head.prev = oldTail
	s.head = oldTail
}

// All returns an iterator over the stack elements from top (index 0) to bottom.
func (s *Stack[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		curr := s.head
		i := 0

		for curr != nil {
			if !yield(i, curr.value) {
				return
			}

			curr = curr.next
			i++
		}
	}
}

// Len returns the current number of items in the stack.
func (s *Stack[T]) Len() int {
	return s.size
}
