package main

import (
	"fmt"
	"html/template"
	"math"
	"os"
	"strings"
	"time"
)

const reportHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Benchmark Report</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:ital,wght@0,100..800;1,100..800&display=swap" rel="stylesheet">
<style>
*, *::before, *::after {
	box-sizing: border-box;
	margin: 0;
	padding: 0;
}

:root {
	--black: #000000;
	--white: #ffffff;
	--border: #e7e3e4;
	--dim: #999999;
}

html, body {
	background: var(--white);
	color: var(--black);
	font-family: 'JetBrains Mono', monospace;
	font-size: 13px;
	line-height: 1.6;
	-webkit-font-smoothing: antialiased;
}

body {
	padding: 48px 40px;
	max-width: 960px;
	margin: 0 auto;
}

header {
	border-bottom: 1px solid var(--border);
	padding-bottom: 24px;
	margin-bottom: 40px;
}

.header-eyebrow {
	font-size: 10px;
	font-weight: 500;
	color: var(--dim);
	letter-spacing: 0.1em;
	text-transform: uppercase;
	margin-bottom: 8px;
}

.header-refs {
	font-size: 20px;
	font-weight: 700;
	letter-spacing: -0.02em;
}

.header-refs .sep {
	font-weight: 300;
	color: var(--dim);
	margin: 0 10px;
}

.header-meta {
	margin-top: 6px;
	font-size: 11px;
	color: var(--dim);
	display: flex;
	gap: 20px;
}

.section-label {
	font-size: 10px;
	font-weight: 600;
	letter-spacing: 0.12em;
	text-transform: uppercase;
	color: var(--dim);
	margin-bottom: 16px;
}

.benchmarks {
	margin-bottom: 48px;
}

.benchmark-block {
	border: 1px solid var(--border);
}

table {
	width: 100%;
	border-collapse: collapse;
}

thead tr.row-groups th {
	padding: 8px 16px 4px;
	font-size: 10px;
	font-weight: 600;
	letter-spacing: 0.1em;
	text-transform: uppercase;
	color: var(--dim);
	text-align: center;
	border-bottom: 1px solid var(--border);
}

thead tr.row-groups th.left  { text-align: left; }
thead tr.row-groups th.right { text-align: right; }
thead tr.row-groups th.sep   { border-left: 1px solid var(--border); }

thead tr.row-subs th {
	padding: 5px 16px 8px;
	font-size: 10px;
	font-weight: 500;
	letter-spacing: 0.08em;
	text-transform: uppercase;
	color: var(--dim);
	text-align: right;
	border-bottom: 1px solid var(--border);
	white-space: nowrap;
}

thead tr.row-subs th.left { text-align: left; }
thead tr.row-subs th.sep  { border-left: 1px solid var(--border); }

tbody tr {
	border-bottom: 1px solid var(--border);
}

tbody tr:last-child {
	border-bottom: none;
}

tbody td {
	padding: 13px 16px;
	font-size: 12px;
	text-align: right;
	font-variant-numeric: tabular-nums;
	vertical-align: middle;
}

tbody td.left { text-align: left; }
tbody td.sep  { border-left: 1px solid var(--border); }

.ref-tag {
	display: inline-block;
	font-size: 10px;
	font-weight: 700;
	letter-spacing: 0.04em;
	padding: 2px 7px;
	border: 1px solid var(--black);
	color: var(--black);
	background: var(--white);
}

.val-strong { font-weight: 700; }
.val-dim    { color: var(--dim); }

.all-failed {
	color: var(--dim);
	font-style: italic;
	font-size: 11px;
}

.failures-badge {
	font-size: 10px;
	font-weight: 600;
	padding: 1px 5px;
	border: 1px solid var(--black);
	color: var(--black);
	margin-left: 6px;
}

.delta          { font-weight: 700; color: var(--black); }
.delta-faster   { color: var(--black); }
.delta-slower   { color: var(--black); }
.delta-na       { color: var(--dim); font-weight: 400; font-style: italic; font-size: 11px; }
.delta-err      { font-size: 11px; font-weight: 400; color: var(--dim); }
.delta-winner   { display: block; font-size: 10px; font-weight: 600; color: var(--dim); letter-spacing: 0.04em; margin-top: 3px; }

.summary {
	border: 1px solid var(--border);
}

.summary-block-header {
	padding: 9px 16px;
	border-bottom: 1px solid var(--border);
	font-size: 11px;
	font-weight: 600;
	letter-spacing: 0.04em;
}

