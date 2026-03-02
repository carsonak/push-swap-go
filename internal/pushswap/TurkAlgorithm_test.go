package pushswap

import (
	"math"
	"slices"
	"testing"
)

// --- helpers ---

// makeStack builds a psStack with the given values (top to bottom) via DoubleStack.A.
func makeStack(vals ...float64) *psStack {
	ds := NewDoubleStack(vals...)
	return &ds.A
}

// stackVals reads all values from a psStack top to bottom.
func stackVals(s *psStack) []float64 {
	out := make([]float64, 0, s.Len())

	for _, v := range s.All() {
		out = append(out, v)
	}

	return out
}

// dsContents reads top-to-bottom values from one side of a DoubleStack for assertions.
func dsContents(ds *DoubleStack, which string) []float64 {
	if which == "A" {
		return stackVals(&ds.A)
	}

	return stackVals(&ds.B)
}

func f64sEqual(a, b []float64) bool {
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

func intSlicesEqualIgnoreOrder(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	aSorted := slices.Clone(a)
	bSorted := slices.Clone(b)
	slices.Sort(aSorted)
	slices.Sort(bSorted)
	return slices.Equal(aSorted, bSorted)
}

// verifyTurkResult applies ops to a fresh DoubleStack and checks A is sorted
// ascending (smallest on top) and B is empty.
func verifyTurkResult(t *testing.T, nums []float64, ops []Operation) {
	t.Helper()
	ds := NewDoubleStack(nums...)

	executeInstructions(ds, ops)
	b := dsContents(ds, "B")

	if len(b) != 0 {
		t.Errorf("B is not empty after applying ops: %v", b)
	}

	a := dsContents(ds, "A")

	if !slices.IsSorted(a) {
		t.Errorf("A is not sorted after applying ops: %v", a)
	}

	// Verify A contains exactly the same values as the input (no data loss).
	want := make([]float64, len(nums))

	copy(want, nums)
	slices.Sort(want)
	if !f64sEqual(a, want) {
		t.Errorf("A values = %v, want %v", a, want)
	}
}

// --- TestFindExtremes (whitebox) ---
func TestFindExtremes(t *testing.T) {
	tests := []struct {
		name     string
		vals     []float64
		findMax  bool
		wantIdxs []int
	}{
		{
			name:     "empty stack - max",
			vals:     []float64{},
			findMax:  true,
			wantIdxs: nil,
		},
		{
			name:     "empty stack - min",
			vals:     []float64{},
			findMax:  false,
			wantIdxs: nil,
		},
		{
			name:     "single element - max",
			vals:     []float64{7},
			findMax:  true,
			wantIdxs: []int{0},
		},
		{
			name:     "single element - min",
			vals:     []float64{7},
			findMax:  false,
			wantIdxs: []int{0},
		},
		{
			name:     "multiple distinct - max",
			vals:     []float64{3, 1, 4, 1, 5, 9},
			findMax:  true,
			wantIdxs: []int{5},
		},
		{
			// Regression: Bug 1 - returned [5] (last index) instead of [1, 3].
			name:     "multiple distinct - min",
			vals:     []float64{3, 1, 4, 1, 5, 9},
			findMax:  false,
			wantIdxs: []int{1, 3},
		},
		{
			// Regression: Bug 1 - returned [2] (last index) instead of [0, 1, 2].
			name:     "all equal - max",
			vals:     []float64{4, 4, 4},
			findMax:  true,
			wantIdxs: []int{0, 1, 2},
		},
		{
			name:     "ascending - max at bottom",
			vals:     []float64{1, 2, 3, 4, 5},
			findMax:  true,
			wantIdxs: []int{4},
		},
		{
			// Regression: Bug 1 - returned [4] (last index) instead of [0].
			name:     "ascending - min at top",
			vals:     []float64{1, 2, 3, 4, 5},
			findMax:  false,
			wantIdxs: []int{0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := makeStack(tt.vals...)
			got := findExtremes(s, tt.findMax)

			if len(got) != len(tt.wantIdxs) {
				t.Fatalf("findExtremes() = %v, want %v", got, tt.wantIdxs)
			}
			for i, idx := range got {
				if idx != tt.wantIdxs[i] {
					t.Errorf("findExtremes()[%d] = %d, want %d", i, idx, tt.wantIdxs[i])
				}
			}
		})
	}
}

// --- TestSortLast3 (whitebox) ---

func TestSortLast3(t *testing.T) {
	tests := []struct {
		name      string
		vals      []float64
		wantOrder []float64 // expected top-to-bottom after sort
		wantOps   []Operation
	}{
		{
			name:      "empty - no-op",
			vals:      []float64{},
			wantOrder: []float64{},
			wantOps:   nil,
		},
		{
			name:      "one element - no-op",
			vals:      []float64{5},
			wantOrder: []float64{5},
			wantOps:   nil,
		},
		{
			name:      "four or more elements - no-op",
			vals:      []float64{4, 3, 2, 1},
			wantOrder: []float64{4, 3, 2, 1},
			wantOps:   nil,
		},
		{
			name:      "two elements - sorted",
			vals:      []float64{1, 2},
			wantOrder: []float64{1, 2},
			wantOps:   nil,
		},
		{
			name:      "two elements - unsorted",
			vals:      []float64{2, 1},
			wantOrder: []float64{1, 2},
			wantOps:   []Operation{SA},
		},
		{
			name:      "three elements - already sorted",
			vals:      []float64{1, 2, 3},
			wantOrder: []float64{1, 2, 3},
			wantOps:   nil,
		},
		{
			name:      "three elements - [2,1,3]",
			vals:      []float64{2, 1, 3},
			wantOrder: []float64{1, 2, 3},
			wantOps:   []Operation{SA},
		},
		{
			// max(3) is at top; RA gives [1, 2, 3] which is already sorted - no swap needed.
			name:      "three elements - [3,1,2] (max on top)",
			vals:      []float64{3, 1, 2},
			wantOrder: []float64{1, 2, 3},
			wantOps:   []Operation{RA},
		},
		{
			name:      "three elements - [1,3,2] (max in middle)",
			vals:      []float64{1, 3, 2},
			wantOrder: []float64{1, 2, 3},
			wantOps:   []Operation{RRA, SA},
		},
		{
			name:      "three elements - [2,3,1]",
			vals:      []float64{2, 3, 1},
			wantOrder: []float64{1, 2, 3},
			wantOps:   []Operation{RRA},
		},
		{
			name:      "three elements - [3,2,1]",
			vals:      []float64{3, 2, 1},
			wantOrder: []float64{1, 2, 3},
			wantOps:   []Operation{RA, SA}, // max(3) at top → RA → [2,1,3] → SA → [1,2,3]
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := makeStack(tt.vals...)
			got := sortLast3(s)

			// Verify returned operations match expectation.
			if len(got) != len(tt.wantOps) {
				t.Errorf("ops = %v, want %v", got, tt.wantOps)
			} else {
				for i, op := range got {
					if op != tt.wantOps[i] {
						t.Errorf("ops[%d] = %q, want %q", i, op, tt.wantOps[i])
					}
				}
			}

			// Verify the stack is actually in the expected order.
			final := stackVals(s)
			if !f64sEqual(final, tt.wantOrder) {
				t.Errorf("stack after sort = %v, want %v", final, tt.wantOrder)
			}
		})
	}
}

