package main

import (
	"fmt"
	"strings"
)

const (
	colFile  = 18
	colStats = 36
	colDelta = 10
)

func printReport(r *BenchmarkReport) {
	totalWidth := 1 + colFile + 1 + colStats + 1 + colStats + 1 + colDelta + 1

	divider := strings.Repeat("─", totalWidth)
	heavyLine := strings.Repeat("═", totalWidth)

	fmt.Println()
	fmt.Println("╔" + heavyLine + "╗")
	title := fmt.Sprintf("  COMPILER BENCHMARK  │  A = %-20s  │  B = %-20s  ", short(r.ShaA, 20), short(r.ShaB, 20))
	fmt.Printf("║ %-*s ║\n", totalWidth-2, title)
	fmt.Println("╠" + heavyLine + "╣")

	hFile := pad("FILE", colFile)
	hA := center("──── Commit A ────── mean · stddev · min · max · mem", colStats)
	hB := center("──── Commit B ────── mean · stddev · min · max · mem", colStats)
	hDelta := center("Δ mean", colDelta)
	fmt.Printf("║ %s │ %s │ %s │ %s ║\n", hFile, hA, hB, hDelta)
	fmt.Println("╠" + heavyLine + "╣")

	indexB := make(map[string]FileStats)
	for _, s := range r.StatsB {
		indexB[s.File] = s
	}

	for i, a := range r.StatsA {
		b, ok := indexB[a.File]

		deltaStr := "  N/A    "
		indicator := " "
		if ok && a.MeanMs > 0 && a.Failures < a.Runs {
			delta := b.MeanMs - a.MeanMs
			pct := (delta / a.MeanMs) * 100
			if pct > 1 {
				indicator = "▲" // B is slower (regression)
				deltaStr = fmt.Sprintf("%s +%.1f%%", indicator, pct)
			} else if pct < -1 {
				indicator = "▼" // B is faster (improvement)
				deltaStr = fmt.Sprintf("%s %.1f%%", indicator, pct)
			} else {
				indicator = "≈"
				deltaStr = fmt.Sprintf("%s  0.0%%", indicator)
			}
			_ = indicator
		}

		cellA := fmtStats(a)
		cellB := "  N/A"
		if ok {
			cellB = fmtStats(b)
		}

		fmt.Printf("║ %-*s │ %-*s │ %-*s │ %-*s ║\n",
			colFile, short(a.File, colFile),
			colStats, cellA,
			colStats, cellB,
			colDelta, deltaStr,
		)

		if a.Failures > 0 || (ok && b.Failures > 0) {
			warnA := fmt.Sprintf("  ⚠ %d/%d runs failed", a.Failures, a.Runs)
			warnB := ""
			if ok && b.Failures > 0 {
				warnB = fmt.Sprintf("  ⚠ %d/%d runs failed", b.Failures, b.Runs)
			}
			fmt.Printf("║ %-*s │ %-*s │ %-*s │ %-*s ║\n",
				colFile, "",
				colStats, warnA,
				colStats, warnB,
				colDelta, "",
			)
		}

		if i < len(r.StatsA)-1 {
			fmt.Println("║" + " " + divider[:colFile] + "─┼─" + divider[:colStats] + "─┼─" + divider[:colStats] + "─┼─" + divider[:colDelta] + " ║")
		}
	}

	fmt.Println("╚" + heavyLine + "╝")
	fmt.Println()
	fmt.Println("  ▲ = B slower (regression)   ▼ = B faster (improvement)   ≈ = no significant change")
	fmt.Println("  Δ mean = (B − A) / A × 100     mem = mean RSS (Linux only)")
	fmt.Println()
}

func fmtStats(s FileStats) string {
	if s.Failures == s.Runs {
		return "  all runs failed"
	}
	return fmt.Sprintf("%7.1fms ±%-6.2f",
		s.MeanMs, s.StdDevMs)
}

func pad(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s + strings.Repeat(" ", n-len(s))
}

func center(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	total := n - len(s)
	left := total / 2
	right := total - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func short(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
