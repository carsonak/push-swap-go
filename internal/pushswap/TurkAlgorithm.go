package pushswap

import (
	"math"
	"slices"
)

type stackID bool

const (
	stackA stackID = false
	stackB stackID = true
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
			extreme = n
		} else if n == extreme {
			indices = append(indices, i)
		}
	}

	return indices
}

// findMaximums returns the indices of all (in-case of duplicates)
// the greatest values in the stack.
func findMaximums(s *psStack) (indices []int) {
	return findExtremes(s, true)
}

// findMinimums returns the indices of all (in-case of duplicates)
// the smallest values in the stack.
func findMinimums(s *psStack) (indices []int) {
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

// identifyStack returns pointer to the stack identified by the given `stackID`.
func identifyStack(ss *DoubleStack, id stackID) *psStack {
	if id == stackA {
		return &ss.A
	}

	return &ss.B
}

// findTargets finds the indices of all the values immediately after the given reference.
func findTargets(s *psStack, ref float64, greater bool) (targetIndices []int) {
	targetVal := math.Inf(1)

	if !greater {
		targetVal = math.Inf(-1)
	}

	// TODO: Apply the following optimisation.
	// OPTIMISATION: We can assume that the target stack `s` is always sorted, therefore
	// we can apply binary search to find the targets.
	for i := range s.Len() {
		val, _ := s.Index(i)

		if (greater && val > ref) || (!greater && val < ref) {
			if val == targetVal {
				targetIndices = append(targetIndices, i)
			} else if (greater && val < targetVal) || (!greater && val > targetVal) {
				targetVal = val
				targetIndices = []int{i}
			}
		}
	}

	return targetIndices
}

// findGreaterTargets finds indices of all the smallest values greater than the
// given reference.
func findGreaterTargets(s *psStack, ref float64) (targetIndices []int) {
	return findTargets(s, ref, true)
}

// findSmallerTargets finds indices of all the greatest values smaller than the
// given reference.
func findSmallerTargets(s *psStack, ref float64) (targetIndices []int) {
	return findTargets(s, ref, false)
}

// shortestRouteToTop calculates the shortest route to move a value at the given
// index on a stack of the given length.
func shortestRouteToTop(index, length int) int {
	if index > length/2 {
		return length - index
	}

	return index
}

// findCheapestTarget returns the cheapest move candidate of the reference index
// `fromIdx` in the target stack.
func findCheapestTarget(ss *DoubleStack, fromIdx int, target stackID) (cheapest moveCandidate) {
	to := identifyStack(ss, target)
	from := identifyStack(ss, !target)
	reference, _ := from.Index(fromIdx)
	var targetIndices []int

	if target == stackB { // Values in B are sorted in descending order.
		targetIndices = findSmallerTargets(to, reference)
		if len(targetIndices) < 1 {
			targetIndices = findMaximums(to)
		}
	} else {
		targetIndices = findGreaterTargets(to, reference)
		if len(targetIndices) < 1 {
			targetIndices = findMinimums(to)
		}
	}

	// OPTIMISATION: For duplicated values, we just pick the indices closest
	// to the ends of the stack.
	if lenTI := len(targetIndices); lenTI > 1 {
		slices.Sort(targetIndices)
		targetIndices = []int{targetIndices[0], targetIndices[lenTI-1]}
	}

	cheapest = moveCandidate{
		fromIdx: fromIdx,
		toIdx:   -1,
		cost:    math.MaxInt,
		target:  target,
	}

	for _, toIndex := range targetIndices {
		rotateCost := cheapest.fromIdx + toIndex
		reverseRotateCost := (from.Len() - cheapest.fromIdx) + (to.Len() - toIndex)
		separateCost := shortestRouteToTop(cheapest.fromIdx, from.Len()) + shortestRouteToTop(toIndex, to.Len())
		minCost := min(rotateCost, reverseRotateCost)

		minCost = min(minCost, separateCost)

		if minCost < cheapest.cost {
			cheapest.toIdx = toIndex
			cheapest.cost = minCost
			switch minCost {
			case rotateCost:
				cheapest.route = RR
			case reverseRotateCost:
				cheapest.route = RRR
			default:
				var separate Operation
				cheapest.route = separate
			}
		}
	}

	return cheapest
}

func findCheapestMove(ss *DoubleStack, to stackID) (cheapest moveCandidate) {
	cheapest = moveCandidate{cost: math.MaxInt, fromIdx: -1, toIdx: -1, target: to}
	from := identifyStack(ss, !to)

	// OPTIMISATION: Iterate from the ends working inwards, as cheaper moves are
	// usually at the ends.
	top, bottom := 0, from.Len()-1
	for top <= bottom {
		fromIdx := -1

		if top < cheapest.cost {
			fromIdx = top
		} else if from.Len()-bottom < cheapest.cost {
			fromIdx = bottom
		}

		if fromIdx >= 0 {
			candidate := findCheapestTarget(ss, fromIdx, to)

			if candidate.cost < cheapest.cost || cheapest.fromIdx < 0 {
				cheapest = candidate
			}
		}

		// Break early since the deeper we go the higher the costs.
		if cheapest.cost < 3 {
			break
		}

		top++
		bottom--
	}

	return cheapest
}

// generateInstructions returns the shortest list of instructions required to move the
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
	opA := RA

	if move.route == RRR || idxA > ss.A.Len()/2 {
		ARotations = ss.A.Len() - idxA
		opA = RRA
	}

	BRotations := idxB
	opB := RB

	if move.route == RRR || idxB > ss.B.Len()/2 {
		BRotations = ss.B.Len() - idxB
		opB = RRB
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

		opA = RA
		opB = RB
		if move.route == RRR {
			opA = RRA
			opB = RRB
		}

		instructions = append(instructions, ops...)
	}

	ops = slices.Repeat([]Operation{opA}, ARotations)
	instructions = append(instructions, ops...)
	ops = slices.Repeat([]Operation{opB}, BRotations)
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

	ss := NewDoubleStack(nums...)

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