// --- TestShortestRouteToTop (whitebox) ---

func TestShortestRouteToTop(t *testing.T) {
	tests := []struct {
		name   string
		index  int
		length int
		want   int
	}{
		{name: "already at top", index: 0, length: 5, want: 0},
		{name: "one from top", index: 1, length: 5, want: 1},
		{name: "exact middle - forward cheaper", index: 2, length: 5, want: 2},
		{name: "just past middle - reverse cheaper", index: 3, length: 5, want: 2},
		{name: "one from bottom", index: 4, length: 5, want: 1},
		{name: "bottom", index: 4, length: 5, want: 1},
		{name: "even length - exact half", index: 3, length: 6, want: 3},
		{name: "even length - past half", index: 4, length: 6, want: 2},
		{name: "single element stack", index: 0, length: 1, want: 0},
		{name: "two element stack - top", index: 0, length: 2, want: 0},
		{name: "two element stack - bottom", index: 1, length: 2, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shortestRouteToTop(tt.index, tt.length)
			if got != tt.want {
				t.Errorf("shortestRouteToTop(%d, %d) = %d, want %d", tt.index, tt.length, got, tt.want)
			}
		})
	}
}

// --- TestFindTargets (whitebox) ---
func TestFindGreaterTargets(t *testing.T) {
	tests := []struct {
		name     string
		vals     []float64
		ref      float64
		wantIdxs []int
	}{
		{
			name:     "empty stack",
			vals:     []float64{},
			ref:      5,
			wantIdxs: nil,
		},
		{
			name:     "no values greater than ref",
			vals:     []float64{1, 2, 3},
			ref:      5,
			wantIdxs: nil,
		},
		{
			name:     "all values greater - returns closest",
			vals:     []float64{3, 4, 5, 2}, // ascending rotated [2,3,4,5], closest (2) is at bottom
			ref:      1,
			wantIdxs: []int{3}, // val=2, the closest > 1
		},
		{
			name:     "single match",
			vals:     []float64{1, 5, 3},
			ref:      4,
			wantIdxs: []int{1}, // val=5 is the only one > 4
		},
		{
			// Regression: Bug 2 - returned [3] (last qualifying idx) instead of [1, 2].
			name:     "duplicate closest values",
			vals:     []float64{1, 3, 3, 5},
			ref:      2,
			wantIdxs: []int{1, 2}, // both val=3 are closest > 2
		},
		{
			// Regression: Bug 2 - returned index 3 (val=4) instead of index 2 (val=2,
			// the actual closest value > 1.5).
			name:     "rotated - closest is not at the top",
			vals:     []float64{3, 4, 2}, // ascending rotated [2,3,4]
			ref:      1.5,
			wantIdxs: []int{2}, // val=2 is the smallest > 1.5
		},
		{
			name:     "ref equals a value - only strictly greater",
			vals:     []float64{3, 3, 5},
			ref:      3,
			wantIdxs: []int{2}, // only val=5 is strictly > 3
		},
		{
			name:     "pivot at top index 0",
			vals:     []float64{1, 2, 3, 4}, // already ascending => pivot=0
			ref:      2.5,
			wantIdxs: []int{2}, // val=3
		},
		{
			name:     "pivot at bottom index",
			vals:     []float64{2, 3, 4, 1}, // ascending rotated => pivot=3 (bottom)
			ref:      0.5,
			wantIdxs: []int{3}, // val=1 comes from second segment [pivot..end]
		},
		{
			name:     "all elements equal",
			vals:     []float64{7, 7, 7},
			ref:      6,
			wantIdxs: []int{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := makeStack(tt.vals...)
			got := findGreaterTargetsInSortedRotated(s, tt.ref)

			if !intSlicesEqualIgnoreOrder(got, tt.wantIdxs) {
				t.Errorf("findGreaterTargets() = %v, want %v", got, tt.wantIdxs)
			}
		})
	}
}

