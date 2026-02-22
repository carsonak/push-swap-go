package dllStack

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		wantLen   int
		wantTop   int
		wantTopOk bool
	}{
		{
			name:      "empty stack",
			items:     []int{},
			wantLen:   0,
			wantTop:   0,
			wantTopOk: false,
		},
		{
			name:      "single element",
			items:     []int{42},
			wantLen:   1,
			wantTop:   42,
			wantTopOk: true,
		},
		{
			name:      "multiple elements",
			items:     []int{1, 2, 3, 4, 5},
			wantLen:   5,
			wantTop:   1,
			wantTopOk: true,
		},
		{
			name:      "two elements",
			items:     []int{10, 20},
			wantLen:   2,
			wantTop:   10,
			wantTopOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.items...)

			if got := s.Len(); got != tt.wantLen {
				t.Errorf("Len() = %v, want %v", got, tt.wantLen)
			}

			gotTop, gotOk := s.Pop()
			if gotOk != tt.wantTopOk {
				t.Errorf("Pop() ok = %v, want %v", gotOk, tt.wantTopOk)
			}
			if gotOk && gotTop != tt.wantTop {
				t.Errorf("Pop() value = %v, want %v", gotTop, tt.wantTop)
			}
		})
	}
}

func TestNewWithCapacity(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		wantLen  int
	}{
		{
			name:     "zero capacity",
			capacity: 0,
			wantLen:  0,
		},
		{
			name:     "small capacity",
			capacity: 5,
			wantLen:  0,
		},
		{
			name:     "large capacity",
			capacity: 100,
			wantLen:  0,
		},
		{
			name:     "capacity of 1",
			capacity: 1,
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewWithCapacity[int](tt.capacity)

			if s == nil {
				t.Fatal("NewWithCapacity returned nil")
			}

			if got := s.Len(); got != tt.wantLen {
				t.Errorf("Len() = %v, want %v", got, tt.wantLen)
			}

			// Test that we can push at least capacity elements
			for i := range tt.capacity {
				s.Push(i)
			}

			if got := s.Len(); got != tt.capacity {
				t.Errorf("After pushing %d elements, Len() = %v, want %v", tt.capacity, got, tt.capacity)
			}

			// Verify the elements are correct (top should be last pushed)
			if tt.capacity > 0 {
				topValue, ok := s.Index(0)
				if !ok {
					t.Error("Index(0) returned false after pushing elements")
				}
				wantTop := tt.capacity - 1
				if topValue != wantTop {
					t.Errorf("Top element = %v, want %v", topValue, wantTop)
				}
			}
		})
	}
}

func TestPushBottom(t *testing.T) {
	tests := []struct {
		name      string
		pushItems []int
		wantLen   int
		wantOrder []int // top to bottom
	}{
		{
			name:      "single element",
			pushItems: []int{42},
			wantLen:   1,
			wantOrder: []int{42},
		},
		{
			name:      "multiple elements in order",
			pushItems: []int{1, 2, 3},
			wantLen:   3,
			wantOrder: []int{1, 2, 3},
		},
		{
			name:      "five elements",
			pushItems: []int{10, 20, 30, 40, 50},
			wantLen:   5,
			wantOrder: []int{10, 20, 30, 40, 50},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[int]()
			for _, v := range tt.pushItems {
				s.PushBottom(v)
			}

			if s.Len() != tt.wantLen {
				t.Errorf("Len() = %v, want %v", s.Len(), tt.wantLen)
			}

			for i, want := range tt.wantOrder {
				got, ok := s.Index(i)
				if !ok {
					t.Fatalf("Index(%d) returned false", i)
				}
				if got != want {
					t.Errorf("Index(%d) = %v, want %v", i, got, want)
				}
			}
		})
	}
}

