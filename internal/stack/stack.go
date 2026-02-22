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

// internalIndex calculates maps the given index to the correct index in the stack.
func internalIndex(i, len int) int {
	return len - 1 - i
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

// Swap switches the positions of the first 2 elements in the stack.
func (s *Stack[T]) Swap() {
	if s.Len() < 2 {
		return
	}

	index0 := internalIndex(0, s.Len())
	index1 := internalIndex(1, s.Len())

	s.items[index0], s.items[index1] = s.items[index1], s.items[index0]
}

// Push adds a new element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the element at the top of the stack.
func (s *Stack[T]) Pop() (T, bool) {
	var topItem T

	if s.Len() < 1 {
		return topItem, false
	}

	topIndex := internalIndex(0, s.Len())

	topItem = s.items[topIndex]
	// Zero out the memory slot to drop references.
	s.items[topIndex] = *new(T)
	s.items = s.items[:topIndex]
	return topItem, true
}

// Index returns the value at the specified index where 0 is the top of the stack.
func (s *Stack[T]) Index(index int) (T, bool) {
	if index < 0 || index >= s.Len() {
		var zero T
		return zero, false
	}

	// Since the end of the slice is the top of the stack we
	// map the user's index (0 = top) to the internal slice index
	index = internalIndex(index, s.Len())
	return s.items[index], true
}

// Rotate shifts all elements up. The top element becomes the bottom element.
func (s *Stack[T]) Rotate() {
	if s.Len() < 2 {
		return
	}

	topIndex := internalIndex(0, s.Len())
	bottomIndex := internalIndex(s.Len()-1, s.Len())
	top := s.items[topIndex]

	// Shift everything else one slot up the stack.
	copy(s.items[bottomIndex+1:], s.items[:topIndex])
	s.items[bottomIndex] = top
}

// ReverseRotate shifts all elements down. The bottom element becomes the top element.
func (s *Stack[T]) ReverseRotate() {
	if s.Len() < 2 {
		return
	}

	topIndex := internalIndex(0, s.Len())
	bottomIndex := internalIndex(s.Len()-1, s.Len())
	bottom := s.items[bottomIndex]

	// Shift everything else one slot down the stack.
	copy(s.items[:topIndex], s.items[bottomIndex+1:])
	s.items[topIndex] = bottom
}

// All returns an iterator over the stack elements from top (index 0) to bottom.
// This allows you to use a standard `for index, value := range stack.All()` loop.
func (s *Stack[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range s.Len() {
			idx := internalIndex(i, s.Len())
			
			if !yield(i, s.items[idx]) {
				return
			}
		}
	}
}

// Len returns the current number of items in the stack.
func (s *Stack[T]) Len() int {
	return len(s.items)
}
