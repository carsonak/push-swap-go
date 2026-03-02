package tests

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// findProjectRoot finds the project root by looking for go.mod
func findProjectRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}

		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatalf("could not find project root (go.mod)")
		}
		wd = parent
	}
}

// buildBinaries compiles the push-swap and checker programs for testing.
func buildBinaries(t *testing.T) (pushSwapPath, checkerPath string) {
	t.Helper()

	tmpDir := t.TempDir()
	pushSwapPath = filepath.Join(tmpDir, "push-swap")
	checkerPath = filepath.Join(tmpDir, "checker")

	projRoot := findProjectRoot(t)

	// Build push-swap
	cmd := exec.Command("go", "build", "-o", pushSwapPath, "./cmd/push-swap")
	cmd.Dir = projRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build push-swap: %v\n%s", err, output)
	}

	// Build checker
	cmd = exec.Command("go", "build", "-o", checkerPath, "./cmd/checker")
	cmd.Dir = projRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build checker: %v\n%s", err, output)
	}

	return pushSwapPath, checkerPath
}

// runPushSwap runs the push-swap program with the given numbers and returns the output.
func runPushSwap(t *testing.T, pushSwapPath string, numbers []string) []string {
	t.Helper()

	input := strings.Join(numbers, " ")
	cmd := exec.Command(pushSwapPath)
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("push-swap failed: %v\nstderr: %s", err, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}
	}
	return lines
}

// runChecker runs the checker program with the given instructions and numbers.
// Returns ("OK", nil) if successful, ("KO", nil) if not sorted, ("", error) if checker fails.
func runChecker(t *testing.T, checkerPath string, instructions []string, numbers []string) (string, error) {
	t.Helper()

	cmd := exec.Command(checkerPath, append([]string{"--"}, numbers...)...)
	instructionsInput := strings.Join(instructions, "\n")
	cmd.Stdin = strings.NewReader(instructionsInput)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := strings.TrimSpace(stdout.String())
	if err != nil {
		// If there's an error, return it with stderr details
		return "", fmt.Errorf("checker failed: %w, stderr: %s", err, stderr.String())
	}

	return output, nil
}

// numSliceToStrings converts a slice of numbers to string representations.
func numSliceToStrings(nums interface{}) []string {
	switch v := nums.(type) {
	case []int:
		strs := make([]string, len(v))
		for i, n := range v {
			strs[i] = fmt.Sprintf("%d", n)
		}
		return strs
	case []float64:
		strs := make([]string, len(v))
		for i, n := range v {
			strs[i] = fmt.Sprintf("%v", n)
		}
		return strs
	default:
		return []string{}
	}
}

// TestPushSwapIntegration tests push-swap and checker together with various inputs.
func TestPushSwapIntegration(t *testing.T) {
	pushSwapPath, checkerPath := buildBinaries(t)

	tests := []struct {
		name    string
		numbers interface{}
		wantOK  bool
	}{
		{"single element", []int{42}, true},
		{"two elements sorted", []int{1, 2}, true},
		{"two elements reverse", []int{2, 1}, true},
		{"three elements sorted", []int{1, 2, 3}, true},
		{"three elements reverse", []int{3, 2, 1}, true},
		{"three elements random", []int{2, 1, 3}, true},
		{"four elements", []int{3, 1, 4, 2}, true},
		{"five elements", []int{5, 2, 4, 1, 3}, true},
		{"five elements sorted", []int{1, 2, 3, 4, 5}, true},
		{"five elements reverse", []int{5, 4, 3, 2, 1}, true},
		{"negative numbers", []int{-5, 3, -1, 0, 2}, true},
		{"floats", []float64{3.5, 1.2, 4.8, 2.1}, true},
		{"mix negative and positive", []int{-10, 5, -3, 0, 8, -7}, true},
		{"large numbers", []int{1000000, -1000000, 500000, -500000}, true},
		{"eight elements", []int{8, 3, 6, 1, 5, 4, 7, 2}, true},
		{"ten elements", []int{10, 5, 8, 2, 9, 1, 7, 3, 6, 4}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			numbers := numSliceToStrings(tt.numbers)
			instructions := runPushSwap(t, pushSwapPath, numbers)
			result, err := runChecker(t, checkerPath, instructions, numbers)

			if err != nil {
				t.Errorf("checker failed: %v", err)
				return
			}

			if tt.wantOK && result != "OK" {
				t.Errorf("expected OK, got %q", result)
			} else if !tt.wantOK && result != "KO" {
				t.Errorf("expected KO, got %q", result)
			}
		})
	}
}

