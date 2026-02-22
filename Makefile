.PHONY: bench perf-setup perf-restore perf-bench perf-bench-prealloc

# Perf tool path
PERF := /usr/lib/linux-tools-6.8.0-100/perf
PERF_PARANOID := /proc/sys/kernel/perf_event_paranoid
PERF_PARANOID_BACKUP := /tmp/perf_paranoid_backup
BENCHMARKS_DIR := ./internal/benchmarks

# Standard benchmark target
bench:
	go test -run=^$$ -bench=. -benchmem $(BENCHMARKS_DIR)

# Setup perf by saving current value and setting to -1
perf-setup:
	@echo "Setting up perf environment..."
	@sudo sh -c 'cat $(PERF_PARANOID) > $(PERF_PARANOID_BACKUP)'
	@sudo sh -c 'echo -1 > $(PERF_PARANOID)'
	@echo "Perf setup complete. Previous value saved to $(PERF_PARANOID_BACKUP)"

# Restore previous perf paranoid value
perf-restore:
	@echo "Restoring perf environment..."
	@sudo sh -c 'cat $(PERF_PARANOID_BACKUP) > $(PERF_PARANOID)'
	@echo "Perf restored to previous value"

# All benchmarks to profile
BENCHMARKS := \
	BenchmarkPrealloc_Slice/No_Pre-allocation \
	BenchmarkPrealloc_Slice/With_Pre-allocation \
	BenchmarkPrealloc_DLL/No_Pre-allocation \
	BenchmarkPrealloc_DLL/With_Pre-allocation \
	BenchmarkRatios/Slice_More_Pushes \
	BenchmarkRatios/DLL_More_Pushes \
	BenchmarkRatios/Slice_More_Pops \
	BenchmarkRatios/DLL_More_Pops \
	BenchmarkRatios/Slice_Equal_PushPop \
	BenchmarkRatios/DLL_Equal_PushPop \
	BenchmarkRatios/Slice_Heavy_Rotations \
	BenchmarkRatios/DLL_Heavy_Rotations \
	BenchmarkStackFrames/Slice_StackFrames \
	BenchmarkStackFrames/DLL_StackFrames

# Run all benchmarks with perf stats (each subtest separately)
perf-bench:
	@echo "Setting up perf environment..."
	@sudo sh -c 'cat $(PERF_PARANOID) > $(PERF_PARANOID_BACKUP)'
	@sudo sh -c 'echo -1 > $(PERF_PARANOID)'
	@for bench in $(BENCHMARKS); do \
		echo "------------------------------------------------\n"; \
		$(PERF) stat -e instructions,cache-references,cache-misses go test -run=^$$ -bench=$$bench -benchmem $(BENCHMARKS_DIR) || true; \
	done
	@echo "Restoring perf environment..."
	@sudo sh -c 'cat $(PERF_PARANOID_BACKUP) > $(PERF_PARANOID)'
