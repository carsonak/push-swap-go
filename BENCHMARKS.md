# Push-Swap Benchmark Suite

This suite benchmarks algorithm and stack behavior for your custom push-swap implementation.
Algorithm benchmarks now live in `internal/benchmarks/algorithm_bench_test.go` and are the default Make target.

## Quick Start

```bash
# default target (same as make bench-algo)
make

# run algorithm benchmarks
make bench-algo

# run stack benchmarks
make bench-stack

# run algorithm benchmarks under perf
make perf-bench-algo

# run stack benchmarks under perf
make perf-bench-stack
```

## Current Scenarios

### 1) Massive Scale (`BenchmarkTurkAlgorithm_MassiveScale`)

- Sizes: `100, 1000, 3000, 5000, 10000, 50000`
- Data: random ints in `[-100000, 100000]`
- Timeout: **10 seconds per iteration**
- Metric: `inst/op` via `b.ReportMetric(...)`

### 2) Floating-Point + Duplicates

- `BenchmarkTurkAlgorithm_StandardFloats`
- `BenchmarkTurkAlgorithm_HighDensityDuplicates`
- Sizes for both: `500, 1000, 1500, 1900` (kept below 2000)

### 3) Nearly Sorted Variants

- `BenchmarkTurkAlgorithm_NearlySorted_TopHeavy`
- `BenchmarkTurkAlgorithm_NearlySorted_MiddleHeavy`
- `BenchmarkTurkAlgorithm_NearlySorted_BottomHeavy`
- Sizes: `500, 750, 1000`
- Cluster size: `10`

### 4) Bit-Depth / Precision Stress

- `BenchmarkTurkAlgorithm_TinyFloats` (fractional range around `0.0001–0.0009`)
- `BenchmarkTurkAlgorithm_MassiveFloats` (values above 32-bit integer range)
- Sizes for both: `500, 1000, 1500, 1900` (kept below 2000)

### 5) Mixed Scenario

- `BenchmarkTurkAlgorithm_AllScenarios`
- Includes representative `1000`-size cases for random, duplicates, nearly-sorted, tiny, and massive floats.

## Metrics

In addition to `ns/op`, `B/op`, and `allocs/op`, benchmarks report:

- `inst/op`: average number of generated push-swap instructions.

Example:

```text
BenchmarkTurkAlgorithm_MassiveScale/1000-8    1    433823935 ns/op    19427 inst/op
```

## Perf Usage

### Make Targets

```bash
make perf-setup
make perf-bench-algo
make perf-bench-stack
make perf-restore
```

Both perf targets write each sub-benchmark result to its own file using the benchmark name, with `/` replaced by `-` and `.txt` suffix.

### Manual perf Commands

```bash
# Basic cache stats
perf stat -e cache-references,cache-misses,L1-dcache-loads,L1-dcache-load-misses,LLC-loads,LLC-load-misses \
  go test -run=^$ -bench=. -benchtime=5s ./internal/benchmarks

# Detailed cache hierarchy
perf stat -e L1-dcache-loads,L1-dcache-load-misses,L1-dcache-stores,L1-icache-load-misses,LLC-loads,LLC-load-misses,LLC-stores,LLC-store-misses,dTLB-loads,dTLB-load-misses \
  go test -run=^$ -bench=BenchmarkTurkAlgorithm_MassiveScale/50000 -benchtime=3s ./internal/benchmarks

# Instructions + branches
perf stat -e instructions,cycles,branches,branch-misses,cache-references,cache-misses \
  go test -run=^$ -bench=. -benchtime=10s ./internal/benchmarks
```

## Makefile Target Map

- `bench-algo`: algorithm benchmarks (default/first target)
- `bench-stack`: stack comparison benchmarks
- `perf-bench-algo`: algorithm benchmarks with perf
- `perf-bench-stack`: stack benchmarks with perf
- `perf-setup` / `perf-restore`: adjust and restore `perf_event_paranoid`

## Files

```text
internal/benchmarks/
├── algorithm_bench_test.go
└── stack_compare_bench_test.go

internal/pushswap/
├── TurkAlgorithm.go
├── operations.go
└── ...
```

## Adding New Algorithms

Use the same function signature style as `pushswap.TurkAlgorithm` and duplicate benchmark patterns in `internal/benchmarks/algorithm_bench_test.go`.
Then add corresponding entries under `ALGO_BENCHMARKS` in `Makefile` for perf iteration.