// TestPushSwapEmptyInput tests push-swap with no numbers.
func TestPushSwapEmptyInput(t *testing.T) {
	pushSwapPath, _ := buildBinaries(t)

	cmd := exec.Command(pushSwapPath)
	cmd.Stdin = strings.NewReader("")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("push-swap failed on empty input: %v\nstderr: %s", err, stderr.String())
	}

	output := strings.TrimSpace(stdout.String())
	if output != "" {
		t.Errorf("expected no output for empty input, got: %q", output)
	}
}

// TestCheckerAlreadySorted tests checker on already sorted input.
func TestCheckerAlreadySorted(t *testing.T) {
	_, checkerPath := buildBinaries(t)

	numbers := []string{"1", "2", "3", "4", "5"}
	result, err := runChecker(t, checkerPath, []string{}, numbers)

	if err != nil {
		t.Errorf("checker failed: %v", err)
	}
	if result != "OK" {
		t.Errorf("expected OK for already sorted input, got %q", result)
	}
}

// TestPushSwapOutputFormat tests that push-swap outputs valid operations.
func TestPushSwapOutputFormat(t *testing.T) {
	pushSwapPath, _ := buildBinaries(t)

	// Operations use lowercase names as defined in operations.go
	validOps := map[string]bool{
		"pa": true, "pb": true, "sa": true, "sb": true, "ss": true,
		"ra": true, "rb": true, "rr": true, "rra": true, "rrb": true, "rrr": true,
	}

	numbers := []string{"5", "3", "1", "4", "2"}
	instructions := runPushSwap(t, pushSwapPath, numbers)

	for _, instr := range instructions {
		instr = strings.TrimSpace(instr)
		if instr == "" {
			continue
		}
		if !validOps[instr] {
			t.Errorf("invalid operation: %q", instr)
		}
	}
}

// TestCheckerWithNegativeNumbers tests integration with negative numbers.
func TestCheckerWithNegativeNumbers(t *testing.T) {
	pushSwapPath, checkerPath := buildBinaries(t)

	numbers := []string{"-5", "10", "-2", "0", "7", "-10", "3"}
	instructions := runPushSwap(t, pushSwapPath, numbers)
	result, err := runChecker(t, checkerPath, instructions, numbers)

	if err != nil {
		t.Errorf("checker failed: %v", err)
		return
	}
	if result != "OK" {
		t.Errorf("expected OK with negative numbers, got %q", result)
	}
}

// TestCheckerWithFloats tests integration with floating-point numbers.
func TestCheckerWithFloats(t *testing.T) {
	pushSwapPath, checkerPath := buildBinaries(t)

	numbers := []string{"3.5", "1.2", "4.8", "2.1", "0.5"}
	instructions := runPushSwap(t, pushSwapPath, numbers)
	result, err := runChecker(t, checkerPath, instructions, numbers)

	if err != nil {
		t.Errorf("checker failed: %v", err)
		return
	}
	if result != "OK" {
		t.Errorf("expected OK with floats, got %q", result)
	}
}

// BenchmarkPushSwapSmall benchmarks push-swap with 10 elements.
func BenchmarkPushSwapSmall(b *testing.B) {
	t := &testing.T{}
	pushSwapPath, _ := buildBinaries(t)

	numbers := []string{"10", "9", "8", "7", "6", "5", "4", "3", "2", "1"}
	input := strings.Join(numbers, " ")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command(pushSwapPath)
		cmd.Stdin = strings.NewReader(input)
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		cmd.Run()
	}
}

// ============================================================================
// Options and Flags Tests
// ============================================================================

