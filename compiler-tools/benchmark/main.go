package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := parseFlags()
	if err != nil {
		return err
	}
	if err := cfg.validate(); err != nil {
		return fmt.Errorf("invalid arguments: %w", err)
	}

	b := &Benchmark{
		BaseRef:    cfg.baseRef,
		HeadRef:    cfg.headRef,
		Target:     cfg.inputTarget,
		Runs:       cfg.runs,
		WarmupRuns: cfg.warmup,
		SleepMs:    cfg.sleepMs,
		WorkDir:    os.TempDir(),
	}

	result, err := b.Run()
	if err != nil {
		return fmt.Errorf("benchmark failed: %w", err)
	}

	printReport(result)
	return nil
}