func TestFindSmallerTargets(t *testing.T) {
	tests := []struct {
		name     string
		vals     []float64
		ref      float64
		wantIdxs []int
	}{
		{
			name:     "empty stack",
			vals:     []float64{},
			ref:      5,
			wantIdxs: nil,
		},
		{
			name:     "no values smaller than ref",
			vals:     []float64{5, 6, 7},
			ref:      3,
			wantIdxs: nil,
		},
		{
			// Regression: Bug 2 - returned [3] (last idx) instead of [2].
			name:     "all values smaller - returns closest",
			vals:     []float64{2, 1, 4, 3}, // descending rotated [4,3,2,1], closest (4) is at index 2
			ref:      5,
			wantIdxs: []int{2}, // val=4, the largest < 5
		},
		{
			name:     "single match",
			vals:     []float64{5, 1, 3},
			ref:      2,
			wantIdxs: []int{1}, // only val=1 (at index 1) is < 2
		},
		{
			// Regression: Bug 2 - returned [3] (last qualifying idx) instead of [1, 2].
			name:     "duplicate closest values",
			vals:     []float64{5, 3, 3, 1},
			ref:      4,
			wantIdxs: []int{1, 2}, // both val=3 are closest < 4
		},
		{
			// Regression: Bug 2 - returned index 3 (val=1) instead of index 2 (val=3,
			// the actual closest value < 3.5).
			name:     "rotated - closest is not at the top",
			vals:     []float64{1, 4, 3, 2}, // descending rotated [4,3,2,1]
			ref:      3.5,
			wantIdxs: []int{2}, // val=3 is the largest < 3.5
		},
		{
			name:     "pivot at top index 0",
			vals:     []float64{9, 7, 5, 3}, // already descending => pivot=0
			ref:      8,
			wantIdxs: []int{1}, // val=7
		},
		{
			name:     "pivot at bottom index",
			vals:     []float64{4, 3, 2, 5}, // descending rotated => pivot=3 (bottom)
			ref:      5.5,
			wantIdxs: []int{3}, // val=5 comes from second segment [pivot..end]
		},
		{
			name:     "all elements equal",
			vals:     []float64{7, 7, 7},
			ref:      8,
			wantIdxs: []int{0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := makeStack(tt.vals...)
			got := findSmallerTargetsInSortedRotated(s, tt.ref)

			if !intSlicesEqualIgnoreOrder(got, tt.wantIdxs) {
				t.Errorf("findSmallerTargets() = %v, want %v", got, tt.wantIdxs)
			}
		})
	}
}