.summary table thead th {
	padding: 6px 16px 8px;
	font-size: 10px;
	font-weight: 500;
	letter-spacing: 0.08em;
	text-transform: uppercase;
	color: var(--dim);
	border-bottom: 1px solid var(--border);
	text-align: left;
}

.summary table thead th:last-child { text-align: right; }

.summary table tbody td {
	padding: 13px 16px;
	font-size: 12px;
	vertical-align: middle;
	text-align: left;
}

.summary table tbody td:last-child { text-align: right; }

.summary table tbody tr {
	border-bottom: 1px solid var(--border);
}

.summary table tbody tr:last-child { border-bottom: none; }

.summary-sentence { color: var(--dim); font-size: 11px; }
.summary-sentence .ref-tag { margin: 0 4px; }

.summary-ratio       { font-size: 14px; font-weight: 800; letter-spacing: -0.02em; }
.summary-uncertainty { font-size: 11px; font-weight: 400; color: var(--dim); margin-left: 4px; }

footer {
	margin-top: 48px;
	padding-top: 16px;
	border-top: 1px solid var(--border);
	font-size: 10px;
	color: var(--dim);
	letter-spacing: 0.04em;
}
</style>
</head>
<body>

<header>
	<div class="header-eyebrow">Benchmark Report</div>
	<div class="header-refs">
		<span>{{.BaseRef}}</span>
		<span class="sep">vs</span>
		<span>{{.HeadRef}}</span>
	</div>
	<div class="header-meta">
		<span>{{.TotalBenchmarks}} benchmark{{if ne .TotalBenchmarks 1}}s{{end}}</span>
		<span>{{.Runs}} runs each</span>
		<span>generated {{.GeneratedAt}}</span>
	</div>
</header>

<div class="benchmarks">
	<div class="section-label">Results</div>

	<div class="benchmark-block">
		<table>
			<thead>
				<tr class="row-groups">
					<th class="left" rowspan="2">benchmark</th>
					<th class="sep" colspan="2">{{.BaseRef}}</th>
					<th class="sep" colspan="2">{{.HeadRef}}</th>
					<th class="sep right" rowspan="2">delta</th>
				</tr>
				<tr class="row-subs">
					<th class="sep">mean</th>
					<th class="sep">stddev</th>
					<th class="sep">mean</th>
					<th class="sep">stddev</th>
				</tr>
			</thead>
			<tbody>
				{{range .Groups}}
				<tr>
					<td class="left">{{.Label}}</td>

					{{if .Base.AllFailed}}
					<td class="sep" colspan="2"><span class="all-failed">all runs failed</span></td>
					{{else}}
					<td class="sep">
						<span class="val-strong">{{.Base.MeanFmt}}</span>
						<span class="val-dim"> ms</span>
						{{if .Base.HasFailures}}<span class="failures-badge">{{.Base.Failures}} failed</span>{{end}}
					</td>
					<td class="sep"><span class="val-dim">{{.Base.StddevFmt}} ms</span></td>
					{{end}}

					{{if .Head.NoData}}
					<td class="sep" colspan="2"><span class="all-failed">no data</span></td>
					{{else if .Head.AllFailed}}
					<td class="sep" colspan="2"><span class="all-failed">all runs failed</span></td>
					{{else}}
					<td class="sep">
						<span class="val-strong">{{.Head.MeanFmt}}</span>
						<span class="val-dim"> ms</span>
						{{if .Head.HasFailures}}<span class="failures-badge">{{.Head.Failures}} failed</span>{{end}}
					</td>
					<td class="sep"><span class="val-dim">{{.Head.StddevFmt}} ms</span></td>
					{{end}}

					<td class="sep">
						{{if .Delta.Available}}
						<span class="delta">{{.Delta.RatioFmt}}×</span>
						<span class="delta-err"> ± {{.Delta.UncertaintyFmt}}</span>
						<span class="delta-winner">{{.Delta.WinnerRef}}</span>
						{{else}}
						<span class="delta-na">—</span>
						{{end}}
					</td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</div>
</div>

{{if .Summaries}}
<div class="summary">
	<div class="summary-block-header">Summary</div>
	<table>
		<thead>
			<tr>
				<th>comparison</th>
				<th>speedup</th>
			</tr>
		</thead>
		<tbody>
			{{range .Summaries}}
			<tr>
				<td>
					<span class="summary-sentence">
						<span class="ref-tag">{{.FasterRef}}</span>
						({{.Label}}) is faster than
						<span class="ref-tag">{{.SlowerRef}}</span>
					</span>
				</td>
				<td>
					<span class="summary-ratio">{{.RatioFmt}}×</span>
					<span class="summary-uncertainty">± {{.UncertaintyFmt}}</span>
				</td>
			</tr>
			{{end}}
		</tbody>
	</table>