func TestPushSwapAllowDuplicatesOption(t *testing.T) {
	pushSwapPath, checkerPath := buildBinaries(t)

	t.Run("without flag reports duplicate parse error", func(t *testing.T) {
		cmd := exec.Command(pushSwapPath)
		cmd.Stdin = strings.NewReader("2 1 2")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			t.Fatalf("push-swap returned unexpected non-zero exit: %v", err)
		}

		if strings.TrimSpace(stdout.String()) != "" {
			t.Errorf("expected no instructions on duplicate parse failure, got %q", stdout.String())
		}
		if !strings.Contains(stderr.String(), "ERROR") {
			t.Errorf("expected duplicate parse error log, got stderr: %q", stderr.String())
		}
	})

	t.Run("with flag accepts duplicates", func(t *testing.T) {
		cmd := exec.Command(pushSwapPath, "-allow-duplicates")
		cmd.Stdin = strings.NewReader("2 1 2")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			t.Fatalf("push-swap failed: %v, stderr: %s", err, stderr.String())
		}

		instructions := strings.Split(strings.TrimSpace(stdout.String()), "\n")
		if len(instructions) == 1 && instructions[0] == "" {
			instructions = []string{}
		}

		checkerCmd := exec.Command(checkerPath, "-allow-duplicates", "--", "2", "1", "2")
		checkerCmd.Stdin = strings.NewReader(strings.Join(instructions, "\n"))

		var checkerOut, checkerErr bytes.Buffer
		checkerCmd.Stdout = &checkerOut
		checkerCmd.Stderr = &checkerErr

		if err := checkerCmd.Run(); err != nil {
			t.Fatalf("checker failed: %v, stderr: %s", err, checkerErr.String())
		}

		if strings.TrimSpace(checkerOut.String()) != "OK" {
			t.Errorf("expected checker OK, got %q", checkerOut.String())
		}
	})
}

func TestCheckerAllowDuplicatesOption(t *testing.T) {
	_, checkerPath := buildBinaries(t)

	t.Run("without flag exits with duplicate parse error", func(t *testing.T) {
		cmd := exec.Command(checkerPath, "--", "1", "1", "2")
		cmd.Stdin = strings.NewReader("")

		if err := cmd.Run(); err == nil {
			t.Fatalf("expected checker to fail on duplicates without -allow-duplicates")
		}
	})

	t.Run("with flag accepts duplicates", func(t *testing.T) {
		cmd := exec.Command(checkerPath, "-allow-duplicates", "--", "1", "1", "2")
		cmd.Stdin = strings.NewReader("")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			t.Fatalf("checker failed: %v, stderr: %s", err, stderr.String())
		}

		if strings.TrimSpace(stdout.String()) != "OK" {
			t.Errorf("expected checker OK, got %q", stdout.String())
		}
	})
}

func TestPushSwapFilesAndAllowDuplicates(t *testing.T) {
	pushSwapPath, checkerPath := buildBinaries(t)
	tmp := t.TempDir()

	input1 := filepath.Join(tmp, "input1.txt")
	output1 := filepath.Join(tmp, "output1.txt")
	input2 := filepath.Join(tmp, "input2.txt")
	output2 := filepath.Join(tmp, "output2.txt")

	if err := os.WriteFile(input1, []byte("2 1 2"), 0644); err != nil {
		t.Fatalf("failed to write input1: %v", err)
	}
	if err := os.WriteFile(input2, []byte("4 3 1 2"), 0644); err != nil {
		t.Fatalf("failed to write input2: %v", err)
	}

	cmd := exec.Command(
		pushSwapPath,
		"-allow-duplicates",
		"-files", input1+","+output1,
		"-files", input2+","+output2,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("push-swap -files failed: %v, stderr: %s", err, stderr.String())
	}

	tests := []struct {
		numbers         []string
		outputPath      string
		allowDuplicates bool
	}{
		{numbers: []string{"2", "1", "2"}, outputPath: output1, allowDuplicates: true},
		{numbers: []string{"4", "3", "1", "2"}, outputPath: output2, allowDuplicates: false},
	}

	for _, tc := range tests {
		content, err := os.ReadFile(tc.outputPath)
		if err != nil {
			t.Fatalf("failed to read output file %s: %v", tc.outputPath, err)
		}

		instructions := strings.Split(strings.TrimSpace(string(content)), "\n")
		if len(instructions) == 1 && instructions[0] == "" {
			instructions = []string{}
		}

		args := []string{"--"}
		if tc.allowDuplicates {
			args = []string{"-allow-duplicates", "--"}
		}
		args = append(args, tc.numbers...)

		checkerCmd := exec.Command(checkerPath, args...)
		checkerCmd.Stdin = strings.NewReader(strings.Join(instructions, "\n"))

		var out, errBuf bytes.Buffer
		checkerCmd.Stdout = &out
		checkerCmd.Stderr = &errBuf

		if err := checkerCmd.Run(); err != nil {
			t.Fatalf("checker validation failed for %s: %v, stderr: %s", tc.outputPath, err, errBuf.String())
		}
		if strings.TrimSpace(out.String()) != "OK" {
			t.Errorf("expected checker OK for %s, got %q", tc.outputPath, out.String())
		}
	}
}

