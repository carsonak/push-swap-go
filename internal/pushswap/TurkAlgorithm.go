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

	maxIndices := findExtremes(sA, true)
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

// findStartIndexRotated returns the start index of a sorted rotated stack.
func findStartIndexRotated(stack *psStack, isAscending bool) (pivotIndex int) {
	topIdx := 0
	botIdx := stack.Len() - 1

	for topIdx < botIdx {
		midIdx := topIdx + (botIdx-topIdx)/2
		mid, _ := stack.Index(midIdx)
		top, _ := stack.Index(topIdx)
		bot, _ := stack.Index(botIdx)

		// Duplicated numbers means we cannot tell which side of `mid` is sorted,
		// therefore we shrink the search area to exclude the duplicated values.
		if mid == bot && top == bot {
			topIdx++
			botIdx--
			continue
		}

		if (!isAscending && mid < bot) || (isAscending && mid > bot) {
			topIdx = midIdx + 1 // The start is to the bot
		} else if (!isAscending && bot < mid) || (isAscending && bot > mid) {
			botIdx = midIdx // The start is to the top or mid
		}
	}

	return topIdx
}

// findTargetsInSortedRotated finds the indices of all the values immediately after the given
// reference in a rotated sorted stack. isAscending specifies the sort order of the stack.
//
// Complexity: O(log n + k), where k is the number of duplicate occurrences of the
// target value.
func findTargetsInSortedRotated(s *psStack, ref float64, greater bool, isAscending bool) []int {
	stackLen := s.Len()
	if stackLen == 0 {
		return nil
	}

	pivot := findStartIndexRotated(s, isAscending)
	searchDirection := isAscending == greater

	// searchSegmentCandidate binary-searches the sorted segment [top, bot] for the value
	// closest to ref on the desired side (smallest value > ref, or largest value < ref).
	searchSegmentCandidate := func(top, bot int) (bestVal float64, bestIdx int, found bool) {
		for top <= bot {
			mid := top + (bot-top)/2
			midVal, _ := s.Index(mid)

			qualifies := (greater && midVal > ref) || (!greater && midVal < ref)
			if qualifies && (!found || (greater && midVal < bestVal) || (!greater && midVal > bestVal)) {
				bestVal, bestIdx, found = midVal, mid, true
			}

			if qualifies != searchDirection {
				// Truth table for (qualifies != searchDirection):
				// isAscending=true,  greater=true  => searchDirection=true,  move right when mid <= ref
				// isAscending=true,  greater=false => searchDirection=false, move right when mid <  ref
				// isAscending=false, greater=true  => searchDirection=false, move right when mid >  ref
				// isAscending=false, greater=false => searchDirection=true,  move right when mid >= ref
				top = mid + 1
			} else {
				// Complement of the table above: move left in the opposite comparator cases.
				bot = mid - 1
			}
		}

		return
	}

	targetVal, targetIdx, found1 := searchSegmentCandidate(0, pivot-1)
	val2, idx2, found2 := searchSegmentCandidate(pivot, stackLen-1)

	if !found1 && !found2 {
		return nil
	}

	if !found1 || (found2 && ((greater && val2 < targetVal) || (!greater && val2 > targetVal))) {
		targetVal, targetIdx = val2, idx2
	}

	// Expand outward from targetIdx to collect all duplicate indices in O(k).
	// Duplicates in a sorted array are contiguous; they may wrap across the rotation
	// point, which modular arithmetic handles correctly.
	targetIndices := []int{targetIdx}

	for lower := (targetIdx + 1) % stackLen; lower != targetIdx; lower = (lower + 1) % stackLen {
		val, _ := s.Index(lower)
		if val != targetVal {
			break
		}

		targetIndices = append(targetIndices, lower)
	}

	if len(targetIndices) < stackLen { // Guard against an all-equal stack to prevent an infinite loop.
		for upper := (targetIdx - 1 + stackLen) % stackLen; upper != targetIdx; upper = (upper - 1 + stackLen) % stackLen {
			val, _ := s.Index(upper)
			if val != targetVal {
				break
			}

			targetIndices = append(targetIndices, upper)
		}
	}

	return targetIndices
}

// findGreaterTargetsInSortedRotated finds indices of all the smallest values greater than the
// given reference in a rotated ascending sorted stack.
func findGreaterTargetsInSortedRotated(s *psStack, ref float64) (targetIndices []int) {
	return findTargetsInSortedRotated(s, ref, true, true)
}

