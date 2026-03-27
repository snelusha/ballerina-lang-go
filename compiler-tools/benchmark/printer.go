package main

import (
	"fmt"
	"strings"
)

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
