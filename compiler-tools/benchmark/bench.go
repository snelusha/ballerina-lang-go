package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type inputMode byte

const (
	modeSingleFile inputMode = iota
	modePackage
	modeDirectory
)

type (
	benchmark struct {
		baseRef    string
		headRef    string
		target     string
		runs       int
		warmupRuns int
		sleepMs    int
		workDir    string
	}
	measurement struct {
		durationMs float64
		failed     bool
	}
	fileStats struct {
		label    string
		runs     int
		failures int
		mean     float64
		stddev   float64
		min      float64
		max      float64
	}
	target struct {
		mode  inputMode
		paths []string
		label string
	}
	benchmarkResult struct {
		baseRef string
		headRef string
		base    []fileStats
		head    []fileStats
	}
)

func (b *benchmark) run() (*benchmarkResult, error) {
	worktreeBase, err := b.checkoutWorktree(b.baseRef, "base")
	if err != nil {
		return nil, err
	}
	defer b.removeWorktree(worktreeBase)

	worktreeHead, err := b.checkoutWorktree(b.headRef, "head")
	if err != nil {
		return nil, err
	}
	defer b.removeWorktree(worktreeHead)

	binaryBase := filepath.Join(worktreeBase, "bal-base")
	binaryHead := filepath.Join(worktreeHead, "bal-head")

	if err := b.buildCompiler(worktreeBase, binaryBase, b.baseRef); err != nil {
		return nil, err
	}
	if err := b.buildCompiler(worktreeHead, binaryHead, b.headRef); err != nil {
		return nil, err
	}

	target, err := resolveTarget(b.target)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve target: %w", err)
	}

	fmt.Printf("Benchmarking %s vs %s - %d runs", b.baseRef, b.headRef, b.runs)
	if b.warmupRuns > 0 {
		fmt.Printf(" (%d warmup)", b.warmupRuns)
	}
	fmt.Println()

	return &benchmarkResult{
		baseRef: b.baseRef,
		headRef: b.headRef,
		base:    b.runBenchmark(binaryBase, *target),
		head:    b.runBenchmark(binaryHead, *target),
	}, nil
}

func (b *benchmark) checkoutWorktree(ref, slot string) (string, error) {
	path := filepath.Join(b.workDir, fmt.Sprintf("bench-%s", slot))
	_ = runCmd(".", "git", "worktree", "remove", "--force", path)

	fmt.Printf("Checking out %q...\n", ref)
	if err := runCmd(".", "git", "worktree", "add", "--detach", path, ref); err != nil {
		return "", fmt.Errorf("failed to checkout %q: %w", ref, err)
	}
	return path, nil
}

func (b *benchmark) removeWorktree(path string) {
	_ = runCmd(".", "git", "worktree", "remove", "--force", path)
}

func (b *benchmark) buildCompiler(worktreePath, outputBinary, ref string) error {
	fmt.Printf("Building compiler for %q...\n", ref)
	if err := runCmd(worktreePath, "go", "build", "-o", outputBinary, "./cli/cmd"); err != nil {
		return fmt.Errorf("failed to build compiler for %q: %w", ref, err)
	}
	return nil
}

func (b *benchmark) runBenchmark(binary string, target target) []fileStats {
	sleep := time.Duration(b.sleepMs) * time.Millisecond
	var stats []fileStats

	for _, inputPath := range target.paths {
		label := filepath.Base(inputPath)
		if target.mode == modePackage {
			label = target.label
		}

		fmt.Printf("Running benchmark for %q...\n", label)

		if b.warmupRuns > 0 {
			fmt.Printf("  Warmup runs (%d)...\n", b.warmupRuns)
			for i := range b.warmupRuns {
				b.invokeCompiler(binary, inputPath)
				if sleep > 0 && i < b.warmupRuns-1 {
					time.Sleep(sleep)
				}
			}
		}

		measurements := make([]measurement, b.runs)
		for i := range b.runs {
			measurements = append(measurements, b.invokeCompiler(binary, inputPath))
			if sleep > 0 && i < b.runs-1 {
				time.Sleep(sleep)
			}
		}

		stats = append(stats, computeStats(label, b.runs, measurements))
	}

	return stats
}

func (b *benchmark) invokeCompiler(binary string, inputPath string) measurement {
	cmd := exec.Command(binary, "run", inputPath)
	cmd.Stdout = nil
	cmd.Stderr = nil

	start := time.Now()
	err := cmd.Run()
	elapsed := time.Since(start)

	failed := err != nil || (cmd.ProcessState != nil && cmd.ProcessState.ExitCode() != 0)
	return measurement{
		durationMs: float64(elapsed.Milliseconds()),
		failed:     failed,
	}
}

func resolveTarget(input string) (*target, error) {
	info, err := os.Stat(input)
	if err != nil {
		return nil, fmt.Errorf("input path %q not found: %w", input, err)
	}

	if !info.IsDir() {
		if filepath.Ext(input) != ".bal" {
			return nil, fmt.Errorf("input file %q is not a ballerina source file", input)
		}
		return &target{
			mode:  modeSingleFile,
			paths: []string{input},
			label: filepath.Base(input),
		}, nil
	}

	if _, err := os.Stat(filepath.Join(input, "Ballerina.toml")); err == nil {
		return &target{
			mode:  modePackage,
			paths: []string{input},
			label: fmt.Sprintf("%s (package)", filepath.Base(input)),
		}, nil
	}

	files, err := collectBalFiles(input)
	if err != nil {
		return nil, fmt.Errorf("failed to collect .bal files from directory %q: %w", input, err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no .bal files found in directory %q", input)
	}

	return &target{
		mode:  modeDirectory,
		paths: files,
		label: fmt.Sprintf("%s (directory)", filepath.Base(input)),
	}, nil
}

func collectBalFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".bal" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory %q: %w", dir, err)
	}

	return files, nil
}

func runCmd(dir, cmd string, args ...string) error {
	var buf bytes.Buffer
	c := exec.Command(cmd, args...)
	c.Dir = dir
	c.Stdout = &buf
	c.Stderr = &buf
	if err := c.Run(); err != nil {
		return fmt.Errorf("%w\n%s", err, buf.String())
	}
	return nil
}

func computeStats(label string, runs int, memeasurements []measurement) fileStats {
	var durations []float64
	failures := 0
	for _, m := range memeasurements {
		durations = append(durations, m.durationMs)
		if m.failed {
			failures++
		}
	}

	m := mean(durations)
	return fileStats{
		label:    label,
		runs:     runs,
		failures: failures,
		mean:     m,
		stddev:   stddev(durations, m),
		min:      min(durations),
		max:      max(durations),
	}
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
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

func min(vals []float64) float64 {
	result := vals[0]
	for _, v := range vals[1:] {
		if v < result {
			result = v
		}
	}
	return result
}

func max(vals []float64) float64 {
	result := vals[0]
	for _, v := range vals[1:] {
		if v > result {
			result = v
		}
	}
	return result
}