// findSmallerTargetsInSortedRotated finds indices of all the greatest values smaller than the
// given reference in a rotated descending sorted stack.
func findSmallerTargetsInSortedRotated(s *psStack, ref float64) (targetIndices []int) {
	return findTargetsInSortedRotated(s, ref, false, false)
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
func findCheapestTarget(stacks *DoubleStack, fromIdx int, target stackID) (cheapest moveCandidate) {
	to := &stacks.B
	from := &stacks.A
	if target == stackA {
		to = &stacks.A
		from = &stacks.B
	}

	reference, _ := from.Index(fromIdx)
	var targetIndices []int

	if target == stackB { // Values in B are sorted in descending order.
		targetIndices = findSmallerTargetsInSortedRotated(to, reference)
		if len(targetIndices) < 1 {
			targetIndices = findMaximums(to)
		}
	} else {
		targetIndices = findGreaterTargetsInSortedRotated(to, reference)
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

func findCheapestMove(stacks *DoubleStack, to stackID) (cheapest moveCandidate) {
	cheapest = moveCandidate{cost: math.MaxInt, fromIdx: -1, toIdx: -1, target: to}
	from := &stacks.A
	if to == stackA {
		from = &stacks.B
	}

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
			candidate := findCheapestTarget(stacks, fromIdx, to)

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
// element at the specified indices to their target positions.
func generateInstructions(stacks *DoubleStack, move moveCandidate) (instructions []Operation) {
	idxA := move.toIdx
	idxB := move.fromIdx

	if move.target == stackB {
		idxA = move.fromIdx
		idxB = move.toIdx
	}

	// Calculate rotations for each stack.
	aRotations := idxA
	opA := RA
	if move.route == RRR || idxA > stacks.A.Len()/2 {
		aRotations = stacks.A.Len() - idxA
		opA = RRA
	}

	bRotations := idxB
	opB := RB
	if move.route == RRR || idxB > stacks.B.Len()/2 {
		bRotations = stacks.B.Len() - idxB
		opB = RRB
	}

	// Handle simultaneous rotations.
	if move.route == RRR || move.route == RR {
		sharedRotations := min(aRotations, bRotations)
		instructions = slices.Repeat([]Operation{move.route}, sharedRotations)
		aRotations -= sharedRotations
		bRotations -= sharedRotations

		// Reset operation types for remaining individual rotations.
		if move.route == RRR {
			opA = RRA
			opB = RRB
		} else {
			opA = RA
			opB = RB
		}
	}

	instructions = append(instructions, slices.Repeat([]Operation{opA}, aRotations)...)
	instructions = append(instructions, slices.Repeat([]Operation{opB}, bRotations)...)
	return instructions
}

func executeInstructions(stacks *DoubleStack, instructions []Operation) {
	for _, op := range instructions {
		switch op {
		case PA:
			stacks.PushToA()
		case PB:
			stacks.PushToB()
		case RA:
			stacks.RotateA()
		case RB:
			stacks.RotateB()
		case RR:
			stacks.RRotate()
		case RRA:
			stacks.ReverseRotateA()
		case RRB:
			stacks.ReverseRotateB()
		case RRR:
			stacks.RReverseRotate()
		case SA:
			stacks.SwapA()
		case SB:
			stacks.SwapB()
		case SS:
			stacks.SSwap()
		}
	}
}

func TurkAlgorithm(nums []float64) []Operation {
	if slices.IsSorted(nums) {
		return nil
	}

	stacks := NewDoubleStack(nums...)

	if len(nums) <= 3 {
		return sortLast3(&stacks.A)
	}

	instructions := []Operation{stacks.PushToB(), stacks.PushToB()}

	for stacks.A.Len() > 3 {
		move := findCheapestMove(stacks, stackB)
		ops := append(generateInstructions(stacks, move), PB)

		executeInstructions(stacks, ops)
		instructions = append(instructions, ops...)
	}

	instructions = append(instructions, sortLast3(&stacks.A)...)

	for stacks.B.Len() > 0 {
		move := findCheapestMove(stacks, stackA)
		ops := append(generateInstructions(stacks, move), PA)

		executeInstructions(stacks, ops)
		instructions = append(instructions, ops...)
	}

	// Rotate stack A until minimum comes to the top.
	minIndices := findMinimums(&stacks.A)
	slices.Sort(minIndices)
	minIdx := minIndices[0]

	var rotations []Operation
	if minIdx < stacks.A.Len()-minIdx {
		rotations = slices.Repeat([]Operation{RA}, minIdx)
	} else {
		rotations = slices.Repeat([]Operation{RRA}, stacks.A.Len()-minIdx)
	}

	return append(instructions, rotations...)
}
