package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Benchmark orchestrates the full compare-two-refs workflow:
//  1. Check out each ref into a temporary git worktree.
//  2. Build the Ballerina compiler from each worktree.
//  3. Run each compiled binary against the target input.
//  4. Return aggregated statistics for both refs.
type Benchmark struct {
	BaseRef    string
	HeadRef    string
	Target     string
	Runs       int
	WarmupRuns int
	SleepMs    int
	WorkDir    string
}

// Run executes the full benchmark pipeline and returns the result.
func (b *Benchmark) Run() (*BenchmarkResult, error) {
	worktreeBase, err := b.checkoutWorktree(b.BaseRef, "base")
	if err != nil {
		return nil, err
	}
	defer b.removeWorktree(worktreeBase)

	worktreeHead, err := b.checkoutWorktree(b.HeadRef, "head")
	if err != nil {
		return nil, err
	}
	defer b.removeWorktree(worktreeHead)

	binaryBase := filepath.Join(b.WorkDir, "bal-base")
	binaryHead := filepath.Join(b.WorkDir, "bal-head")

	if err := b.buildCompiler(worktreeBase, binaryBase, b.BaseRef); err != nil {
		return nil, err
	}
	if err := b.buildCompiler(worktreeHead, binaryHead, b.HeadRef); err != nil {
		return nil, err
	}

	target, err := resolveTarget(b.Target)
	if err != nil {
		return nil, fmt.Errorf("resolving input target: %w", err)
	}

	fmt.Printf("\n%s Benchmarking %s vs %s — %d runs",
		cyan("◆"), bold(b.BaseRef), bold(b.HeadRef), b.Runs)
	if b.WarmupRuns > 0 {
		fmt.Printf(" (%d warmup)", b.WarmupRuns)
	}
	fmt.Println()

	statsBase := b.benchmarkBinary(binaryBase, *target)
	statsHead := b.benchmarkBinary(binaryHead, *target)

	return &BenchmarkResult{
		BaseRef: b.BaseRef,
		HeadRef: b.HeadRef,
		Base:    statsBase,
		Head:    statsHead,
	}, nil
}

// checkoutWorktree creates a detached git worktree for the given ref.
func (b *Benchmark) checkoutWorktree(ref, slot string) (string, error) {
	path := filepath.Join(b.WorkDir, fmt.Sprintf("balbench-%s-%s", slot, ref))
	_ = execQuiet(".", "git", "worktree", "remove", "--force", path)

	fmt.Printf("  %s Checking out %s …\n", dim("→"), bold(ref))
	if err := execQuiet(".", "git", "worktree", "add", "--detach", path, ref); err != nil {
		return "", fmt.Errorf("git worktree add failed for %s: %w", ref, err)
	}
	return path, nil
}

// removeWorktree removes a git worktree, silently ignoring errors.
func (b *Benchmark) removeWorktree(path string) {
	_ = execQuiet(".", "git", "worktree", "remove", "--force", path)
}

// buildCompiler compiles the Go compiler binary from the given worktree.
func (b *Benchmark) buildCompiler(worktreePath, outputBinary, ref string) error {
	fmt.Printf("  %s Building compiler for %s …\n", dim("→"), bold(ref))
	if err := execQuiet(worktreePath, "go", "build", "-o", outputBinary, "./cli/cmd"); err != nil {
		return fmt.Errorf("build failed for %s: %w", ref, err)
	}
	return nil
}

// benchmarkBinary runs the given binary against every path in the target
// and returns one FileStats per path.
func (b *Benchmark) benchmarkBinary(binary string, target Target) []FileStats {
	sleep := time.Duration(b.SleepMs) * time.Millisecond
	var results []FileStats

	for _, inputPath := range target.Paths {
		label := filepath.Base(inputPath)
		if target.Mode == modePackage {
			label = target.Label
		}

		fmt.Printf("\n  %s %s\n", cyan("■"), bold(label))

		if b.WarmupRuns > 0 {
			fmt.Printf("  %s warming up (%d runs) …\n", dim("·"), b.WarmupRuns)
			for i := range b.WarmupRuns {
				b.invokeCompiler(binary, inputPath)
				if sleep > 0 && i < b.WarmupRuns-1 {
					time.Sleep(sleep)
				}
			}
		}

		measurements := make([]Measurement, 0, b.Runs)
		for i := range b.Runs {
			m := b.invokeCompiler(binary, inputPath)
			measurements = append(measurements, m)
			if sleep > 0 && i < b.Runs-1 { // correctly guards against trailing sleep
				time.Sleep(sleep)
			}
		}

		results = append(results, computeStats(label, b.Runs, measurements))
	}

	return results
}

// invokeCompiler runs the compiler binary once against the given file and
// returns a Measurement. Compiler output is suppressed.
func (b *Benchmark) invokeCompiler(binary, inputPath string) Measurement {
	cmd := exec.Command(binary, "run", inputPath)
	cmd.Stdout = os.Stderr // suppress; redirect to stderr sink only if debugging
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	elapsed := time.Since(start)

	failed := err != nil || (cmd.ProcessState != nil && cmd.ProcessState.ExitCode() != 0)
	return Measurement{
		DurationMs: float64(elapsed.Milliseconds()),
		Failed:     failed,
	}
}
