package benchmarks

/*
PERFORMANCE PROFILING WITH PERF STAT

1) Basic cache stats:
perf stat -e cache-references,cache-misses,L1-dcache-loads,L1-dcache-load-misses,LLC-loads,LLC-load-misses \
  go test -run=^$ -bench=. -benchtime=5s ./internal/benchmarks

2) Detailed cache hierarchy:
perf stat -e L1-dcache-loads,L1-dcache-load-misses,L1-dcache-stores,L1-icache-load-misses,LLC-loads,LLC-load-misses,LLC-stores,LLC-store-misses,dTLB-loads,dTLB-load-misses \
  go test -run=^$ -bench=BenchmarkTurkAlgorithm_MassiveScale/50000 -benchtime=3s ./internal/benchmarks

3) Instructions + branches:
perf stat -e instructions,cycles,branches,branch-misses,cache-references,cache-misses \
  go test -run=^$ -bench=. -benchtime=10s ./internal/benchmarks
*/

import (
	"cmp"
	"context"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"push-swap-go/internal/pushswap"
)

const benchmarkIterationTimeout = 10 * time.Second

type AlgorithmFunc[T cmp.Ordered] func([]T) []pushswap.Operation

func runWithTimeout[T cmp.Ordered](ctx context.Context, algo AlgorithmFunc[T], data []T) (int, bool) {
	type result struct{ count int }
	resultChan := make(chan result, 1)

	go func() {
		dataCopy := make([]T, len(data))
		copy(dataCopy, data)
		ops := algo(dataCopy)
		resultChan <- result{count: len(ops)}
	}()

	select {
	case res := <-resultChan:
		return res.count, false
	case <-ctx.Done():
		return -1, true
	}
}

func generateRandomInts(n, min, max int, seed int64) []int {
	rng := rand.New(rand.NewSource(seed))
	data := make([]int, n)

	for i := range data {
		data[i] = rng.Intn(max-min+1) + min
	}

	return data
}

func generateRandomFloats(n int, min, max float64, seed int64) []float64 {
	rng := rand.New(rand.NewSource(seed))
	data := make([]float64, n)

	for i := range data {
		data[i] = min + rng.Float64()*(max-min)
	}

	return data
}

func generateFloatsWithDuplicates(n, dupPercent int, seed int64) []float64 {
	rng := rand.New(rand.NewSource(seed))
	data := make([]float64, n)
	duplicateValue := rng.Float64() * 10000.0
	numDuplicates := (n * dupPercent) / 100

	for i := 0; i < numDuplicates; i++ {
		data[i] = duplicateValue
	}

	for i := numDuplicates; i < n; i++ {
		data[i] = rng.Float64() * 10000.0
	}

	rng.Shuffle(n, func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})

	return data
}

func generateNearlySorted(n, clusterSize int, position string, seed int64) []float64 {
	rng := rand.New(rand.NewSource(seed))
	data := make([]float64, n)

	for i := range data {
		data[i] = float64(i)
	}

	cluster := make([]float64, clusterSize)

	for i := range cluster {
		cluster[i] = rng.Float64() * float64(n)
	}

	insertPos := 0
	switch position {
	case "middle":
		insertPos = n / 2
	case "bottom":
		insertPos = n - clusterSize
	}

	copy(data[insertPos:insertPos+clusterSize], cluster)
	return data
}

func generateTinyFloats(n int, seed int64) []float64 {
	rng := rand.New(rand.NewSource(seed))
	data := make([]float64, n)

	for i := range data {
		data[i] = 0.0001 + rng.Float64()*0.0008
	}

	return data
}

func generateMassiveFloats(n int, seed int64) []float64 {
	rng := rand.New(rand.NewSource(seed))
	data := make([]float64, n)

	for i := range data {
		data[i] = math.MaxInt32 + rng.Float64()*math.MaxInt32
	}

	return data
}

func reportInstructions(b *testing.B, total, completed int) {
	if completed == 0 {
		return
	}

	avg := float64(total) / float64(completed)
	b.ReportMetric(avg, "inst/op")
}

