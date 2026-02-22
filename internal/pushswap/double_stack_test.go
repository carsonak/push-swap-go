package pushswap

import (
	"testing"
)

// stackContents reads a stack's elements into a slice from top (index 0) to bottom.
// Using only the public Index/Len API to avoid coupling to the internal implementation.
func stackContents(ds *DStack, which string) []int {
	var s interface {
		Len() int
		Index(int) (int, bool)
	}
	if which == "A" {
		s = &ds.A
	} else {
		s = &ds.B
	}

	out := make([]int, s.Len())
	for i := range s.Len() {
		v, _ := s.Index(i)
		out[i] = v
	}
	return out
}

func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestPushToA(t *testing.T) {
	tests := []struct {
		name   string
		initA  []int
		initB  []int
		wantA  []int
		wantB  []int
		wantOp Operation
	}{
		{
			name:   "push top of B onto A",
			initA:  []int{3, 4},
			initB:  []int{1, 2},
			wantA:  []int{1, 3, 4},
			wantB:  []int{2},
			wantOp: PA,
		},
		{
			name:   "push into empty A",
			initA:  []int{},
			initB:  []int{5},
			wantA:  []int{5},
			wantB:  []int{},
			wantOp: PA,
		},
		{
			name:   "B empty returns invalid, A unchanged",
			initA:  []int{1, 2},
			initB:  []int{},
			wantA:  []int{1, 2},
			wantB:  []int{},
			wantOp: Invalid,
		},
		{
			name:   "both empty returns invalid",
			initA:  []int{},
			initB:  []int{},
			wantA:  []int{},
			wantB:  []int{},
			wantOp: Invalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New(tt.initA...)
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}

			gotOp := ds.PushToA()

			if gotOp != tt.wantOp {
				t.Errorf("PushToA() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.wantA) {
				t.Errorf("A = %v, want %v", gotA, tt.wantA)
			}
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.wantB) {
				t.Errorf("B = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestPushToB(t *testing.T) {
	tests := []struct {
		name   string
		initA  []int
		initB  []int
		wantA  []int
		wantB  []int
		wantOp Operation
	}{
		{
			name:   "push top of A onto B",
			initA:  []int{1, 2},
			initB:  []int{3, 4},
			wantA:  []int{2},
			wantB:  []int{1, 3, 4},
			wantOp: PB,
		},
		{
			name:   "push into empty B",
			initA:  []int{5},
			initB:  []int{},
			wantA:  []int{},
			wantB:  []int{5},
			wantOp: PB,
		},
		{
			name:   "A empty returns invalid, B unchanged",
			initA:  []int{},
			initB:  []int{1, 2},
			wantA:  []int{},
			wantB:  []int{1, 2},
			wantOp: Invalid,
		},
		{
			name:   "both empty returns invalid",
			initA:  []int{},
			initB:  []int{},
			wantA:  []int{},
			wantB:  []int{},
			wantOp: Invalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New(tt.initA...)
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}

			gotOp := ds.PushToB()

			if gotOp != tt.wantOp {
				t.Errorf("PushToB() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.wantA) {
				t.Errorf("A = %v, want %v", gotA, tt.wantA)
			}
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.wantB) {
				t.Errorf("B = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

// TestPushToAAndPushToBInverse verifies that pb then pa is a no-op.
func TestPushToAAndPushToBInverse(t *testing.T) {
	tests := []struct {
		name  string
		initA []int
	}{
		{name: "single element", initA: []int{1}},
		{name: "multiple elements", initA: []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New(tt.initA...)
			ds.PushToB()
			ds.PushToA()

			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.initA) {
				t.Errorf("A after pb+pa = %v, want %v", gotA, tt.initA)
			}
			if gotB := stackContents(ds, "B"); len(gotB) != 0 {
				t.Errorf("B after pb+pa = %v, want []", gotB)
			}
		})
	}
}

func TestSwapA(t *testing.T) {
	tests := []struct {
		name   string
		initA  []int
		wantA  []int
		wantOp Operation
	}{
		{
			name:   "swaps top two elements",
			initA:  []int{1, 2, 3},
			wantA:  []int{2, 1, 3},
			wantOp: SA,
		},
		{
			name:   "two elements swaps both",
			initA:  []int{1, 2},
			wantA:  []int{2, 1},
			wantOp: SA,
		},
		{
			name:   "single element unchanged",
			initA:  []int{42},
			wantA:  []int{42},
			wantOp: SA,
		},
		{
			name:   "empty stack unchanged",
			initA:  []int{},
			wantA:  []int{},
			wantOp: SA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New(tt.initA...)

			gotOp := ds.SwapA()

			if gotOp != tt.wantOp {
				t.Errorf("SwapA() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.wantA) {
				t.Errorf("A = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}

func TestSwapB(t *testing.T) {
	tests := []struct {
		name   string
		initB  []int
		wantB  []int
		wantOp Operation
	}{
		{
			name:   "swaps top two elements",
			initB:  []int{1, 2, 3},
			wantB:  []int{2, 1, 3},
			wantOp: SB,
		},
		{
			name:   "two elements swaps both",
			initB:  []int{1, 2},
			wantB:  []int{2, 1},
			wantOp: SB,
		},
		{
			name:   "single element unchanged",
			initB:  []int{42},
			wantB:  []int{42},
			wantOp: SB,
		},
		{
			name:   "empty stack unchanged",
			initB:  []int{},
			wantB:  []int{},
			wantOp: SB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New()
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}

			gotOp := ds.SwapB()

			if gotOp != tt.wantOp {
				t.Errorf("SwapB() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.wantB) {
				t.Errorf("B = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestSSwap(t *testing.T) {
	tests := []struct {
		name   string
		initA  []int
		initB  []int
		wantA  []int
		wantB  []int
		wantOp Operation
	}{
		{
			name:   "swaps both stacks",
			initA:  []int{1, 2, 3},
			initB:  []int{4, 5, 6},
			wantA:  []int{2, 1, 3},
			wantB:  []int{5, 4, 6},
			wantOp: SS,
		},
		{
			name:   "one empty stack",
			initA:  []int{1, 2},
			initB:  []int{},
			wantA:  []int{2, 1},
			wantB:  []int{},
			wantOp: SS,
		},
		{
			name:   "both empty",
			initA:  []int{},
			initB:  []int{},
			wantA:  []int{},
			wantB:  []int{},
			wantOp: SS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New(tt.initA...)
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}

			gotOp := ds.SSwap()

			if gotOp != tt.wantOp {
				t.Errorf("SSwap() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.wantA) {
				t.Errorf("A = %v, want %v", gotA, tt.wantA)
			}
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.wantB) {
				t.Errorf("B = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestRotateA(t *testing.T) {
	tests := []struct {
		name    string
		initA   []int
		numRots int
		wantA   []int
		wantOp  Operation
	}{
		{
			name:    "top becomes bottom",
			initA:   []int{1, 2, 3},
			numRots: 1,
			wantA:   []int{2, 3, 1},
			wantOp:  RA,
		},
		{
			name:    "two elements",
			initA:   []int{1, 2},
			numRots: 1,
			wantA:   []int{2, 1},
			wantOp:  RA,
		},
		{
			name:    "single element unchanged",
			initA:   []int{42},
			numRots: 1,
			wantA:   []int{42},
			wantOp:  RA,
		},
		{
			name:    "empty stack unchanged",
			initA:   []int{},
			numRots: 1,
			wantA:   []int{},
			wantOp:  RA,
		},
		{
			name:    "full cycle restores original",
			initA:   []int{1, 2, 3, 4, 5},
			numRots: 5,
			wantA:   []int{1, 2, 3, 4, 5},
			wantOp:  RA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New(tt.initA...)
			var gotOp Operation
			for range tt.numRots {
				gotOp = ds.RotateA()
			}

			if gotOp != tt.wantOp {
				t.Errorf("RotateA() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.wantA) {
				t.Errorf("A = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}

func TestRotateB(t *testing.T) {
	tests := []struct {
		name    string
		initB   []int
		numRots int
		wantB   []int
		wantOp  Operation
	}{
		{
			name:    "top becomes bottom",
			initB:   []int{1, 2, 3},
			numRots: 1,
			wantB:   []int{2, 3, 1},
			wantOp:  RB,
		},
		{
			name:    "two elements",
			initB:   []int{1, 2},
			numRots: 1,
			wantB:   []int{2, 1},
			wantOp:  RB,
		},
		{
			name:    "single element unchanged",
			initB:   []int{42},
			numRots: 1,
			wantB:   []int{42},
			wantOp:  RB,
		},
		{
			name:    "empty stack unchanged",
			initB:   []int{},
			numRots: 1,
			wantB:   []int{},
			wantOp:  RB,
		},
		{
			name:    "full cycle restores original",
			initB:   []int{1, 2, 3, 4, 5},
			numRots: 5,
			wantB:   []int{1, 2, 3, 4, 5},
			wantOp:  RB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New()
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}
			var gotOp Operation
			for range tt.numRots {
				gotOp = ds.RotateB()
			}

			if gotOp != tt.wantOp {
				t.Errorf("RotateB() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.wantB) {
				t.Errorf("B = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestRRotate(t *testing.T) {
	tests := []struct {
		name   string
		initA  []int
		initB  []int
		wantA  []int
		wantB  []int
		wantOp Operation
	}{
		{
			name:   "rotates both stacks",
			initA:  []int{1, 2, 3},
			initB:  []int{4, 5, 6},
			wantA:  []int{2, 3, 1},
			wantB:  []int{5, 6, 4},
			wantOp: RR,
		},
		{
			name:   "one empty stack",
			initA:  []int{1, 2, 3},
			initB:  []int{},
			wantA:  []int{2, 3, 1},
			wantB:  []int{},
			wantOp: RR,
		},
		{
			name:   "both empty",
			initA:  []int{},
			initB:  []int{},
			wantA:  []int{},
			wantB:  []int{},
			wantOp: RR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New(tt.initA...)
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}

			gotOp := ds.RRotate()

			if gotOp != tt.wantOp {
				t.Errorf("RRotate() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.wantA) {
				t.Errorf("A = %v, want %v", gotA, tt.wantA)
			}
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.wantB) {
				t.Errorf("B = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestReverseRotateA(t *testing.T) {
	tests := []struct {
		name    string
		initA   []int
		numRots int
		wantA   []int
		wantOp  Operation
	}{
		{
			name:    "bottom becomes top",
			initA:   []int{1, 2, 3},
			numRots: 1,
			wantA:   []int{3, 1, 2},
			wantOp:  RRA,
		},
		{
			name:    "two elements",
			initA:   []int{1, 2},
			numRots: 1,
			wantA:   []int{2, 1},
			wantOp:  RRA,
		},
		{
			name:    "single element unchanged",
			initA:   []int{42},
			numRots: 1,
			wantA:   []int{42},
			wantOp:  RRA,
		},
		{
			name:    "empty stack unchanged",
			initA:   []int{},
			numRots: 1,
			wantA:   []int{},
			wantOp:  RRA,
		},
		{
			name:    "full cycle restores original",
			initA:   []int{1, 2, 3, 4, 5},
			numRots: 5,
			wantA:   []int{1, 2, 3, 4, 5},
			wantOp:  RRA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New(tt.initA...)
			var gotOp Operation
			for range tt.numRots {
				gotOp = ds.ReverseRotateA()
			}

			if gotOp != tt.wantOp {
				t.Errorf("ReverseRotateA() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.wantA) {
				t.Errorf("A = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}

func TestReverseRotateB(t *testing.T) {
	tests := []struct {
		name    string
		initB   []int
		numRots int
		wantB   []int
		wantOp  Operation
	}{
		{
			name:    "bottom becomes top",
			initB:   []int{1, 2, 3},
			numRots: 1,
			wantB:   []int{3, 1, 2},
			wantOp:  RRB,
		},
		{
			name:    "two elements",
			initB:   []int{1, 2},
			numRots: 1,
			wantB:   []int{2, 1},
			wantOp:  RRB,
		},
		{
			name:    "single element unchanged",
			initB:   []int{42},
			numRots: 1,
			wantB:   []int{42},
			wantOp:  RRB,
		},
		{
			name:    "empty stack unchanged",
			initB:   []int{},
			numRots: 1,
			wantB:   []int{},
			wantOp:  RRB,
		},
		{
			name:    "full cycle restores original",
			initB:   []int{1, 2, 3, 4, 5},
			numRots: 5,
			wantB:   []int{1, 2, 3, 4, 5},
			wantOp:  RRB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New()
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}
			var gotOp Operation
			for range tt.numRots {
				gotOp = ds.ReverseRotateB()
			}

			if gotOp != tt.wantOp {
				t.Errorf("ReverseRotateB() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.wantB) {
				t.Errorf("B = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestRReverseRotate(t *testing.T) {
	tests := []struct {
		name   string
		initA  []int
		initB  []int
		wantA  []int
		wantB  []int
		wantOp Operation
	}{
		{
			name:   "reverse-rotates both stacks",
			initA:  []int{1, 2, 3},
			initB:  []int{4, 5, 6},
			wantA:  []int{3, 1, 2},
			wantB:  []int{6, 4, 5},
			wantOp: RRR,
		},
		{
			name:   "one empty stack",
			initA:  []int{1, 2, 3},
			initB:  []int{},
			wantA:  []int{3, 1, 2},
			wantB:  []int{},
			wantOp: RRR,
		},
		{
			name:   "both empty",
			initA:  []int{},
			initB:  []int{},
			wantA:  []int{},
			wantB:  []int{},
			wantOp: RRR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := New(tt.initA...)
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}

			gotOp := ds.RReverseRotate()

			if gotOp != tt.wantOp {
				t.Errorf("RReverseRotate() = %q, want %q", gotOp, tt.wantOp)
			}
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.wantA) {
				t.Errorf("A = %v, want %v", gotA, tt.wantA)
			}
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.wantB) {
				t.Errorf("B = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

// TestRotateAndReverseRotateAreInverse verifies ra+rra, rb+rrb, and rr+rrr all restore original state.
func TestRotateAndReverseRotateAreInverse(t *testing.T) {
	tests := []struct {
		name  string
		initA []int
		initB []int
	}{
		{name: "single elements", initA: []int{1}, initB: []int{2}},
		{name: "multiple elements", initA: []int{1, 2, 3}, initB: []int{4, 5, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name+" ra then rra", func(t *testing.T) {
			ds := New(tt.initA...)
			ds.RotateA()
			ds.ReverseRotateA()
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.initA) {
				t.Errorf("A after ra+rra = %v, want %v", gotA, tt.initA)
			}
		})

		t.Run(tt.name+" rb then rrb", func(t *testing.T) {
			ds := New()
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}
			ds.RotateB()
			ds.ReverseRotateB()
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.initB) {
				t.Errorf("B after rb+rrb = %v, want %v", gotB, tt.initB)
			}
		})

		t.Run(tt.name+" rr then rrr", func(t *testing.T) {
			ds := New(tt.initA...)
			for _, v := range tt.initB {
				ds.B.PushBottom(v)
			}
			ds.RRotate()
			ds.RReverseRotate()
			if gotA := stackContents(ds, "A"); !slicesEqual(gotA, tt.initA) {
				t.Errorf("A after rr+rrr = %v, want %v", gotA, tt.initA)
			}
			if gotB := stackContents(ds, "B"); !slicesEqual(gotB, tt.initB) {
				t.Errorf("B after rr+rrr = %v, want %v", gotB, tt.initB)
			}
		})
	}
}

// TestSSwapEquivalentToSAAndSB verifies ss produces the same result as sa then sb independently.
func TestSSwapEquivalentToSAAndSB(t *testing.T) {
	initA := []int{1, 2, 3}
	initB := []int{4, 5, 6}

	ds1 := New(initA...)
	for _, v := range initB {
		ds1.B.Push(v)
	}
	ds1.SSwap()

	ds2 := New(initA...)
	for _, v := range initB {
		ds2.B.Push(v)
	}
	ds2.SwapA()
	ds2.SwapB()

	if gotA1, gotA2 := stackContents(ds1, "A"), stackContents(ds2, "A"); !slicesEqual(gotA1, gotA2) {
		t.Errorf("ss A = %v, sa A = %v; want equal", gotA1, gotA2)
	}
	if gotB1, gotB2 := stackContents(ds1, "B"), stackContents(ds2, "B"); !slicesEqual(gotB1, gotB2) {
		t.Errorf("ss B = %v, sb B = %v; want equal", gotB1, gotB2)
	}
}

// TestRRotateEquivalentToRAAndRB verifies rr produces the same result as ra then rb independently.
func TestRRotateEquivalentToRAAndRB(t *testing.T) {
	initA := []int{1, 2, 3}
	initB := []int{4, 5, 6}

	ds1 := New(initA...)
	for _, v := range initB {
		ds1.B.Push(v)
	}
	ds1.RRotate()

	ds2 := New(initA...)
	for _, v := range initB {
		ds2.B.Push(v)
	}
	ds2.RotateA()
	ds2.RotateB()

	if gotA1, gotA2 := stackContents(ds1, "A"), stackContents(ds2, "A"); !slicesEqual(gotA1, gotA2) {
		t.Errorf("rr A = %v, ra A = %v; want equal", gotA1, gotA2)
	}
	if gotB1, gotB2 := stackContents(ds1, "B"), stackContents(ds2, "B"); !slicesEqual(gotB1, gotB2) {
		t.Errorf("rr B = %v, rb B = %v; want equal", gotB1, gotB2)
	}
}

// TestRReverseRotateEquivalentToRRAAndRRB verifies rrr produces the same result as rra then rrb independently.
func TestRReverseRotateEquivalentToRRAAndRRB(t *testing.T) {
	initA := []int{1, 2, 3}
	initB := []int{4, 5, 6}

	ds1 := New(initA...)
	for _, v := range initB {
		ds1.B.Push(v)
	}
	ds1.RReverseRotate()

	ds2 := New(initA...)
	for _, v := range initB {
		ds2.B.Push(v)
	}
	ds2.ReverseRotateA()
	ds2.ReverseRotateB()

	if gotA1, gotA2 := stackContents(ds1, "A"), stackContents(ds2, "A"); !slicesEqual(gotA1, gotA2) {
		t.Errorf("rrr A = %v, rra A = %v; want equal", gotA1, gotA2)
	}
	if gotB1, gotB2 := stackContents(ds1, "B"), stackContents(ds2, "B"); !slicesEqual(gotB1, gotB2) {
		t.Errorf("rrr B = %v, rrb B = %v; want equal", gotB1, gotB2)
	}
}
