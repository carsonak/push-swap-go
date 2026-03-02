package pushswap

import (
	"cmp"

	stack "push-swap-go/internal/dllStack"
)

// DoubleStack is a pair of stacks for the push-swap program.
// T must satisfy cmp.Ordered so that elements can be compared with <, >, and ==.
type DoubleStack[T cmp.Ordered] struct {
	A stack.Stack[T]
	B stack.Stack[T]
}

// NewDoubleStack initialises a DoubleStack with an optional list of values.
func NewDoubleStack[T cmp.Ordered](nums ...T) *DoubleStack[T] {
	return &DoubleStack[T]{
		A: *stack.New(nums...),
		B: *stack.NewWithCapacity[T](len(nums)),
	}
}

func (ds *DoubleStack[T]) PushToA() Operation {
	val, success := ds.B.Pop()
	if !success {
		return Invalid
	}

	ds.A.Push(val)
	return PA
}

func (ds *DoubleStack[T]) PushToB() Operation {
	val, success := ds.A.Pop()
	if !success {
		return Invalid
	}

	ds.B.Push(val)
	return PB
}

func (ds *DoubleStack[T]) SwapA() Operation {
	ds.A.Swap()
	return SA
}

func (ds *DoubleStack[T]) SwapB() Operation {
	ds.B.Swap()
	return SB
}

func (ds *DoubleStack[T]) SSwap() Operation {
	ds.A.Swap()
	ds.B.Swap()
	return SS
}

func (ds *DoubleStack[T]) RotateA() Operation {
	ds.A.Rotate()
	return RA
}

func (ds *DoubleStack[T]) RotateB() Operation {
	ds.B.Rotate()
	return RB
}

func (ds *DoubleStack[T]) RRotate() Operation {
	ds.A.Rotate()
	ds.B.Rotate()
	return RR
}

func (ds *DoubleStack[T]) ReverseRotateA() Operation {
	ds.A.ReverseRotate()
	return RRA
}

func (ds *DoubleStack[T]) ReverseRotateB() Operation {
	ds.B.ReverseRotate()
	return RRB
}

func (ds *DoubleStack[T]) RReverseRotate() Operation {
	ds.A.ReverseRotate()
	ds.B.ReverseRotate()
	return RRR
}

// func (ds *DoubleStack[T]) String() string {
// 	var output strings.Builder

// 	output.WriteString(fmt.Sprintln("{"))
// 	output.WriteString(fmt.Sprintln("\tA:", ds.A))
// 	output.WriteString(fmt.Sprintln("\tB:", ds.B))
// 	output.WriteString(fmt.Sprintln("}"))
// 	return output.String()
// }

func (ds *DoubleStack[T]) ExecuteInstructions(instructions []Operation) {
	for _, op := range instructions {
		switch op {
		case PA:
			ds.PushToA()
		case PB:
			ds.PushToB()
		case RA:
			ds.RotateA()
		case RB:
			ds.RotateB()
		case RR:
			ds.RRotate()
		case RRA:
			ds.ReverseRotateA()
		case RRB:
			ds.ReverseRotateB()
		case RRR:
			ds.RReverseRotate()
		case SA:
			ds.SwapA()
		case SB:
			ds.SwapB()
		case SS:
			ds.SSwap()
		}
	}
}