// TestPopAndReuse verifies that after popping, subsequent Push calls correctly
// reuse the freed node slots and the stack remains consistent.
func TestPopAndReuse(t *testing.T) {
	tests := []struct {
		name       string
		operations func(s *Stack[int])
		wantOrder  []int // top to bottom; nil means empty
	}{
		{
			name: "pop one then push reuses slot",
			operations: func(s *Stack[int]) {
				s.Push(1)
				s.Push(2)
				s.Pop() // frees the node that held 2
				s.Push(3)
			},
			wantOrder: []int{3, 1},
		},
		{
			name: "pop all then push again",
			operations: func(s *Stack[int]) {
				s.Push(10)
				s.Push(20)
				s.Push(30)
				s.Pop()
				s.Pop()
				s.Pop()
				s.Push(99)
			},
			wantOrder: []int{99},
		},
		{
			name: "interleaved push and pop preserves order",
			operations: func(s *Stack[int]) {
				s.Push(1)
				s.Push(2)
				s.Push(3)
				s.Pop() // removes 3
				s.Push(4)
				s.Push(5)
			},
			wantOrder: []int{5, 4, 2, 1},
		},
		{
			name: "pop all then push preserves empty between",
			operations: func(s *Stack[int]) {
				s.Push(7)
				s.Pop()
				// Stack is empty; make sure the next operation doesn't panic.
				s.Push(8)
				s.Push(9)
			},
			wantOrder: []int{9, 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[int]()
			tt.operations(s)

			if s.Len() != len(tt.wantOrder) {
				t.Errorf("Len() = %v, want %v", s.Len(), len(tt.wantOrder))
			}
			for i, want := range tt.wantOrder {
				got, ok := s.Index(i)
				if !ok {
					t.Fatalf("Index(%d) returned false", i)
				}
				if got != want {
					t.Errorf("Index(%d) = %v, want %v", i, got, want)
				}
			}
		})
	}
}

func TestSwap(t *testing.T) {
	tests := []struct {
		name          string
		initial       []int
		wantAfterSwap []int
		wantLen       int
	}{
		{
			name:          "empty stack",
			initial:       []int{},
			wantAfterSwap: []int{},
			wantLen:       0,
		},
		{
			name:          "single element",
			initial:       []int{42},
			wantAfterSwap: []int{42},
			wantLen:       1,
		},
		{
			name:          "two elements",
			initial:       []int{1, 2},
			wantAfterSwap: []int{2, 1},
			wantLen:       2,
		},
		{
			name:          "three elements",
			initial:       []int{1, 2, 3},
			wantAfterSwap: []int{2, 1, 3},
			wantLen:       3,
		},
		{
			name:          "five elements",
			initial:       []int{1, 2, 3, 4, 5},
			wantAfterSwap: []int{2, 1, 3, 4, 5},
			wantLen:       5,
		},
		{
			name:          "large stack",
			initial:       []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			wantAfterSwap: []int{20, 10, 30, 40, 50, 60, 70, 80, 90, 100},
			wantLen:       10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.initial...)
			s.Swap()

			if s.Len() != tt.wantLen {
				t.Errorf("Len() = %v, want %v", s.Len(), tt.wantLen)
			}

			for i := range len(tt.wantAfterSwap) {
				gotValue, gotOk := s.Index(i)
				if !gotOk {
					t.Fatalf("Index(%d) returned false", i)
				}
				if gotValue != tt.wantAfterSwap[i] {
					t.Errorf("Index(%d) = %v, want %v", i, gotValue, tt.wantAfterSwap[i])
				}
			}
		})
	}
}

func TestPush(t *testing.T) {
	tests := []struct {
		name         string
		initial      []int
		pushValues   []int
		wantLen      int
		wantTopValue int
	}{
		{
			name:         "push to empty stack",
			initial:      []int{},
			pushValues:   []int{42},
			wantLen:      1,
			wantTopValue: 42,
		},
		{
			name:         "push multiple to empty stack",
			initial:      []int{},
			pushValues:   []int{1, 2, 3},
			wantLen:      3,
			wantTopValue: 3,
		},
		{
			name:         "push to non-empty stack",
			initial:      []int{10, 20},
			pushValues:   []int{30},
			wantLen:      3,
			wantTopValue: 30,
		},
		{
			name:         "push multiple to non-empty stack",
			initial:      []int{1},
			pushValues:   []int{2, 3, 4, 5},
			wantLen:      5,
			wantTopValue: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.initial...)

			for _, v := range tt.pushValues {
				s.Push(v)
			}

			if got := s.Len(); got != tt.wantLen {
				t.Errorf("Len() = %v, want %v", got, tt.wantLen)
			}

			gotTop, gotOk := s.Index(0)
			if !gotOk {
				t.Fatal("Index(0) returned false, want true")
			}
			if gotTop != tt.wantTopValue {
				t.Errorf("Index(0) = %v, want %v", gotTop, tt.wantTopValue)
			}
		})
	}
}