// --- TestGenerateInstructions (whitebox) ---
func TestGenerateInstructions(t *testing.T) {
	// buildDS creates a DoubleStack with specific A and B stack sizes (filled
	// with placeholder values) so we can test rotation index arithmetic.
	buildDS := func(aVals, bVals []float64) *DoubleStack {
		ds := NewDoubleStack(aVals...)
		for _, v := range bVals {
			ds.B.PushBottom(v)
		}
		return ds
	}

	tests := []struct {
		name    string
		aVals   []float64
		bVals   []float64
		move    moveCandidate
		wantOps []Operation
	}{
		{
			// Indices at top: no rotations needed at all.
			name:    "both already at top - no ops",
			aVals:   []float64{1, 2, 3, 4, 5},
			bVals:   []float64{6, 7, 8, 9, 10},
			move:    moveCandidate{fromIdx: 0, toIdx: 0, target: stackB},
			wantOps: []Operation{},
		},
		{
			// route=RR: forward-rotate both simultaneously, then A alone.
			name:  "route RR - combined then A solo",
			aVals: []float64{1, 2, 3, 4, 5},
			bVals: []float64{6, 7, 8, 9, 10},
			move:  moveCandidate{fromIdx: 3, toIdx: 2, target: stackB, route: RR},
			// idxA=fromIdx=3, idxB=toIdx=2
			// ARotations=3 (3 ≤ 5/2 is false? 3 > 2 yes, but route!=RRR, so
			// condition: route==RRR(false) || idxA>5/2(3>2=true) → ARotations=5-3=2
			// BRotations: route==RRR(false) || idxB>5/2(2>2=false) → BRotations=2
			// Combined: min(2,2)=2 RR, ARotations=0, BRotations=0
			wantOps: []Operation{RR, RR},
		},
		{
			// route=RR with unequal distances: combined ops, then B solo.
			name:  "route RR - combined then B solo",
			aVals: []float64{1, 2, 3, 4, 5},
			bVals: []float64{6, 7, 8, 9, 10},
			move:  moveCandidate{fromIdx: 1, toIdx: 3, target: stackB, route: RR},
			// idxA=fromIdx=1, idxB=toIdx=3
			// ARotations: !RRR, 1 ≤ 2 → ARotations=1
			// BRotations: !RRR, 3 > 2 → BRotations=5-3=2
			// Combined: min(1,2)=1 RR, ARotations-=1→0, BRotations=1
			wantOps: []Operation{RR, RB},
		},
		{
			// route=RRR with equal distances: only combined ops.
			name:  "route RRR - equal reverse distances",
			aVals: []float64{1, 2, 3, 4, 5},
			bVals: []float64{6, 7, 8, 9, 10},
			move:  moveCandidate{fromIdx: 4, toIdx: 4, target: stackB, route: RRR},
			// idxA=4, idxB=4; ARotations=5-4=1; BRotations=5-4=1
			// Combined: 1 RRR; no remaining
			wantOps: []Operation{RRR},
		},
		{
			// Regression: Bug 3 - emitted [RRR, RA, RA] instead of [RRR, RRA, RRA].
			name:  "route RRR - remaining A rotations use RRA",
			aVals: []float64{1, 2, 3, 4, 5},
			bVals: []float64{6, 7, 8, 9, 10},
			move:  moveCandidate{fromIdx: 2, toIdx: 4, target: stackB, route: RRR},
			// idxA=2, idxB=4; ARotations=5-2=3; BRotations=5-4=1
			// Combined: 1 RRR; ARotations-=1→2; BRotations=0
			wantOps: []Operation{RRR, RRA, RRA},
		},
		{
			// route="" (separate), forward only: RA + RB.
			name:  "route separate - both forward",
			aVals: []float64{1, 2, 3, 4, 5},
			bVals: []float64{6, 7, 8, 9, 10},
			move:  moveCandidate{fromIdx: 2, toIdx: 1, target: stackB},
			// idxA=2, idxB=1; !RRR, 2 ≤ 2 → ARotations=2; !RRR, 1 ≤ 2 → BRotations=1
			wantOps: []Operation{RA, RA, RB},
		},
		{
			// Regression: Bug 3 - emitted [RA, RA, RB] instead of [RRA, RRA, RB].
			name:  "route separate - A past midpoint uses RRA",
			aVals: []float64{1, 2, 3, 4, 5},
			bVals: []float64{6, 7, 8, 9, 10},
			move:  moveCandidate{fromIdx: 3, toIdx: 1, target: stackB},
			// idxA=3, idxB=1; !RRR, 3>2 → ARotations=5-3=2; !RRR, 1≤2 → BRotations=1
			wantOps: []Operation{RRA, RRA, RB},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := buildDS(tt.aVals, tt.bVals)
			got := generateInstructions(ds, tt.move)

			if len(got) != len(tt.wantOps) {
				t.Errorf("generateInstructions() = %v (len=%d), want %v (len=%d)", got, len(got), tt.wantOps, len(tt.wantOps))
				return
			}
			for i, op := range got {
				if op != tt.wantOps[i] {
					t.Errorf("op[%d] = %q, want %q", i, op, tt.wantOps[i])
				}
			}
		})
	}
}

