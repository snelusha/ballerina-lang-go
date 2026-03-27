package main

import "math"

// Measurement holds the timing of a single compiler invocation.
type Measurement struct {
	DurationMs float64
	Failed     bool
}

// FileStats holds the aggregated benchmark statistics for one input file.
type FileStats struct {
	// Label is the display name (file base name or package label).
	Label string

	// Runs is the total number of timed runs attempted.
	Runs int

	// Failures is the count of runs that exited non-zero or errored.
	Failures int

	Mean   float64
	StdDev float64
	Min    float64
	Max    float64
}

// BenchmarkResult is the top-level output of a benchmark run.
// It is intentionally kept as a plain data structure so it can be
// consumed by both the terminal reporter and a future HTML renderer.
type BenchmarkResult struct {
	BaseRef string
	HeadRef string
	Base    []FileStats
	Head    []FileStats
}

func computeStats(label string, runs int, measurements []Measurement) FileStats {
	failures := 0
	var durations []float64

	for _, m := range measurements {
		if m.Failed {
			failures++
		} else {
			durations = append(durations, m.DurationMs)
		}
	}

	s := FileStats{
		Label:    label,
		Runs:     runs,
		Failures: failures,
	}

	if len(durations) > 0 {
		s.Mean = meanOf(durations)
		s.StdDev = stddevOf(durations, s.Mean)
		s.Min = minOf(durations)
		s.Max = maxOf(durations)
	}

	return s
}

func meanOf(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func stddevOf(vals []float64, m float64) float64 {
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

func minOf(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	result := vals[0]
	for _, v := range vals[1:] {
		if v < result {
			result = v
		}
	}
	return result
}

func maxOf(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	result := vals[0]
	for _, v := range vals[1:] {
		if v > result {
			result = v
		}
	}
	return result
}
