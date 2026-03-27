package main

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type config struct {
	runs        int
	warmup      int
	sleep       int
	refA        string
	refB        string
	inputTarget string
}

func (c *config) validate() error {
	if c.refA == "" {
		return errors.New("refA cannot be empty")
	}
	if c.refB == "" {
		return errors.New("refB cannot be empty")
	}
	if c.refA == c.refB {
		return errors.New("refA and refB must be different")
	}
	if _, err := os.Stat(c.inputTarget); os.IsNotExist(err) {
		return fmt.Errorf("inputTarget %q does not exist", c.inputTarget)
	}
	if c.runs <= 0 {
		return errors.New("--runs must be greater than 0")
	}
	if c.warmup < 0 {
		return errors.New("--warmup must be 0 or greater")
	}
	if c.warmup >= c.runs {
		return errors.New("--warmup must be less than --runs")
	}
	if c.sleep < 0 {
		return errors.New("--sleep must be 0 or greater")
	}
	return nil
}

func parseArgs() (*config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	runs := fs.Int("runs", 10, "Number of benchmark runs")
	warmup := fs.Int("warmup", 0, "Number of warmup runs")
	sleep := fs.Int("sleep", 0, "Sleep between runs in milliseconds")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <ref-a> <ref-b> <input-target>\n", os.Args[0])
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	args := fs.Args()
	if len(args) < 3 {
		fs.Usage()
		return nil, errors.New("missing required arguments: <refA> <refB> <inputTarget>")
	}

	return &config{
		runs:        *runs,
		warmup:      *warmup,
		sleep:       *sleep,
		refA:        args[0],
		refB:        args[1],
		inputTarget: args[2],
	}, nil
}

func printReport(r *BenchmarkResult) {
	fmt.Println()
	fmt.Printf("  %s\n\n", bold("Benchmark Results"))

	indexB := make(map[string]FileStats)
	for _, s := range r.StatsB {
		indexB[s.File] = s
	}

	for _, a := range r.StatsA {
		b, ok := indexB[a.File]

		// ── File label ──
		fmt.Printf("  %s %s\n", cyan("◆"), bold(a.File))

		// ── Commit A row ──
		fmt.Printf("  %s  %s\n",
			lbl(fmt.Sprintf("A %-18s", short(r.RefA, 18))),
			fmtStatLine(a),
		)

		// ── Commit B row ──
		if ok {
			fmt.Printf("  %s  %s\n",
				lbl(fmt.Sprintf("B %-18s", short(r.RefB, 18))),
				fmtStatLine(b),
			)
		} else {
			fmt.Printf("  %s  %s\n", lbl(fmt.Sprintf("B %-18s", short(r.RefB, 18))), dim("no data"))
		}

		// ── Delta summary line ──
		if ok && a.Mean > 0 && a.Failures < a.Runs && b.Failures < b.Runs {
			delta := b.Mean - a.Mean
			pct := (delta / a.Mean) * 100
			switch {
			case pct > 1:
				fmt.Printf("  %s  %s is %s slower   %s\n",
					strings.Repeat(" ", 22),
					bold("B"), red(fmt.Sprintf("+%.2f%%", pct)),
					dim("regression"),
				)
			case pct < -1:
				fmt.Printf("  %s  %s is %s faster   %s\n",
					strings.Repeat(" ", 22),
					bold("B"), green(fmt.Sprintf("%.2f%%", pct)),
					dim("improvement"),
				)
			default:
				fmt.Printf("  %s  %s\n",
					strings.Repeat(" ", 22),
					dim("no significant change"),
				)
			}
		}

		// ── Failure warnings ──
		if a.Failures > 0 {
			fmt.Printf("  %s  A: %d/%d runs failed\n", strings.Repeat(" ", 22), a.Failures, a.Runs)
		}
		if ok && b.Failures > 0 {
			fmt.Printf("  %s  B: %d/%d runs failed\n", strings.Repeat(" ", 22), b.Failures, b.Runs)
		}

		fmt.Println()
	}

	fmt.Printf("  %s\n", dim("mean ± stddev   [min … max]   mem RSS"))
	fmt.Printf("  %s\n\n", dim("Δ relative to A  ·  green = faster  ·  red = slower"))
}

func fmtStatLine(s FileStats) string {
	if s.Failures == s.Runs {
		return red("all runs failed")
	}
	return fmt.Sprintf("%s ± %s   %s",
		bold(fmt.Sprintf("%.1fms", s.Mean)),
		fmt.Sprintf("%.2fms", s.StdDev),
		dim(fmt.Sprintf("[%.1f … %.1f]", s.Min, s.Max)),
	)
}

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiDim    = "\033[2m"
	ansiCyan   = "\033[36m"
	ansiGreen  = "\033[32m"
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
)

func bold(s string) string  { return ansiBold + s + ansiReset }
func dim(s string) string   { return ansiDim + s + ansiReset }
func cyan(s string) string  { return ansiCyan + ansiBold + s + ansiReset }
func green(s string) string { return ansiGreen + ansiBold + s + ansiReset }
func red(s string) string   { return ansiRed + ansiBold + s + ansiReset }
func lbl(s string) string   { return ansiYellow + s + ansiReset }

func short(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

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
