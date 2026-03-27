package main

import (
	"fmt"
	"strings"
)

// printReport writes a hyperfine-style comparison to stdout.
// BenchmarkResult is kept as a pure data type so a future HTML renderer
// (or JSON exporter) can consume the same struct without touching this file.
func printReport(r *BenchmarkResult) {
	fmt.Println()
	fmt.Printf("  %s\n\n", bold("Benchmark Results"))

	// Index head stats by label for O(1) lookup.
	headByLabel := make(map[string]FileStats, len(r.Head))
	for _, s := range r.Head {
		headByLabel[s.Label] = s
	}

	for _, base := range r.Base {
		head, hasHead := headByLabel[base.Label]

		// ── File heading ──────────────────────────────────────────────────────
		fmt.Printf("  %s %s\n", cyan("◆"), bold(base.Label))

		// ── Per-ref timing rows ───────────────────────────────────────────────
		fmt.Printf("  %s  %s\n", refLabel("base", r.BaseRef), formatStatLine(base))
		if hasHead {
			fmt.Printf("  %s  %s\n", refLabel("head", r.HeadRef), formatStatLine(head))
		} else {
			fmt.Printf("  %s  %s\n", refLabel("head", r.HeadRef), dim("no data"))
		}

		// ── Delta summary ─────────────────────────────────────────────────────
		const indent = "                        "
		if hasHead && base.Mean > 0 && base.Failures < base.Runs && head.Failures < head.Runs {
			delta := head.Mean - base.Mean
			pct := (delta / base.Mean) * 100

			switch {
			case pct > 1:
				fmt.Printf("  %s%s is %s slower  %s\n",
					indent, bold("head"), red(fmt.Sprintf("+%.2f%%", pct)), dim("▲ regression"),
				)
			case pct < -1:
				fmt.Printf("  %s%s is %s faster  %s\n",
					indent, bold("head"), green(fmt.Sprintf("%.2f%%", pct)), dim("▼ improvement"),
				)
			default:
				fmt.Printf("  %s%s\n", indent, dim("no significant difference"))
			}
		}

		// ── Failure warnings ──────────────────────────────────────────────────
		if base.Failures > 0 {
			fmt.Printf("  %sbase: %s\n", indent,
				red(fmt.Sprintf("%d/%d runs failed", base.Failures, base.Runs)))
		}
		if hasHead && head.Failures > 0 {
			fmt.Printf("  %shead: %s\n", indent,
				red(fmt.Sprintf("%d/%d runs failed", head.Failures, head.Runs)))
		}

		fmt.Println()
	}

	fmt.Printf("  %s\n", dim("mean ± σ   [min … max]"))
	fmt.Printf("  %s\n\n", dim("Δ relative to base  ·  "+green("▼ faster")+"  "+red("▲ slower")))
}

// formatStatLine returns a single formatted timing line for one FileStats.
func formatStatLine(s FileStats) string {
	if s.Failures == s.Runs {
		return red("all runs failed")
	}
	return fmt.Sprintf("%s ± %s   %s",
		bold(fmt.Sprintf("%.1f ms", s.Mean)),
		fmt.Sprintf("%.2f ms", s.StdDev),
		dim(fmt.Sprintf("[%.1f … %.1f ms]", s.Min, s.Max)),
	)
}

// refLabel formats the left-hand column for a ref row, e.g. "base  main  ".
func refLabel(slot, ref string) string {
	label := fmt.Sprintf("%-4s  %-16s", slot, truncate(ref, 16))
	return ansiYellow + label + ansiReset
}

// truncate shortens s to n runes, appending "…" if truncated.
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s + strings.Repeat(" ", n-len(runes))
	}
	return string(runes[:n-1]) + "…"
}

// ── ANSI helpers ─────────────────────────────────────────────────────────────

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