func runTimedBenchmark[T cmp.Ordered](b *testing.B, datasets [][]T, algo AlgorithmFunc[T]) {
	b.ResetTimer()
	totalInstructions := 0
	completed := 0
	timedOut := false

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), benchmarkIterationTimeout)
		count, timeoutReached := runWithTimeout(ctx, algo, datasets[i])
		cancel()

		if timeoutReached {
			timedOut = true
			break
		}

		totalInstructions += count
		completed++
	}

	reportInstructions(b, totalInstructions, completed)
	if timedOut {
		b.ReportMetric(1, "timeout")
	} else {
		b.ReportMetric(0, "timeout")
	}
}

func BenchmarkTurkAlgorithm_MassiveScale(b *testing.B) {
	sizes := []int{100, 1000, 3000, 5000, 10000, 50000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			datasets := make([][]int, b.N)
			for i := 0; i < b.N; i++ {
				datasets[i] = generateRandomInts(size, -100000, 100000, int64(i))
			}

			runTimedBenchmark(b, datasets, pushswap.TurkAlgorithm[int])
		})
	}
}

func BenchmarkTurkAlgorithm_StandardFloats(b *testing.B) {
	sizes := []int{500, 1000, 1500, 1900}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			datasets := make([][]float64, b.N)
			for i := 0; i < b.N; i++ {
				datasets[i] = generateRandomFloats(size, -10000, 10000, int64(i))
			}

			runTimedBenchmark(b, datasets, pushswap.TurkAlgorithm[float64])
		})
	}
}

func BenchmarkTurkAlgorithm_HighDensityDuplicates(b *testing.B) {
	sizes := []int{500, 1000, 1500, 1900}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			datasets := make([][]float64, b.N)
			for i := 0; i < b.N; i++ {
				datasets[i] = generateFloatsWithDuplicates(size, 50, int64(i))
			}

			runTimedBenchmark(b, datasets, pushswap.TurkAlgorithm[float64])
		})
	}
}

func BenchmarkTurkAlgorithm_NearlySorted_TopHeavy(b *testing.B) {
	benchmarkNearlySorted(b, "top")
}

func BenchmarkTurkAlgorithm_NearlySorted_MiddleHeavy(b *testing.B) {
	benchmarkNearlySorted(b, "middle")
}

func BenchmarkTurkAlgorithm_NearlySorted_BottomHeavy(b *testing.B) {
	benchmarkNearlySorted(b, "bottom")
}

func benchmarkNearlySorted(b *testing.B, position string) {
	sizes := []int{500, 750, 1000}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			datasets := make([][]float64, b.N)
			for i := 0; i < b.N; i++ {
				datasets[i] = generateNearlySorted(size, 10, position, int64(i))
			}

			runTimedBenchmark(b, datasets, pushswap.TurkAlgorithm[float64])
		})
	}
}

func BenchmarkTurkAlgorithm_TinyFloats(b *testing.B) {
	sizes := []int{500, 1000, 1500, 1900}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			datasets := make([][]float64, b.N)
			for i := 0; i < b.N; i++ {
				datasets[i] = generateTinyFloats(size, int64(i))
			}

			runTimedBenchmark(b, datasets, pushswap.TurkAlgorithm[float64])
		})
	}
}

func BenchmarkTurkAlgorithm_MassiveFloats(b *testing.B) {
	sizes := []int{500, 1000, 1500, 1900}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			datasets := make([][]float64, b.N)
			for i := 0; i < b.N; i++ {
				datasets[i] = generateMassiveFloats(size, int64(i))
			}

			runTimedBenchmark(b, datasets, pushswap.TurkAlgorithm[float64])
		})
	}
}

func BenchmarkTurkAlgorithm_AllScenarios(b *testing.B) {
	scenarios := []struct {
		name      string
		generator func(seed int64) []float64
	}{
		{"Random_1000", func(seed int64) []float64 { return generateRandomFloats(1000, -10000, 10000, seed) }},
		{"Duplicates_1000", func(seed int64) []float64 { return generateFloatsWithDuplicates(1000, 50, seed) }},
		{"NearlySorted_1000", func(seed int64) []float64 { return generateNearlySorted(1000, 10, "middle", seed) }},
		{"TinyFloats_1000", func(seed int64) []float64 { return generateTinyFloats(1000, seed) }},
		{"MassiveFloats_1000", func(seed int64) []float64 { return generateMassiveFloats(1000, seed) }},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			datasets := make([][]float64, b.N)
			for i := 0; i < b.N; i++ {
				datasets[i] = scenario.generator(int64(i))
			}

			runTimedBenchmark(b, datasets, pushswap.TurkAlgorithm[float64])
		})
	}
}
