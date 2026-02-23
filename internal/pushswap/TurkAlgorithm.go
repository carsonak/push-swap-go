package pushswap

import (
	"math"
	"slices"
)

type stackID byte

const (
	stackA stackID = iota
	stackB
)

type moveCandidate struct {
	fromIdx int
	toIdx   int
	cost    int
	target  stackID
	route   Operation
}

// findExtremes returns the indices of all (in-case of duplicates)
// the smallest/largest values in the stack.
func findExtremes(s *psStack, findMax bool) (indices []int) {
	extreme := math.Inf(1)

	if findMax {
		extreme = math.Inf(-1)
	}

	for i, n := range s.All() {
		if (findMax && n > extreme) || (!findMax && n < extreme) {
			indices = []int{i}
		} else if n == extreme {
			indices = append(indices, i)
		}
	}

	return indices
}

// findMaximums returns the indices of all (in-case of duplicates)
// the greatest values in the stack.
func findMaximums(s *psStack) []int {
	return findExtremes(s, true)
}

// findMinimums returns the indices of all (in-case of duplicates)
// the smallest values in the stack.
func findMinimums(s *psStack) []int {
	return findExtremes(s, false)
}

// sortLast3 sorts a stack with at-most 3 values.
func sortLast3(sA *psStack) (instructions []Operation) {
	if sA.Len() < 2 || sA.Len() > 3 {
		return nil
	}

	maxIndices := findMaximums(sA)
	maxIdx := maxIndices[0]

	if sA.Len() > 2 && len(maxIndices) < 3 {
		// Rotate the max value to the bottom
		switch maxIdx {
		case 0:
			sA.Rotate()
			instructions = append(instructions, RA)
		case 1:
			sA.ReverseRotate()
			instructions = append(instructions, RRA)
		}
	}

	n0, _ := sA.Index(0)
	n1, _ := sA.Index(1)

	if n0 > n1 {
		sA.Swap()
		instructions = append(instructions, SA)
	}

	return instructions
}

func identifyStack(ss *DoubleStack, id stackID) *psStack {
	if id == stackA {
		return &ss.A
	}

	return &ss.B
}

func findCheapestMove(ss *DoubleStack, to stackID) moveCandidate {
	panic("Not Implemented")
}

func generateInstructions(ss *DoubleStack, move moveCandidate) (instructions []Operation) {
	// Identify target stack.

	idxA := move.toIdx
	idxB := move.fromIdx

	if move.target == stackB {
		idxA = move.fromIdx
		idxB = move.toIdx
	}

	// Calculate number of rotations on each stack.

	ARotations := idxA
	if move.route == RRR || idxA > ss.A.Len()/2 {
		ARotations = ss.A.Len() - idxA
	}

	BRotations := idxB
	if move.route == RRR || idxB > ss.B.Len()/2 {
		BRotations = ss.B.Len() - idxB
	}

	// Generate operations.

	var ops []Operation

	if move.route == RRR || move.route == RR {
		if ARotations > BRotations {
			ARotations -= BRotations
			ops = slices.Repeat([]Operation{move.route}, BRotations)
			BRotations = 0
		} else {
			BRotations -= ARotations
			ops = slices.Repeat([]Operation{move.route}, ARotations)
			ARotations = 0
		}

		instructions = append(instructions, ops...)
	}

	ops = slices.Repeat([]Operation{RA}, ARotations)
	instructions = append(instructions, ops...)
	ops = slices.Repeat([]Operation{RB}, BRotations)
	return append(instructions, ops...)
}

func executeInstructions(ss *DoubleStack, instructions []Operation) {
	for _, op := range instructions {
		switch op {
		case PA:
			ss.PushToA()
		case PB:
			ss.PushToB()
		case RA:
			ss.RotateA()
		case RB:
			ss.RotateB()
		case RR:
			ss.RRotate()
		case RRA:
			ss.ReverseRotateA()
		case RRB:
			ss.ReverseRotateB()
		case RRR:
			ss.RReverseRotate()
		case SA:
			ss.SwapA()
		case SB:
			ss.SwapB()
		case SS:
			ss.SSwap()
		}
	}
}

func TurkAlgorithm(nums []float64) []Operation {
	if slices.IsSorted(nums) {
		return nil
	}

	ss := New(nums...)

	if len(nums) <= 3 {
		return sortLast3(&ss.A)
	}

	instructions := []Operation{ss.PushToB(), ss.PushToB()}

	for ss.A.Len() > 3 {
		move := findCheapestMove(ss, stackB)
		ops := append(generateInstructions(ss, move), PB)

		executeInstructions(ss, ops)
		instructions = append(instructions, ops...)
	}

	instructions = append(instructions, sortLast3(&ss.A)...)

	for ss.B.Len() > 0 {
		move := findCheapestMove(ss, stackA)
		ops := append(generateInstructions(ss, move), PA)

		executeInstructions(ss, ops)
		instructions = append(instructions, ops...)
	}

	// Rotate stackA till minimum comes to the top
	var rotations []Operation
	mins := findMinimums(&ss.A)

	slices.Sort(mins)

	if mins[0] < (ss.A.Len() - mins[0]) {
		rotations = slices.Repeat([]Operation{RA}, mins[0])
	} else {
		rotations = slices.Repeat([]Operation{RRA}, ss.A.Len()-mins[0])
	}

	return append(instructions, rotations...)
}
