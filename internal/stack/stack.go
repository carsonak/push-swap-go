package stack

import (
	"iter"
	"slices"
)

// Stack represents a Last-In-First-Out (LIFO) data structure.
type Stack[T any] struct {
	// The last item in the slice will be the top of the stack.
	items []T
}

// New creates a `Stack`. You can optionally pass initialising elements.
// The first element passed will be at the top of the stack,
// and the last element will be at the bottom.
func New[T any](items ...T) *Stack[T] {
	// Create a copy of the data to own it.
	cpy := make([]T, len(items))

	copy(cpy, items)
	slices.Reverse(cpy)
	return &Stack[T]{
		items: cpy,
	}
}

// NewWithCapacity creates an empty stack pre-allocated to hold at least
// `capacity` elements without triggering memory allocation on push.
func NewWithCapacity[T any](capacity int) *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0, capacity),
	}
}

// Push adds a new element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the element at the top of the stack.
func (s *Stack[T]) Pop() (T, bool) {
	var topItem T

	if len(s.items) < 1 {
		return topItem, false
	}

	topIndex := len(s.items) - 1

	topItem = s.items[topIndex]
	// Zero out the memory slot to drop references.
	s.items[topIndex] = *new(T)
	s.items = s.items[:topIndex]
	return topItem, true
}

// Index returns the value at the specified index where 0 is the top of the stack.
func (s *Stack[T]) Index(index int) (T, bool) {
	if index < 0 || index >= len(s.items) {
		var zero T
		return zero, false
	}

	// Since the end of the slice is the top of the stack we
	// map the user's index (0 = top) to the internal slice index
	internalIndex := len(s.items) - 1 - index
	return s.items[internalIndex], true
}

// Rotate shifts all elements up. The top element becomes the bottom element.
func (s *Stack[T]) Rotate() {
	if len(s.items) <= 1 {
		return
	}

	top := s.items[len(s.items)-1]

	// Shift everything else one slot up the stack.
	copy(s.items[1:], s.items[:len(s.items)-1])
	s.items[0] = top
}

// ReverseRotate shifts all elements down. The bottom element becomes the top element.
func (s *Stack[T]) ReverseRotate() {
	if len(s.items) <= 1 {
		return
	}

	bottom := s.items[0]

	// Shift everything else one slot down the stack.
	copy(s.items[:len(s.items)-1], s.items[1:])
	s.items[len(s.items)-1] = bottom
}

// All returns an iterator over the stack elements from top (index 0) to bottom.
// This allows you to use a standard `for index, value := range stack.All()` loop.
func (s *Stack[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range s.Len() {
			val, _ := s.Index(i)
			if !yield(i, val) {
				return
			}
		}
	}
}

// Len returns the current number of items in the stack.
func (s *Stack[T]) Len() int {
	return len(s.items)
}
