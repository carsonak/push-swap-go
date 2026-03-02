.PHONY: bench-algo bench-stack perf-setup perf-restore perf-bench-stack perf-bench-algo

# Perf tool path
PERF := /usr/lib/linux-tools-6.8.0-100/perf
PERF_PARANOID := /proc/sys/kernel/perf_event_paranoid
PERF_PARANOID_BACKUP := /tmp/perf_paranoid_backup
BENCHMARKS_DIR := ./internal/benchmarks

# Algorithm benchmark target - runs push-swap algorithm benchmarks
bench-algo:
	go test -run=^$$ -bench=BenchmarkTurkAlgorithm_ -benchmem -benchtime=3s $(BENCHMARKS_DIR)

# Stack benchmark target - runs stack comparison benchmarks
bench-stack:
	go test -run=^$$ -bench='BenchmarkPrealloc_|BenchmarkRatios|BenchmarkStackFrames' -benchmem $(BENCHMARKS_DIR)

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
STACK_BENCHMARKS := \
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

# Run stack benchmarks with perf stats (each subtest separately)
perf-bench-stack:
	@echo "Setting up perf environment..."
	@sudo sh -c 'cat $(PERF_PARANOID) > $(PERF_PARANOID_BACKUP)'
	@sudo sh -c 'echo -1 > $(PERF_PARANOID)'
	@for bench in $(STACK_BENCHMARKS); do \
		output_file=$$(echo "$$bench" | tr '/' '-').txt; \
		echo "Running $$bench with perf stats (output to $$output_file)..."; \
		$(PERF) stat -e instructions,cache-references,cache-misses --output=$$output_file go test -run=^$$ -bench=$$bench ./... || true; \
	done
	@echo "Restoring perf environment..."
	@sudo sh -c 'cat $(PERF_PARANOID_BACKUP) > $(PERF_PARANOID)'

# Algorithm benchmarks to profile
ALGO_BENCHMARKS := \
	BenchmarkTurkAlgorithm_MassiveScale/100 \
	BenchmarkTurkAlgorithm_MassiveScale/1000 \
	BenchmarkTurkAlgorithm_MassiveScale/3000 \
	BenchmarkTurkAlgorithm_MassiveScale/5000 \
	BenchmarkTurkAlgorithm_MassiveScale/10000 \
	BenchmarkTurkAlgorithm_MassiveScale/50000 \
	BenchmarkTurkAlgorithm_StandardFloats/500 \
	BenchmarkTurkAlgorithm_StandardFloats/1000 \
	BenchmarkTurkAlgorithm_StandardFloats/1500 \
	BenchmarkTurkAlgorithm_StandardFloats/1900 \
	BenchmarkTurkAlgorithm_HighDensityDuplicates/500 \
	BenchmarkTurkAlgorithm_HighDensityDuplicates/1000 \
	BenchmarkTurkAlgorithm_HighDensityDuplicates/1500 \
	BenchmarkTurkAlgorithm_HighDensityDuplicates/1900 \
	BenchmarkTurkAlgorithm_NearlySorted_TopHeavy/500 \
	BenchmarkTurkAlgorithm_NearlySorted_TopHeavy/750 \
	BenchmarkTurkAlgorithm_NearlySorted_TopHeavy/1000 \
	BenchmarkTurkAlgorithm_NearlySorted_MiddleHeavy/500 \
	BenchmarkTurkAlgorithm_NearlySorted_MiddleHeavy/750 \
	BenchmarkTurkAlgorithm_NearlySorted_MiddleHeavy/1000 \
	BenchmarkTurkAlgorithm_NearlySorted_BottomHeavy/500 \
	BenchmarkTurkAlgorithm_NearlySorted_BottomHeavy/750 \
	BenchmarkTurkAlgorithm_NearlySorted_BottomHeavy/1000 \
	BenchmarkTurkAlgorithm_TinyFloats/500 \
	BenchmarkTurkAlgorithm_TinyFloats/1000 \
	BenchmarkTurkAlgorithm_TinyFloats/1500 \
	BenchmarkTurkAlgorithm_TinyFloats/1900 \
	BenchmarkTurkAlgorithm_MassiveFloats/500 \
	BenchmarkTurkAlgorithm_MassiveFloats/1000 \
	BenchmarkTurkAlgorithm_MassiveFloats/1500 \
	BenchmarkTurkAlgorithm_MassiveFloats/1900

# Run algorithm benchmarks with perf stats
perf-bench-algo:
	@echo "Setting up perf environment..."
	@sudo sh -c 'cat $(PERF_PARANOID) > $(PERF_PARANOID_BACKUP)'
	@sudo sh -c 'echo -1 > $(PERF_PARANOID)'
	@for bench in $(ALGO_BENCHMARKS); do \
		output_file=$$(echo "$$bench" | tr '/' '-').txt; \
		echo "Running $$bench with perf stats (output to $$output_file)..."; \
		$(PERF) stat -e instructions,cache-references,cache-misses,LLC-loads,LLC-load-misses --output=$$output_file go test -run=^$$ -bench=$$bench -benchtime=3s $(BENCHMARKS_DIR) || true; \
	done
	@echo "Restoring perf environment..."
	@sudo sh -c 'cat $(PERF_PARANOID_BACKUP) > $(PERF_PARANOID)'