// --- TestExecuteInstructions (whitebox) ---

func TestExecuteInstructions(t *testing.T) {
	// Each case starts from A=[1,2,3], B=[4,5,6] and applies the ops,
	// then checks the expected A and B state.
	initA := []float64{1, 2, 3}
	initB := []float64{4, 5, 6}

	buildDS := func() *DoubleStack {
		ds := NewDoubleStack(initA...)
		for _, v := range initB {
			ds.B.PushBottom(v)
		}
		return ds
	}

	tests := []struct {
		name  string
		ops   []Operation
		wantA []float64
		wantB []float64
	}{
		{
			name:  "no ops - unchanged",
			ops:   nil,
			wantA: []float64{1, 2, 3},
			wantB: []float64{4, 5, 6},
		},
		{
			name:  "SA swaps top two of A",
			ops:   []Operation{SA},
			wantA: []float64{2, 1, 3},
			wantB: []float64{4, 5, 6},
		},
		{
			name:  "SB swaps top two of B",
			ops:   []Operation{SB},
			wantA: []float64{1, 2, 3},
			wantB: []float64{5, 4, 6},
		},
		{
			name:  "SS swaps top two of both",
			ops:   []Operation{SS},
			wantA: []float64{2, 1, 3},
			wantB: []float64{5, 4, 6},
		},
		{
			name:  "RA rotates A",
			ops:   []Operation{RA},
			wantA: []float64{2, 3, 1},
			wantB: []float64{4, 5, 6},
		},
		{
			name:  "RB rotates B",
			ops:   []Operation{RB},
			wantA: []float64{1, 2, 3},
			wantB: []float64{5, 6, 4},
		},
		{
			name:  "RR rotates both",
			ops:   []Operation{RR},
			wantA: []float64{2, 3, 1},
			wantB: []float64{5, 6, 4},
		},
		{
			name:  "RRA reverse-rotates A",
			ops:   []Operation{RRA},
			wantA: []float64{3, 1, 2},
			wantB: []float64{4, 5, 6},
		},
		{
			name:  "RRB reverse-rotates B",
			ops:   []Operation{RRB},
			wantA: []float64{1, 2, 3},
			wantB: []float64{6, 4, 5},
		},
		{
			name:  "RRR reverse-rotates both",
			ops:   []Operation{RRR},
			wantA: []float64{3, 1, 2},
			wantB: []float64{6, 4, 5},
		},
		{
			name:  "PB pushes top of A to B",
			ops:   []Operation{PB},
			wantA: []float64{2, 3},
			wantB: []float64{1, 4, 5, 6},
		},
		{
			name:  "PA pushes top of B to A",
			ops:   []Operation{PA},
			wantA: []float64{4, 1, 2, 3},
			wantB: []float64{5, 6},
		},
		{
			name:  "invalid op is ignored",
			ops:   []Operation{Invalid},
			wantA: []float64{1, 2, 3},
			wantB: []float64{4, 5, 6},
		},
		{
			name:  "sequence: PB then PA round-trips top element",
			ops:   []Operation{PB, PA},
			wantA: []float64{1, 2, 3},
			wantB: []float64{4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := buildDS()
			executeInstructions(ds, tt.ops)

			gotA := dsContents(ds, "A")
			if !f64sEqual(gotA, tt.wantA) {
				t.Errorf("A = %v, want %v", gotA, tt.wantA)
			}
			gotB := dsContents(ds, "B")
			if !f64sEqual(gotB, tt.wantB) {
				t.Errorf("B = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

// --- TestTurkAlgorithm (blackbox) ---
//
// This is the primary integration test. For each case the returned instruction
// sequence is verified by re-executing it on a fresh DoubleStack and confirming
// that A ends up sorted ascending with B empty.

func TestTurkAlgorithm(t *testing.T) {
	tests := []struct {
		name    string
		input   []float64
		wantNil bool // whether nil (empty) operations are expected
	}{
		{
			name:    "nil / empty input",
			input:   nil,
			wantNil: true,
		},
		{
			name:    "single element",
			input:   []float64{1},
			wantNil: true,
		},
		{
			name:    "already sorted - two elements",
			input:   []float64{1, 2},
			wantNil: true,
		},
		{
			name:    "already sorted - five elements",
			input:   []float64{1, 2, 3, 4, 5},
			wantNil: true,
		},
		{
			name:  "two elements - reverse order",
			input: []float64{2, 1},
		},
		{
			name:  "three elements - reverse order",
			input: []float64{3, 2, 1},
		},
		{
			name:  "three elements - various permutation",
			input: []float64{2, 3, 1},
		},
		{
			name:  "three elements - max on top",
			input: []float64{3, 1, 2},
		},
		{
			name:  "four elements",
			input: []float64{4, 2, 3, 1},
		},
		{
			name:  "five elements - reverse",
			input: []float64{5, 4, 3, 2, 1},
		},
		{
			name:  "five elements - arbitrary",
			input: []float64{3, 1, 4, 5, 2},
		},
		{
			name:  "eight elements",
			input: []float64{8, 3, 6, 1, 7, 2, 5, 4},
		},
		{
			name:  "ten elements",
			input: []float64{10, 1, 9, 2, 8, 3, 7, 4, 6, 5},
		},
		{
			name:  "negative numbers",
			input: []float64{-3, 0, -1, 2, -5},
		},
		{
			name:  "float values",
			input: []float64{3.5, 1.1, 2.7, 0.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ops := TurkAlgorithm(tt.input)

			if tt.wantNil {
				if ops != nil {
					t.Errorf("TurkAlgorithm() = %v, want nil", ops)
				}

				return
			}

			verifyTurkResult(t, tt.input, ops)
		})
	}
}

// TestTurkAlgorithmProducesValidOpsOnly checks that no Invalid operation is
// emitted by TurkAlgorithm.
func TestTurkAlgorithmProducesValidOpsOnly(t *testing.T) {
	validOps := map[Operation]bool{
		PA: true, PB: true, SA: true, SB: true, SS: true,
		RA: true, RB: true, RR: true, RRA: true, RRB: true, RRR: true,
	}

	inputs := [][]float64{
		{5, 4, 3, 2, 1},
		{3, 1, 4, 5, 2},
		{8, 3, 6, 1, 7, 2, 5, 4},
	}

	for _, input := range inputs {
		ops := TurkAlgorithm(input)
		for _, op := range ops {
			if !validOps[op] {
				t.Errorf("TurkAlgorithm(%v) emitted invalid op %q", input, op)
			}
		}
	}
}

// --- TestFindCheapestMove (whitebox) ---
func TestFindCheapestMove(t *testing.T) {
	tests := []struct {
		name          string
		aVals         []float64
		bVals         []float64
		to            stackID
		wantValidMove bool // fromIdx and toIdx must both be ≥ 0
	}{
		{
			// Regression: Bug 5 - loop guard `top < bottom` skipped the only
			// iteration when source had one element, returning invalid sentinel.
			name:          "single element in source (B→A) - must find valid move",
			aVals:         []float64{1, 3, 5},
			bVals:         []float64{4},
			to:            stackA,
			wantValidMove: true,
		},
		{
			name:          "multiple elements in source (B→A)",
			aVals:         []float64{1, 3, 5},
			bVals:         []float64{2, 4},
			to:            stackA,
			wantValidMove: true,
		},
		{
			// Regression: Bug 5 - same loop guard issue as the B→A case.
			name:          "single element in source (A→B) - must find valid move",
			aVals:         []float64{3},
			bVals:         []float64{1, 4, 6},
			to:            stackB,
			wantValidMove: true,
		},
		{
			name:          "two elements in source (A→B)",
			aVals:         []float64{3, 7},
			bVals:         []float64{1, 4, 6},
			to:            stackB,
			wantValidMove: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDoubleStack(tt.aVals...)
			for _, v := range tt.bVals {
				ds.B.PushBottom(v)
			}

			got := findCheapestMove(ds, tt.to)

			if tt.wantValidMove && (got.fromIdx < 0 || got.toIdx < 0) {
				t.Errorf("findCheapestMove() returned invalid move {fromIdx:%d, toIdx:%d, cost:%d}, want a valid move",
					got.fromIdx, got.toIdx, got.cost)
			}
		})
	}
}

// --- TestFindCheapestTarget (whitebox) ---

func TestFindCheapestTarget(t *testing.T) {
	tests := []struct {
		name        string
		aVals       []float64
		bVals       []float64
		fromIdx     int
		target      stackID
		wantToIdx   int
		wantMaxCost int  // cost must be ≤ this; -1 means any valid target found
		wantNoMove  bool // toIdx == -1 means no valid target exists
	}{
		{
			// No element in B is smaller than A[0]=10, so wraps to maximum of B.
			name:        "no smaller target - wraps to maximum",
			aVals:       []float64{10, 1, 2},
			bVals:       []float64{3, 4, 5},
			fromIdx:     0,
			target:      stackB,
			wantToIdx:   2, // val=5 is maximum in B
			wantMaxCost: math.MaxInt,
		},
		{
			// A[0]=2; B has only one value greater (5 at index 0).
			name:        "single valid greater target",
			aVals:       []float64{2, 8},
			bVals:       []float64{5},
			fromIdx:     0,
			target:      stackB,
			wantToIdx:   0,
			wantMaxCost: math.MaxInt,
		},
		{
			// B[0]=5; A has values greater: find insertion point.
			name:        "find greater target in A",
			aVals:       []float64{1, 3, 7},
			bVals:       []float64{5, 2},
			fromIdx:     0,
			target:      stackA,
			wantToIdx:   2, // val=7 is smallest > 5
			wantMaxCost: math.MaxInt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := NewDoubleStack(tt.aVals...)
			for _, v := range tt.bVals {
				ds.B.PushBottom(v)
			}

			got := findCheapestTarget(ds, tt.fromIdx, tt.target)

			if tt.wantNoMove {
				if got.toIdx != -1 {
					t.Errorf("expected no valid target (toIdx=-1), got toIdx=%d cost=%d", got.toIdx, got.cost)
				}
				return
			}

			if got.toIdx != tt.wantToIdx {
				t.Errorf("toIdx = %d, want %d", got.toIdx, tt.wantToIdx)
			}
			if got.cost > tt.wantMaxCost {
				t.Errorf("cost = %d, want ≤ %d", got.cost, tt.wantMaxCost)
			}
		})
	}
}

// ============================================================================
// Regression Test References
// ============================================================================
//
// Bug 1: findExtremes never updated `extreme`, so every comparison held against
//        ±Inf and `indices` was reset on each iteration, always returning only
//        the last element's index.
//
// Bug 2: findTargets reset targetIndices on ANY differing value, so the last
//        qualifying value always won instead of the closest one.
//
// Bug 3: generateInstructions always emitted RA/RB for remaining individual
//        rotations even when the route (RRR) or element position (index > length/2)
//        required the reverse direction (RRA/RRB).
//
// Bug 4: stale fromIdx was reused across iterations when neither cost-improvement
//        condition fired (fromIdx not reset at loop start).
//
// Bug 5: loop guard was `top < bottom`, so the body never ran for a single-element
//        source stack (top==bottom==0), returning the invalid sentinel
//        {fromIdx:-1, toIdx:-1} and causing a panic in generateInstructions.
