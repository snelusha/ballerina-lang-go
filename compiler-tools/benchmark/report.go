package main

import (
	"fmt"
	"math"
)

func printReport(r *benchmarkResult) {
	headByLabel := make(map[string]fileStats, len(r.head))
	for _, s := range r.head {
		headByLabel[s.label] = s
	}

	benchmarkIndex := 1

	for _, base := range r.base {
		head, hasHead := headByLabel[base.label]

		fmt.Println()
		fmt.Printf("  Benchmark %d: %s (%s)\n",
			benchmarkIndex, cyan(r.baseRef), base.label)
		benchmarkIndex++

		printTimingBlock(base)

		fmt.Println()
		fmt.Printf("  Benchmark %d: %s (%s)\n",
			benchmarkIndex, cyan(r.headRef), base.label)
		benchmarkIndex++

		if hasHead {
			printTimingBlock(head)
		} else {
			fmt.Printf("  %s\n", dim("no data"))
		}
	}

	fmt.Println()
	fmt.Printf("  %s\n", bold("Summary"))

	for _, base := range r.base {
		head, hasHead := headByLabel[base.label]
		if !hasHead || base.mean <= 0 || base.failures == base.runs || head.failures == head.runs {
			continue
		}

		ratio, ratioErr := speedupRatio(base, head)
		fasterRef, slowerRef := r.headRef, r.baseRef
		if ratio < 1.0 {
			ratio = 1.0 / ratio
			ratioErr = ratioErr / (ratio * ratio)
			fasterRef, slowerRef = r.baseRef, r.headRef
		}

		fmt.Printf("  %s (%s) ran\n", cyan(fasterRef), base.label)
		fmt.Printf("  %s times faster than %s (%s)\n",
			green(fmt.Sprintf("%.2f ± %.2f", ratio, ratioErr)),
			cyan(slowerRef),
			base.label,
		)
	}

	fmt.Println()
}

func printTimingBlock(s fileStats) {
	if s.failures == s.runs {
		fmt.Printf("  %s\n", red(fmt.Sprintf("all %d runs failed", s.runs)))
		return
	}

	suffix := ""
	if s.failures > 0 {
		suffix = "  " + red(fmt.Sprintf("%d failed", s.failures))
	}

	fmt.Printf("  %s %s ± %s%s\n",
		dim("Time (mean ± σ):  "),
		bold(fmt.Sprintf("%7.1f ms", s.mean)),
		fmt.Sprintf("%7.1f ms", s.stddev),
		suffix,
	)
	fmt.Printf("  %s %s … %s    %s\n",
		dim("Range (min … max):"),
		cyan(fmt.Sprintf("%7.1f ms", s.min)),
		cyan(fmt.Sprintf("%.1f ms", s.max)),
		dim(fmt.Sprintf("%d runs", s.runs-s.failures)),
	)
}

func speedupRatio(base, head fileStats) (ratio, uncertainty float64) {
	ratio = base.mean / head.mean
	relBase := base.stddev / base.mean
	relHead := head.stddev / head.mean
	uncertainty = ratio * math.Sqrt(relBase*relBase+relHead*relHead)
	return ratio, uncertainty
}

const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiDim   = "\033[2m"
	ansiCyan  = "\033[36m"
	ansiGreen = "\033[32m"
	ansiRed   = "\033[31m"
)

func bold(s string) string  { return ansiBold + s + ansiReset }
func dim(s string) string   { return ansiDim + s + ansiReset }
func cyan(s string) string  { return ansiCyan + s + ansiReset }
func green(s string) string { return ansiGreen + ansiBold + s + ansiReset }
func red(s string) string   { return ansiRed + ansiBold + s + ansiReset }