func TestPop(t *testing.T) {
	tests := []struct {
		name       string
		initial    []int
		numPops    int
		wantValues []int
		wantOks    []bool
		wantLen    int
	}{
		{
			name:       "pop from empty stack",
			initial:    []int{},
			numPops:    1,
			wantValues: []int{0},
			wantOks:    []bool{false},
			wantLen:    0,
		},
		{
			name:       "pop single element",
			initial:    []int{42},
			numPops:    1,
			wantValues: []int{42},
			wantOks:    []bool{true},
			wantLen:    0,
		},
		{
			name:       "pop all elements",
			initial:    []int{1, 2, 3},
			numPops:    3,
			wantValues: []int{1, 2, 3},
			wantOks:    []bool{true, true, true},
			wantLen:    0,
		},
		{
			name:       "pop more than available",
			initial:    []int{10, 20},
			numPops:    3,
			wantValues: []int{10, 20, 0},
			wantOks:    []bool{true, true, false},
			wantLen:    0,
		},
		{
			name:       "pop partial",
			initial:    []int{1, 2, 3, 4, 5},
			numPops:    2,
			wantValues: []int{1, 2},
			wantOks:    []bool{true, true},
			wantLen:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.initial...)

			for i := range tt.numPops {
				gotValue, gotOk := s.Pop()

				if gotOk != tt.wantOks[i] {
					t.Errorf("Pop() %d ok = %v, want %v", i, gotOk, tt.wantOks[i])
				}
				if gotValue != tt.wantValues[i] {
					t.Errorf("Pop() %d value = %v, want %v", i, gotValue, tt.wantValues[i])
				}
			}

			if got := s.Len(); got != tt.wantLen {
				t.Errorf("Len() = %v, want %v", got, tt.wantLen)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	tests := []struct {
		name      string
		initial   []int
		index     int
		wantValue int
		wantOk    bool
	}{
		{
			name:      "empty stack",
			initial:   []int{},
			index:     0,
			wantValue: 0,
			wantOk:    false,
		},
		{
			name:      "index 0 (top)",
			initial:   []int{1, 2, 3},
			index:     0,
			wantValue: 1,
			wantOk:    true,
		},
		{
			name:      "index middle",
			initial:   []int{1, 2, 3},
			index:     1,
			wantValue: 2,
			wantOk:    true,
		},
		{
			name:      "index last (bottom)",
			initial:   []int{1, 2, 3},
			index:     2,
			wantValue: 3,
			wantOk:    true,
		},
		{
			name:      "negative index",
			initial:   []int{1, 2, 3},
			index:     -1,
			wantValue: 0,
			wantOk:    false,
		},
		{
			name:      "index out of bounds (too large)",
			initial:   []int{1, 2, 3},
			index:     3,
			wantValue: 0,
			wantOk:    false,
		},
		{
			name:      "single element stack",
			initial:   []int{42},
			index:     0,
			wantValue: 42,
			wantOk:    true,
		},
		{
			name:      "single element out of bounds",
			initial:   []int{42},
			index:     1,
			wantValue: 0,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.initial...)

			gotValue, gotOk := s.Index(tt.index)

			if gotOk != tt.wantOk {
				t.Errorf("Index(%d) ok = %v, want %v", tt.index, gotOk, tt.wantOk)
			}
			if gotValue != tt.wantValue {
				t.Errorf("Index(%d) value = %v, want %v", tt.index, gotValue, tt.wantValue)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	tests := []struct {
		name         string
		initial      []int
		wantAfterRot []int
	}{
		{
			name:         "empty stack",
			initial:      []int{},
			wantAfterRot: []int{},
		},
		{
			name:         "single element",
			initial:      []int{42},
			wantAfterRot: []int{42},
		},
		{
			name:         "two elements",
			initial:      []int{1, 2},
			wantAfterRot: []int{2, 1},
		},
		{
			name:         "three elements",
			initial:      []int{1, 2, 3},
			wantAfterRot: []int{2, 3, 1},
		},
		{
			name:         "five elements",
			initial:      []int{1, 2, 3, 4, 5},
			wantAfterRot: []int{2, 3, 4, 5, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.initial...)
			s.Rotate()

			if s.Len() != len(tt.wantAfterRot) {
				t.Errorf("Len() = %v, want %v", s.Len(), len(tt.wantAfterRot))
			}

			for i := range len(tt.wantAfterRot) {
				gotValue, gotOk := s.Index(i)
				if !gotOk {
					t.Fatalf("Index(%d) returned false", i)
				}
				if gotValue != tt.wantAfterRot[i] {
					t.Errorf("Index(%d) = %v, want %v", i, gotValue, tt.wantAfterRot[i])
				}
			}
		})
	}
}

func TestReverseRotate(t *testing.T) {
	tests := []struct {
		name         string
		initial      []int
		wantAfterRot []int
	}{
		{
			name:         "empty stack",
			initial:      []int{},
			wantAfterRot: []int{},
		},
		{
			name:         "single element",
			initial:      []int{42},
			wantAfterRot: []int{42},
		},
		{
			name:         "two elements",
			initial:      []int{1, 2},
			wantAfterRot: []int{2, 1},
		},
		{
			name:         "three elements",
			initial:      []int{1, 2, 3},
			wantAfterRot: []int{3, 1, 2},
		},
		{
			name:         "five elements",
			initial:      []int{1, 2, 3, 4, 5},
			wantAfterRot: []int{5, 1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.initial...)
			s.ReverseRotate()

			if s.Len() != len(tt.wantAfterRot) {
				t.Errorf("Len() = %v, want %v", s.Len(), len(tt.wantAfterRot))
			}

			for i := range len(tt.wantAfterRot) {
				gotValue, gotOk := s.Index(i)
				if !gotOk {
					t.Fatalf("Index(%d) returned false", i)
				}
				if gotValue != tt.wantAfterRot[i] {
					t.Errorf("Index(%d) = %v, want %v", i, gotValue, tt.wantAfterRot[i])
				}
			}
		})
	}
}

// TestRotateFullCycle verifies that rotating N elements N times returns them to their original positions.
func TestRotateFullCycle(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
	}{
		{name: "single element", initial: []int{1}},
		{name: "two elements", initial: []int{1, 2}},
		{name: "three elements", initial: []int{1, 2, 3}},
		{name: "four elements", initial: []int{1, 2, 3, 4}},
		{name: "five elements", initial: []int{1, 2, 3, 4, 5}},
		{name: "seven elements", initial: []int{1, 2, 3, 4, 5, 6, 7}},
		{name: "ten elements", initial: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, tt := range tests {
		// Test with New() initialization
		t.Run(tt.name+" via New", func(t *testing.T) {
			s := New(tt.initial...)
			n := s.Len()

			// Rotate N times
			for range n {
				s.Rotate()
			}

			// Verify all elements are back in their original positions
			if s.Len() != len(tt.initial) {
				t.Errorf("Len() = %v, want %v", s.Len(), len(tt.initial))
			}
			for i, want := range tt.initial {
				got, ok := s.Index(i)
				if !ok {
					t.Fatalf("Index(%d) returned false", i)
				}
				if got != want {
					t.Errorf("After %d rotations, Index(%d) = %v, want %v", n, i, got, want)
				}
			}
		})

		// Test with Push() initialization
		t.Run(tt.name+" via Push", func(t *testing.T) {
			s := New[int]()
			// Push in reverse order so that when we read from index 0, we get the first element
			for i := len(tt.initial) - 1; i >= 0; i-- {
				s.Push(tt.initial[i])
			}
			n := s.Len()

			// Rotate N times
			for range n {
				s.Rotate()
			}

			// Verify all elements are back in their original positions
			if s.Len() != len(tt.initial) {
				t.Errorf("Len() = %v, want %v", s.Len(), len(tt.initial))
			}
			for i, want := range tt.initial {
				got, ok := s.Index(i)
				if !ok {
					t.Fatalf("Index(%d) returned false", i)
				}
				if got != want {
					t.Errorf("After %d rotations, Index(%d) = %v, want %v", n, i, got, want)
				}
			}
		})
	}
}

// TestReverseRotateFullCycle verifies that reverse-rotating N elements N times returns them to their original positions.
func TestReverseRotateFullCycle(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
	}{
		{name: "single element", initial: []int{1}},
		{name: "two elements", initial: []int{1, 2}},
		{name: "three elements", initial: []int{1, 2, 3}},
		{name: "four elements", initial: []int{1, 2, 3, 4}},
		{name: "five elements", initial: []int{1, 2, 3, 4, 5}},
		{name: "seven elements", initial: []int{1, 2, 3, 4, 5, 6, 7}},
		{name: "ten elements", initial: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, tt := range tests {
		// Test with New() initialization
		t.Run(tt.name+" via New", func(t *testing.T) {
			s := New(tt.initial...)
			n := s.Len()

			// Reverse-rotate N times
			for range n {
				s.ReverseRotate()
			}

			// Verify all elements are back in their original positions
			if s.Len() != len(tt.initial) {
				t.Errorf("Len() = %v, want %v", s.Len(), len(tt.initial))
			}
			for i, want := range tt.initial {
				got, ok := s.Index(i)
				if !ok {
					t.Fatalf("Index(%d) returned false", i)
				}
				if got != want {
					t.Errorf("After %d reverse-rotations, Index(%d) = %v, want %v", n, i, got, want)
				}
			}
		})

		// Test with Push() initialization
		t.Run(tt.name+" via Push", func(t *testing.T) {
			s := New[int]()
			// Push in reverse order so that when we read from index 0, we get the first element
			for i := len(tt.initial) - 1; i >= 0; i-- {
				s.Push(tt.initial[i])
			}
			n := s.Len()

			// Reverse-rotate N times
			for range n {
				s.ReverseRotate()
			}

			// Verify all elements are back in their original positions
			if s.Len() != len(tt.initial) {
				t.Errorf("Len() = %v, want %v", s.Len(), len(tt.initial))
			}
			for i, want := range tt.initial {
				got, ok := s.Index(i)
				if !ok {
					t.Fatalf("Index(%d) returned false", i)
				}
				if got != want {
					t.Errorf("After %d reverse-rotations, Index(%d) = %v, want %v", n, i, got, want)
				}
			}
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		wantSeq []int
	}{
		{
			name:    "empty stack",
			initial: []int{},
			wantSeq: []int{},
		},
		{
			name:    "single element",
			initial: []int{42},
			wantSeq: []int{42},
		},
		{
			name:    "multiple elements",
			initial: []int{1, 2, 3, 4, 5},
			wantSeq: []int{1, 2, 3, 4, 5},
		},
		{
			name:    "two elements",
			initial: []int{10, 20},
			wantSeq: []int{10, 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.initial...)

			var got []int
			for _, val := range s.All() {
				got = append(got, val)
			}

			if len(got) != len(tt.wantSeq) {
				t.Errorf("All() returned %d elements, want %d", len(got), len(tt.wantSeq))
			}

			for i, v := range got {
				if v != tt.wantSeq[i] {
					t.Errorf("All() element %d = %v, want %v", i, v, tt.wantSeq[i])
				}
			}
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name    string
		initial []int
		want    int
	}{
		{
			name:    "empty stack",
			initial: []int{},
			want:    0,
		},
		{
			name:    "single element",
			initial: []int{42},
			want:    1,
		},
		{
			name:    "multiple elements",
			initial: []int{1, 2, 3, 4, 5},
			want:    5,
		},
		{
			name:    "two elements",
			initial: []int{10, 20},
			want:    2,
		},
		{
			name:    "ten elements",
			initial: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.initial...)

			if got := s.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCombinedOperations(t *testing.T) {
	tests := []struct {
		name       string
		operations func(*Stack[int])
		wantTop    int
		wantLen    int
		wantOk     bool
	}{
		{
			name: "push and pop sequence",
			operations: func(s *Stack[int]) {
				s.Push(1)
				s.Push(2)
				s.Pop()
				s.Push(3)
			},
			wantTop: 3,
			wantLen: 2,
			wantOk:  true,
		},
		{
			name: "rotate after pushes",
			operations: func(s *Stack[int]) {
				s.Push(1)
				s.Push(2)
				s.Push(3)
				s.Rotate()
			},
			wantTop: 2,
			wantLen: 3,
			wantOk:  true,
		},
		{
			name: "reverse rotate after pushes",
			operations: func(s *Stack[int]) {
				s.Push(1)
				s.Push(2)
				s.Push(3)
				s.ReverseRotate()
			},
			wantTop: 1,
			wantLen: 3,
			wantOk:  true,
		},
		{
			name: "complex sequence",
			operations: func(s *Stack[int]) {
				s.Push(1)
				s.Push(2)
				s.Rotate()
				s.Push(3)
				s.Pop()
				s.ReverseRotate()
			},
			wantTop: 2,
			wantLen: 2,
			wantOk:  true,
		},
		{
			name: "pop all then check empty",
			operations: func(s *Stack[int]) {
				s.Push(1)
				s.Push(2)
				s.Pop()
				s.Pop()
			},
			wantTop: 0,
			wantLen: 0,
			wantOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New[int]()
			tt.operations(s)

			if got := s.Len(); got != tt.wantLen {
				t.Errorf("Len() = %v, want %v", got, tt.wantLen)
			}

			gotTop, gotOk := s.Index(0)
			if gotOk != tt.wantOk {
				t.Errorf("Index(0) ok = %v, want %v", gotOk, tt.wantOk)
			}
			if gotOk && gotTop != tt.wantTop {
				t.Errorf("Index(0) = %v, want %v", gotTop, tt.wantTop)
			}
		})
	}
}
