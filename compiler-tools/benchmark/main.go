package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type InputMode byte

const (
	SingleFile InputMode = iota
	Package
	Directory
)

type Target struct {
	Mode  InputMode
	Paths []string
	Label string
}

func resolveTarget(input string) (*Target, error) {
	info, err := os.Stat(input)
	if err != nil {
		return nil, fmt.Errorf("input path %q not found: %w", input, err)
	}

	if !info.IsDir() {
		if filepath.Ext(input) != ".bal" {
			return nil, fmt.Errorf("input file %q must have .bal extension", input)
		}
		return &Target{
			Mode:  SingleFile,
			Paths: []string{input},
			Label: filepath.Base(input),
		}, nil
	}

	tomlPath := filepath.Join(input, "Ballerina.toml")
	if _, err := os.Stat(tomlPath); err == nil {
		return &Target{
			Mode:  Package,
			Paths: []string{input},
			Label: fmt.Sprintf("%s (package)", filepath.Base(input)),
		}, nil
	}

	var files []string
	err = filepath.Walk(input, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".bal" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error scanning directory %q: %w", input, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no .bal files found in directory %q", input)
	}

	return &Target{
		Mode:  Directory,
		Paths: files,
		Label: fmt.Sprintf("%s /", filepath.Base(input)),
	}, nil
}

type Benchmark struct {
	RefA        string
	RefB        string
	InputTarget string
	Runs        int
	WarmupRuns  int
	Sleep       time.Duration
	WorkDir     string
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
	Mean     float64
	StdDev   float64
	Min      float64
	Max      float64
}

type BenchmarkResult struct {
	RefA   string
	RefB   string
	StatsA []FileStats
	StatsB []FileStats
}

func (b *Benchmark) addWorktree(ref, label string) (string, error) {
	path := filepath.Join(b.WorkDir, fmt.Sprintf("%s-%s", label, ref))
	b.removeWorktree(path)

	fmt.Printf("  [worktree] Checking out %s → %s\n", ref, path)
	if err := runCmd(".", "git", "worktree", "add", "--detach", path, ref); err != nil {
		return "", fmt.Errorf("worktree add failed for %s: %w", ref, err)
	}
	return path, nil
}

func (b *Benchmark) removeWorktree(path string) {
	_ = runCmdSilent(".", "git", "worktree", "remove", "--force", path)
}

func (b *Benchmark) buildCompiler(path, output string) error {
	fmt.Printf("  [build] %s\n        → %s\n", path, output)
	return runCmd(path, "go", "build", "-o", output, "./cli/cmd")
}

func (b *Benchmark) runOnce(binary, path string) RunResult {
	fmt.Printf("    [run] %s\n", path)
	cmd := exec.Command(binary, "run", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	if err != nil {
		fmt.Printf("    [error] %v (exit code: %d)\n", err, exitCode)
	}

	return RunResult{
		Duration: duration,
		ExitCode: exitCode,
		Error:    err,
	}
}

func (b *Benchmark) benchmarkBinary(binary string, target Target) []FileStats {
	var stats []FileStats

	for _, f := range target.Paths {
		var durations []float64
		failures := 0

		fmt.Printf("  [benchmark] %s\n", filepath.Base(f))

		if b.WarmupRuns > 0 {
			fmt.Printf("    [warmup] %d runs\n", b.WarmupRuns)
			for i := range b.WarmupRuns {
				b.runOnce(binary, f)
				if b.Sleep > 0 && i < b.WarmupRuns-1 {
					time.Sleep(b.Sleep)
				}
			}
		}

		for i := range b.Runs {
			result := b.runOnce(binary, f)
			if result.Error != nil || result.ExitCode != 0 {
				failures = failures + 1
			} else {
				durations = append(durations, float64(result.Duration.Milliseconds()))
			}

			if b.Sleep > 0 && i < b.WarmupRuns-1 {
				time.Sleep(b.Sleep)
			}
		}

		label := filepath.Base(f)
		if target.Mode == Package {
			label = target.Label
		}

		sts := FileStats{
			File:     label,
			Runs:     b.Runs,
			Failures: failures,
		}

		if len(durations) > 0 {
			sts.Mean = mean(durations)
			sts.StdDev = stddev(durations, sts.Mean)
			sts.Min = min(durations)
			sts.Max = max(durations)
		}

		stats = append(stats, sts)
	}

	return stats
}

func (b *Benchmark) Run() (*BenchmarkResult, error) {
	pathA, err := b.addWorktree(b.RefA, "a")
	if err != nil {
		return nil, err
	}
	pathB, err := b.addWorktree(b.RefB, "b")
	if err != nil {
		b.removeWorktree(pathA)
		return nil, err
	}
	defer func() {
		b.removeWorktree(pathA)
		b.removeWorktree(pathB)
	}()

	binaryA := filepath.Join(b.WorkDir, "bal-a")
	binaryB := filepath.Join(b.WorkDir, "bal-b")

	if err := b.buildCompiler(pathA, binaryA); err != nil {
		return nil, fmt.Errorf("failed to build compiler for %s: %w", b.RefA, err)
	}
	if err := b.buildCompiler(pathB, binaryB); err != nil {
		return nil, fmt.Errorf("failed to build compiler for %s: %w", b.RefB, err)
	}

	target, err := resolveTarget(b.InputTarget)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve input target: %w", err)
	}

	fmt.Printf("Benchmarking %s vs %s on %d runs\n", b.RefA, b.RefB, b.Runs)
	statsA := b.benchmarkBinary(binaryA, *target)
	statsB := b.benchmarkBinary(binaryB, *target)

	return &BenchmarkResult{
		RefA:   b.RefA,
		RefB:   b.RefB,
		StatsA: statsA,
		StatsB: statsB,
	}, nil
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

func min(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	min := vals[0]
	for _, v := range vals[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func max(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	max := vals[0]
	for _, v := range vals[1:] {
		if v > max {
			max = v
		}
	}
	return max
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

func runCmd(dir, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func runCmdSilent(dir, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Dir = dir
	return c.Run()
}

func run() error {
	cfg, err := parseArgs()
	if err != nil {
		return err
	}

	if err := cfg.validate(); err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}

	benchmark := &Benchmark{
		RefA:        cfg.refA,
		RefB:        cfg.refB,
		InputTarget: cfg.inputTarget,
		Runs:        cfg.runs,
		WarmupRuns:  cfg.warmup,
		Sleep:       time.Duration(cfg.sleep) * time.Millisecond,
		WorkDir:     os.TempDir(),
	}

	result, err := benchmark.Run()
	if err != nil {
		return fmt.Errorf("benchmark failed: %w", err)
	}

	printReport(result)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
