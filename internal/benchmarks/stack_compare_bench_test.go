package benchmarks

import (
	"math/rand"
	"testing"

	dllStack "push-swap-go/internal/dll_stack"
	sliceStack "push-swap-go/internal/slice_stack"
)

// Stacker interface allows us to run the exact same benchmark code on both types
type Stacker[T any] interface {
	Push(T)
	Pop() (T, bool)
	Rotate()
	ReverseRotate()
}

type opType int

const (
	opPush opType = iota
	opPop
	opRotate
	opRevRotate
)

// --- Scenario 1: Pre-allocation impact ---

func BenchmarkPrealloc_Slice(b *testing.B) {
	b.Run("No_Pre-allocation", func(b *testing.B) {
		for range b.N {
			s := sliceStack.New[int]()

			for j := range 1000 {
				s.Push(j)
			}
		}
	})

	b.Run("With_Pre-allocation", func(b *testing.B) {
		for range b.N {
			s := sliceStack.NewWithCapacity[int](1000)

			for j := range 1000 {
				s.Push(j)
			}
		}
	})
}

func BenchmarkPrealloc_DLL(b *testing.B) {
	b.Run("No_Pre-allocation", func(b *testing.B) {
		for range b.N {
			s := dllStack.New[int]()

			for j := range 1000 {
				s.Push(j)
			}
		}
	})

	b.Run("With_Pre-allocation", func(b *testing.B) {
		for range b.N {
			s := dllStack.NewWithCapacity[int](1000)

			for j := range 1000 {
				s.Push(j)
			}
		}
	})
}

// Scenario 2 & 3 Helper for executing operation lists
func runOpsBenchmark(b *testing.B, ops []opType, s Stacker[int]) {
	b.ResetTimer()
	for range b.N {
		for _, op := range ops {
			switch op {
			case opPush:
				s.Push(1)
			case opPop:
				s.Pop()
			case opRotate:
				s.Rotate()
			case opRevRotate:
				s.ReverseRotate()
			}
		}
	}
}

// generateOps creates a slice of operations based on given probabilities
func generateOps(count int, pushW, popW, rotW, revRotW int) []opType {
	total := pushW + popW + rotW + revRotW
	ops := make([]opType, count)

	for i := range count {
		r := rand.Intn(total)
		if r < pushW {
			ops[i] = opPush
		} else if r < pushW+popW {
			ops[i] = opPop
		} else if r < pushW+popW+rotW {
			ops[i] = opRotate
		} else {
			ops[i] = opRevRotate
		}
	}

	return ops
}

// --- Scenario 2: Bench different ratios of operations on the stacks ---

func BenchmarkRatios(b *testing.B) {
	scenarios := []struct {
		name string
		ops  []opType
	}{
		{"More_Pushes", generateOps(10000, 70, 30, 0, 0)},
		{"More_Pops", generateOps(10000, 30, 70, 0, 0)},
		{"Equal_PushPop", generateOps(10000, 50, 50, 0, 0)},
		{"Heavy_Rotations", generateOps(10000, 40, 40, 10, 10)},
	}

	for _, sc := range scenarios {
		b.Run("Slice_"+sc.name, func(b *testing.B) {
			s := sliceStack.NewWithCapacity[int](10000)

			runOpsBenchmark(b, sc.ops, s)
		})

		b.Run("DLL_"+sc.name, func(b *testing.B) {
			s := dllStack.NewWithCapacity[int](10000)

			runOpsBenchmark(b, sc.ops, s)
		})
	}
}

// --- Scenario 3: Stack Frame Simulation (LIFO bursts) ---

func BenchmarkStackFrames(b *testing.B) {
	// Simulate calling into functions (pushes) and returning (pops)
	ops := make([]opType, 0, 10000)

	for range 100 { // 100 frames
		depth := rand.Intn(50) + 10 // push 10-60 variables per frame

		for range depth {
			ops = append(ops, opPush)
		}

		for range depth {
			ops = append(ops, opPop)
		}
	}

	b.Run("Slice_StackFrames", func(b *testing.B) {
		s := sliceStack.NewWithCapacity[int](10000)

		runOpsBenchmark(b, ops, s)
	})

	b.Run("DLL_StackFrames", func(b *testing.B) {
		s := dllStack.NewWithCapacity[int](10000)

		runOpsBenchmark(b, ops, s)
	})
}
