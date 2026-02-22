package pushswap

import stack "push-swap-go/internal/dllStack"

// DStack is a pair of stacks for the push-swap program.
type DStack struct {
	A stack.Stack[float64]
	B stack.Stack[float64]
}

// New initialises a `DStack` with an optional list of numbers.
func New(nums ...float64) *DStack {
	return &DStack{
		A: *stack.New(nums...),
		B: *stack.NewWithCapacity[float64](len(nums)),
	}
}

func (ds *DStack) PushToA() Operation {
	val, success := ds.B.Pop()
	if !success {
		return Invalid
	}

	ds.A.Push(val)
	return PA
}

func (ds *DStack) PushToB() Operation {
	val, success := ds.A.Pop()
	if !success {
		return Invalid
	}

	ds.B.Push(val)
	return PB
}

func (ds *DStack) SwapA() Operation {
	ds.A.Swap()
	return SA
}

func (ds *DStack) SwapB() Operation {
	ds.B.Swap()
	return SB
}

func (ds *DStack) SSwap() Operation {
	ds.A.Swap()
	ds.B.Swap()
	return SS
}

func (ds *DStack) RotateA() Operation {
	ds.A.Rotate()
	return RA
}

func (ds *DStack) RotateB() Operation {
	ds.B.Rotate()
	return RB
}

func (ds *DStack) RRotate() Operation {
	ds.A.Rotate()
	ds.B.Rotate()
	return RR
}

func (ds *DStack) ReverseRotateA() Operation {
	ds.A.ReverseRotate()
	return RRA
}

func (ds *DStack) ReverseRotateB() Operation {
	ds.B.ReverseRotate()
	return RRB
}

func (ds *DStack) RReverseRotate() Operation {
	ds.A.ReverseRotate()
	ds.B.ReverseRotate()
	return RRR
}
