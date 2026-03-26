package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type BenchmarkTool struct {
	ShaA string
	ShaB string

	TestFile string

	Runs         int
	WarmupRuns   int
	SleepBetween time.Duration

	WorkDir string
}

type RunResult struct {
	Duration time.Duration
	ExitCode int
	Error    error
}

type FileStats struct {
	File     string
	Runs     int
	Failures int

	MeanMs   float64
	StdDevMs float64
}

type BenchmarkReport struct {
	ShaA string
	ShaB string

	StatsA []FileStats
	StatsB []FileStats
}

func NewBenchmarkTool(shaA, shaB, testFile string, runs, warmup int, sleep time.Duration) *BenchmarkTool {
	return &BenchmarkTool{
		ShaA:         shaA,
		ShaB:         shaB,
		TestFile:     testFile,
		Runs:         runs,
		WarmupRuns:   warmup,
		SleepBetween: sleep,
		WorkDir:      os.TempDir(),
	}
}

func (b *BenchmarkTool) addWorktree(sha, label string) (string, error) {
	path := filepath.Join(b.WorkDir, "bench-worktree-"+label)
	b.removeWorktree(path)

	fmt.Printf("  [worktree] Checking out %s → %s\n", sha, path)
	if err := runCmd(".", false, "git", "worktree", "add", "--detach", path, sha); err != nil {
		return "", fmt.Errorf("worktree add failed for %s: %w", sha, err)
	}
	return path, nil
}

func (b *BenchmarkTool) removeWorktree(path string) {
	_ = runCmd(".", true, "git", "worktree", "remove", "--force", path)
}

func (b *BenchmarkTool) buildCompiler(worktreePath, outputBinary string) error {
	fmt.Printf("  [build] %s\n        → %s\n", worktreePath, outputBinary)
	return runCmd(worktreePath, false, "go", "build", "-o", outputBinary, "./cli/cmd")
}

func (b *BenchmarkTool) runOnce(binary, sourceFile string) RunResult {
	cmd := exec.Command(binary, "run", sourceFile)
	cmd.Stdout = nil
	cmd.Stderr = nil

	start := time.Now()
	err := cmd.Run()
	elapsed := time.Since(start)

	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	return RunResult{
		Duration: elapsed,
		ExitCode: exitCode,
		Error:    err,
	}
}

func (b *BenchmarkTool) benchmarkBinary(binary, sourceFile string) []FileStats {
	var durations []float64
	failures := 0

	fmt.Printf("    Benchmarking %s on %s\n", binary, sourceFile)
	if b.WarmupRuns > 0 {
		fmt.Printf("      Warmup runs: %d\n", b.WarmupRuns)
		for i := 0; i < b.WarmupRuns; i++ {
			result := b.runOnce(binary, sourceFile)
			if result.Error != nil {
				fmt.Printf("        Warmup run %d failed: %v\n", i+1, result.Error)
			}
			time.Sleep(b.SleepBetween)
		}
		fmt.Println()
	}

	for i := 0; i < b.Runs; i++ {
		result := b.runOnce(binary, sourceFile)
		if result.Error != nil {
			fmt.Printf("        Run %d failed: %v\n", i+1, result.Error)
			failures++
		} else {
			durations = append(durations, float64(result.Duration.Milliseconds()))
			fmt.Printf("        Run %d: %.2f ms\n", i+1, durations[len(durations)-1])
		}
		if b.SleepBetween > 0 && i < b.Runs-1 {
			time.Sleep(b.SleepBetween)
		}
	}

	stats := FileStats{
		File:     sourceFile,
		Runs:     b.Runs,
		Failures: failures,
	}

	if len(durations) > 0 {
		stats.MeanMs = mean(durations)
		stats.StdDevMs = stddev(durations, stats.MeanMs)
	}

	return []FileStats{stats}
}

func (b *BenchmarkTool) Run() (*BenchmarkReport, error) {
	pathA, err := b.addWorktree(b.ShaA, "A")
	if err != nil {
		return nil, err
	}
	pathB, err := b.addWorktree(b.ShaB, "B")
	if err != nil {
		b.removeWorktree(pathA)
		return nil, err
	}
	defer func() {
		b.removeWorktree(pathA)
		b.removeWorktree(pathB)
	}()

	binaryA := filepath.Join(b.WorkDir, "bench-worktree-A.bin")
	binaryB := filepath.Join(b.WorkDir, "bench-worktree-B.bin")

	if err := b.buildCompiler(pathA, binaryA); err != nil {
		return nil, fmt.Errorf("failed to build compiler for %s: %w", b.ShaA, err)
	}
	if err := b.buildCompiler(pathB, binaryB); err != nil {
		return nil, fmt.Errorf("failed to build compiler for %s: %w", b.ShaB, err)
	}

	statsA := b.benchmarkBinary(binaryA, b.TestFile)
	statsB := b.benchmarkBinary(binaryB, b.TestFile)

	report := &BenchmarkReport{
		ShaA:   b.ShaA,
		ShaB:   b.ShaB,
		StatsA: statsA,
		StatsB: statsB,
	}
	return report, nil
}

func mean(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func stddev(vals []float64, m float64) float64 {
	if len(vals) < 2 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		d := v - m
		sum += d * d
	}
	return math.Sqrt(sum / float64(len(vals)-1))
}

func runCmd(dir string, silent bool, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Dir = dir
	if !silent {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}
	return command.Run()
}

func resolveRef(ref string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--verify", ref)
	cmd.Dir = "."
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("cannot resolve git ref %q - not a valid SHA, tag, branch, or ref", ref)
	}
	return strings.TrimSpace(string(output)), nil
}

func main() {
	shaA, err := resolveRef("HEAD")
	if err != nil {
		fmt.Printf("Error resolving ref: %v\n", err)
		return
	}
	shaB, err := resolveRef("HEAD~1")
	if err != nil {
		fmt.Printf("Error resolving ref: %v\n", err)
		return
	}

	b := NewBenchmarkTool(shaA, shaB, "test.bal", 10, 2, time.Second)

	result, err := b.Run()
	if err != nil {
		fmt.Printf("Benchmark failed: %v\n", err)
		return
	}
	printReport(result)
}
