package dllStack

import "iter"

// node represents a single element in the doubly linked list.
type node[T any] struct {
	value T
	next  int // index of the next node, -1 means nil
	prev  int // index of the previous node, -1 means nil
}

// Stack represents a Last-In-First-Out (LIFO) data structure.
type Stack[T any] struct {
	head  int
	tail  int
	size  int
	nodes []node[T] // The underlying array storing nodes
	free  []int     // Stack of available indices from popped nodes
}

// New creates a `Stack`.
func New[T any](items ...T) *Stack[T] {
	s := &Stack[T]{
		head: -1,
		tail: -1,
	}

	if len(items) > 0 {
		s.nodes = make([]node[T], 0, len(items))
	}

	for _, item := range items {
		s.PushBottom(item) // Helper function to initialize in correct order
	}

	return s
}

// NewWithCapacity pre-allocates a contiguous block of memory.
func NewWithCapacity[T any](capacity int) *Stack[T] {
	s := &Stack[T]{
		head: -1,
		tail: -1,
	}

	if capacity > 0 {
		s.nodes = make([]node[T], 0, capacity)
		s.free = make([]int, 0, capacity)
	}

	return s
}

// PushBottom is a helper to initialize the stack with items in correct order
func (s *Stack[T]) PushBottom(v T) {
	idx := len(s.nodes)

	s.nodes = append(s.nodes, node[T]{value: v, next: -1, prev: -1})
	if s.size == 0 {
		s.head = idx
		s.tail = idx
	} else {
		s.nodes[s.tail].next = idx
		s.nodes[idx].prev = s.tail
		s.tail = idx
	}
	s.size++
}

// Swap switches the positions of the first 2 elements in the stack.
func (s *Stack[T]) Swap() {
	if s.size < 2 {
		return
	}

	firstIdx := s.head
	secondIdx := s.nodes[firstIdx].next

	// Update first node's pointers
	s.nodes[firstIdx].next = s.nodes[secondIdx].next

	// Update the third node's prev pointer (if it exists)
	if s.nodes[secondIdx].next != -1 {
		s.nodes[s.nodes[secondIdx].next].prev = firstIdx
	} else {
		s.tail = firstIdx // Stack only had 2 items
	}

	// Update second node's pointers
	s.nodes[secondIdx].prev = -1
	s.nodes[secondIdx].next = firstIdx
	s.nodes[firstIdx].prev = secondIdx

	s.head = secondIdx
}

func (s *Stack[T]) Push(v T) {
	var idx int

	if len(s.free) > 0 {
		// Reuse a popped node's index
		idx = s.free[len(s.free)-1]
		s.free = s.free[:len(s.free)-1]
		s.nodes[idx] = node[T]{value: v, next: -1, prev: -1}
	} else {
		// Allocate a new node
		idx = len(s.nodes)
		s.nodes = append(s.nodes, node[T]{value: v, next: -1, prev: -1})
	}

	if s.size == 0 {
		s.head = idx
		s.tail = idx
	} else {
		s.nodes[idx].next = s.head
		s.nodes[s.head].prev = idx
		s.head = idx
	}

	s.size++
}

// Pop removes and returns the element at the top of the stack.
func (s *Stack[T]) Pop() (T, bool) {
	if s.size == 0 {
		var zero T
		return zero, false
	}

	idx := s.head
	poppedNode := &s.nodes[idx]
	val := poppedNode.value

	s.head = poppedNode.next

	if s.head == -1 {
		s.tail = -1
	} else {
		s.nodes[s.head].prev = -1
	}

	s.size--

	// Clear the node data to prevent memory leaks if T contains pointers
	*poppedNode = node[T]{next: -1, prev: -1}

	// Add the index to the free pool
	s.free = append(s.free, idx)

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
		curr = s.nodes[curr].next
	}

	return s.nodes[curr].value, true
}

// Rotate shifts all elements up. The top element becomes the bottom element.
func (s *Stack[T]) Rotate() {
	if s.size < 2 {
		return
	}

	oldHeadIdx := s.head

	// Detach head
	s.head = s.nodes[oldHeadIdx].next
	s.nodes[s.head].prev = -1

	// Attach oldHead to tail
	s.nodes[oldHeadIdx].next = -1
	s.nodes[oldHeadIdx].prev = s.tail

	s.nodes[s.tail].next = oldHeadIdx
	s.tail = oldHeadIdx
}

// ReverseRotate shifts all elements down. The bottom element becomes the top element.
func (s *Stack[T]) ReverseRotate() {
	if s.size < 2 {
		return
	}

	oldTailIdx := s.tail

	// Detach tail
	s.tail = s.nodes[oldTailIdx].prev
	s.nodes[s.tail].next = -1

	// Attach oldTail to head
	s.nodes[oldTailIdx].prev = -1
	s.nodes[oldTailIdx].next = s.head

	s.nodes[s.head].prev = oldTailIdx
	s.head = oldTailIdx
}

// All returns an iterator over the stack elements from top (index 0) to bottom.
func (s *Stack[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		curr := s.head
		i := 0

		for curr != -1 {
			if !yield(i, s.nodes[curr].value) {
				return
			}

			curr = s.nodes[curr].next
			i++
		}
	}
}

// Len returns the current number of items in the stack.
func (s *Stack[T]) Len() int {
	return s.size
}
