package stack

import (
	"iter"
)

// Stack represents a Last-In-First-Out (LIFO) data structure.
// Implemented as a ring buffer to achieve O(1) rotations.
type Stack[T any] struct {
	items []T
	head  int // index of the top element
	size  int // current number of elements
}

// internalIndex calculates the actual slice index for a given logical stack index.
func (s *Stack[T]) internalIndex(i int) int {
	return (s.head + i) % len(s.items)
}

// New creates a `Stack`. You can optionally pass initialising elements.
// The first element passed will be at the top of the stack.
func New[T any](items ...T) *Stack[T] {
	cpy := make([]T, len(items))

	copy(cpy, items)
	return &Stack[T]{
		items: cpy,
		head:  0,
		size:  len(items),
	}
}

// NewWithCapacity creates an empty stack pre-allocated to hold at least
// `capacity` elements without triggering memory allocation on push.
func NewWithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		items: make([]T, capacity),
		head:  0,
		size:  0,
	}
}

// PushBottom adds a new element to the bottom of the stack.
// This is a helper to initialize the stack with items in correct order.
func (s *Stack[T]) PushBottom(v T) {
	if s.size == len(s.items) {
		s.resize()
	}

	if s.size == 0 {
		s.items[s.head] = v
	} else {
		bottomIdx := (s.head + s.size) % len(s.items)
		s.items[bottomIdx] = v
	}

	s.size++
}

// resize doubles the capacity of the underlying array when full.
func (s *Stack[T]) resize() {
	newCap := len(s.items) * 2
	if newCap == 0 {
		newCap = 8 // Default initial capacity
	}

	newItems := make([]T, newCap)

	for i := range s.size {
		newItems[i] = s.items[s.internalIndex(i)]
	}

	s.items = newItems
	s.head = 0
}

// Swap switches the positions of the first 2 elements in the stack.
func (s *Stack[T]) Swap() {
	if s.size < 2 {
		return
	}

	idx0 := s.head
	idx1 := s.internalIndex(1)

	s.items[idx0], s.items[idx1] = s.items[idx1], s.items[idx0]
}

// Push adds a new element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	if s.size == len(s.items) {
		s.resize()
	}

	// Move head backwards (wrapping around) to prepend the new top.
	s.head = (s.head - 1 + len(s.items)) % len(s.items)
	s.items[s.head] = v
	s.size++
}

// Pop removes and returns the element at the top of the stack.
func (s *Stack[T]) Pop() (T, bool) {
	if s.size == 0 {
		var zero T

		return zero, false
	}

	topItem := s.items[s.head]

	// Zero out memory to drop references
	s.items[s.head] = *new(T)

	s.head = s.internalIndex(1)
	s.size--

	return topItem, true
}

// Index returns the value at the specified index where 0 is the top of the stack.
func (s *Stack[T]) Index(index int) (T, bool) {
	if index < 0 || index >= s.size {
		var zero T

		return zero, false
	}

	return s.items[s.internalIndex(index)], true
}

// Rotate shifts all elements up. The top element becomes the bottom element.
func (s *Stack[T]) Rotate() {
	if s.size < 2 {
		return
	}

	// Copy the top element into the slot just after the current bottom so it
	// remains reachable once head advances. When the buffer is full the copy
	// lands on head itself (a no-op), which is also correct.
	newBottomIdx := (s.internalIndex(s.size-1) + 1) % len(s.items)
	s.items[newBottomIdx] = s.items[s.head]
	s.head = s.internalIndex(1)
}

// ReverseRotate shifts all elements down. The bottom element becomes the top element.
func (s *Stack[T]) ReverseRotate() {
	if s.size < 2 {
		return
	}
	// Copy the bottom element into the slot just before head so it becomes
	// the new top once head retreats. When the buffer is full the copy
	// lands on the bottom index itself (a no-op), which is also correct.
	newHeadIdx := (s.head - 1 + len(s.items)) % len(s.items)
	s.items[newHeadIdx] = s.items[s.internalIndex(s.size-1)]
	s.head = newHeadIdx
}

// All returns an iterator over the stack elements from top (index 0) to bottom.
func (s *Stack[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range s.size {
			if !yield(i, s.items[s.internalIndex(i)]) {
				return
			}
		}
	}
}

// Len returns the current number of items in the stack.
func (s *Stack[T]) Len() int {
	return s.size
}
