package pushswap

import stack "push-swap-go/internal/dllStack"

type psStack = stack.Stack[float64]

// DoubleStack is a pair of stacks for the push-swap program.
type DoubleStack struct {
	A psStack
	B psStack
}

// NewDoubleStack initialises a `DStack` with an optional list of numbers.
func NewDoubleStack(nums ...float64) *DoubleStack {
	return &DoubleStack{
		A: *stack.New(nums...),
		B: *stack.NewWithCapacity[float64](len(nums)),
	}
}

func (ds *DoubleStack) PushToA() Operation {
	val, success := ds.B.Pop()
	if !success {
		return Invalid
	}

	ds.A.Push(val)
	return PA
}

func (ds *DoubleStack) PushToB() Operation {
	val, success := ds.A.Pop()
	if !success {
		return Invalid
	}

	ds.B.Push(val)
	return PB
}

func (ds *DoubleStack) SwapA() Operation {
	ds.A.Swap()
	return SA
}

func (ds *DoubleStack) SwapB() Operation {
	ds.B.Swap()
	return SB
}

func (ds *DoubleStack) SSwap() Operation {
	ds.A.Swap()
	ds.B.Swap()
	return SS
}

func (ds *DoubleStack) RotateA() Operation {
	ds.A.Rotate()
	return RA
}

func (ds *DoubleStack) RotateB() Operation {
	ds.B.Rotate()
	return RB
}

func (ds *DoubleStack) RRotate() Operation {
	ds.A.Rotate()
	ds.B.Rotate()
	return RR
}

func (ds *DoubleStack) ReverseRotateA() Operation {
	ds.A.ReverseRotate()
	return RRA
}

func (ds *DoubleStack) ReverseRotateB() Operation {
	ds.B.ReverseRotate()
	return RRB
}

func (ds *DoubleStack) RReverseRotate() Operation {
	ds.A.ReverseRotate()
	ds.B.ReverseRotate()
	return RRR
}
