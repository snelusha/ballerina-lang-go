package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type config struct {
	baseRef     string
	headRef     string
	inputTarget string
	runs        int
	warmup      int
	sleepMs     int
}

func (c *config) validate() error {
	if c.baseRef == "" {
		return errors.New("base ref cannot be empty")
	}
	if c.headRef == "" {
		return errors.New("head ref cannot be empty")
	}
	if c.baseRef == c.headRef {
		return errors.New("base and head refs must be different")
	}
	if _, err := os.Stat(c.inputTarget); os.IsNotExist(err) {
		return fmt.Errorf("input target %q does not exist", c.inputTarget)
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
	if c.sleepMs < 0 {
		return errors.New("--sleep must be 0 or greater")
	}
	return nil
}

func parseFlags() (*config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	runs := fs.Int("runs", 10, "Number of benchmark runs per command")
	warmup := fs.Int("warmup", 0, "Number of warmup runs (excluded from results)")
	sleepMs := fs.Int("sleep", 0, "Milliseconds to sleep between runs")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <base-ref> <head-ref> <input>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	args := fs.Args()
	if len(args) < 3 {
		fs.Usage()
		return nil, errors.New("missing required positional arguments: <base-ref> <head-ref> <input>")
	}

	return &config{
		baseRef:     args[0],
		headRef:     args[1],
		inputTarget: args[2],
		runs:        *runs,
		warmup:      *warmup,
		sleepMs:     *sleepMs,
	}, nil
}