</div>
{{end}}

<footer>bal-bench &mdash; benchmark report</footer>

</body>
</html>`

// htmlReportData is the top-level data passed to reportHTMLTemplate.
type htmlReportData struct {
	BaseRef         string
	HeadRef         string
	TotalBenchmarks int
	Runs            int
	GeneratedAt     string
	Groups          []htmlGroup
	Summaries       []htmlSummary
}

// htmlGroup holds the rendered stats for one benchmark target (one file/package).
type htmlGroup struct {
	Label string
	Base  htmlStats
	Head  htmlStats
	Delta htmlDelta
}

// htmlStats holds pre-formatted strings for a single ref's measurements.
type htmlStats struct {
	MeanFmt     string
	StddevFmt   string
	Failures    int
	HasFailures bool
	AllFailed   bool
	NoData      bool
}

// htmlDelta holds the pre-formatted speedup ratio between base and head.
type htmlDelta struct {
	Available      bool
	HeadFaster     bool
	WinnerRef      string // the ref that ran faster
	RatioFmt       string
	UncertaintyFmt string
}

// htmlSummary holds one row of the summary table.
type htmlSummary struct {
	Label          string
	FasterRef      string
	SlowerRef      string
	RatioFmt       string
	UncertaintyFmt string
}

func toHTMLStats(s fileStats) htmlStats {
	return htmlStats{
		MeanFmt:     fmt.Sprintf("%.1f", s.mean),
		StddevFmt:   fmt.Sprintf("%.1f", s.stddev),
		Failures:    s.failures,
		HasFailures: s.failures > 0,
		AllFailed:   s.failures == s.runs,
	}
}

func buildHTMLReport(r *benchmarkResult) htmlReportData {
	headByLabel := make(map[string]fileStats, len(r.head))
	for _, s := range r.head {
		headByLabel[s.label] = s
	}

	// Infer shared run count from base; all targets run the same number of times.
	runs := 0
	if len(r.base) > 0 {
		runs = r.base[0].runs
	}

	var groups []htmlGroup
	for _, base := range r.base {
		head, hasHead := headByLabel[base.label]

		var headStats htmlStats
		if hasHead {
			headStats = toHTMLStats(head)
		} else {
			headStats = htmlStats{NoData: true}
		}

		delta := htmlDelta{}
		if hasHead && base.mean > 0 && base.failures < base.runs && head.failures < head.runs {
			ratio, ratioErr := speedupRatio(base, head)
			// ratio >= 1.0 means head is faster (base.mean / head.mean >= 1)
			headFaster := ratio >= 1.0
			winnerRef := r.headRef
			if !headFaster {
				winnerRef = r.baseRef
			}
			if ratio < 1.0 {
				ratio = 1.0 / ratio
				ratioErr = ratioErr / (ratio * ratio)
			}
			delta = htmlDelta{
				Available:      true,
				HeadFaster:     headFaster,
				WinnerRef:      winnerRef,
				RatioFmt:       fmt.Sprintf("%.2f", ratio),
				UncertaintyFmt: fmt.Sprintf("%.2f", math.Abs(ratioErr)),
			}
		}

		groups = append(groups, htmlGroup{
			Label: base.label,
			Base:  toHTMLStats(base),
			Head:  headStats,
			Delta: delta,
		})
	}

	var summaries []htmlSummary
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

		summaries = append(summaries, htmlSummary{
			Label:          base.label,
			FasterRef:      fasterRef,
			SlowerRef:      slowerRef,
			RatioFmt:       fmt.Sprintf("%.2f", ratio),
			UncertaintyFmt: fmt.Sprintf("%.2f", math.Abs(ratioErr)),
		})
	}

	return htmlReportData{
		BaseRef:         r.baseRef,
		HeadRef:         r.headRef,
		TotalBenchmarks: len(groups),
		Runs:            runs,
		GeneratedAt:     time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
		Groups:          groups,
		Summaries:       summaries,
	}
}

func writeHTMLReport(r *benchmarkResult, outputPath string) error {
	tmpl, err := template.New("report").Parse(reportHTMLTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse html template: %w", err)
	}

	data := buildHTMLReport(r)

	var sb strings.Builder
	if err := tmpl.Execute(&sb, data); err != nil {
		return fmt.Errorf("failed to render html report: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(sb.String()), 0o644); err != nil {
		return fmt.Errorf("failed to write html report to %q: %w", outputPath, err)
	}

	return nil
}
