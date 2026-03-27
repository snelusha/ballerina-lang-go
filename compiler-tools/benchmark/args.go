package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
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