func TestCheckerFilesAndAllowDuplicates(t *testing.T) {
	pushSwapPath, checkerPath := buildBinaries(t)
	tmp := t.TempDir()

	numbers1 := []string{"2", "1", "2"}
	numbers2 := []string{"5", "3", "1", "4", "2"}

	generateInstructions := func(numbers []string, allowDup bool, path string) {
		args := []string{}
		if allowDup {
			args = append(args, "-allow-duplicates")
		}
		cmd := exec.Command(pushSwapPath, args...)
		cmd.Stdin = strings.NewReader(strings.Join(numbers, " "))

		var out, errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf

		if err := cmd.Run(); err != nil {
			t.Fatalf("push-swap failed: %v, stderr: %s", err, errBuf.String())
		}
		if err := os.WriteFile(path, out.Bytes(), 0644); err != nil {
			t.Fatalf("failed to write instructions file: %v", err)
		}
	}

	inst1 := filepath.Join(tmp, "inst1.txt")
	inst2 := filepath.Join(tmp, "inst2.txt")
	num1 := filepath.Join(tmp, "nums1.txt")
	num2 := filepath.Join(tmp, "nums2.txt")

	if err := os.WriteFile(num1, []byte(strings.Join(numbers1, " ")), 0644); err != nil {
		t.Fatalf("failed to write nums1: %v", err)
	}
	if err := os.WriteFile(num2, []byte(strings.Join(numbers2, " ")), 0644); err != nil {
		t.Fatalf("failed to write nums2: %v", err)
	}

	generateInstructions(numbers1, true, inst1)
	generateInstructions(numbers2, false, inst2)

	checkerCmd := exec.Command(
		checkerPath,
		"-allow-duplicates",
		"-files", inst1+","+num1,
		"-files", inst2+","+num2,
	)

	var stdout, stderr bytes.Buffer
	checkerCmd.Stdout = &stdout
	checkerCmd.Stderr = &stderr

	if err := checkerCmd.Run(); err != nil {
		t.Fatalf("checker -files failed: %v, stderr: %s", err, stderr.String())
	}

	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 status lines, got %d: %q", len(lines), stdout.String())
	}
	if lines[0] != "OK" || lines[1] != "OK" {
		t.Errorf("expected both statuses OK, got %q", stdout.String())
	}
}

// BenchmarkPushSwapMedium benchmarks push-swap with 100 elements.
func BenchmarkPushSwapMedium(b *testing.B) {
	t := &testing.T{}
	pushSwapPath, _ := buildBinaries(t)

	numbers := make([]string, 100)
	for i := 0; i < 100; i++ {
		numbers[i] = fmt.Sprintf("%d", 100-i)
	}
	input := strings.Join(numbers, " ")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command(pushSwapPath)
		cmd.Stdin = strings.NewReader(input)
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
		cmd.Run()
	}
}
