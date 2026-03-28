package main

import (
	"fmt"
	"html/template"
	"math"
	"os"
	"strings"
	"time"
)

const htmlTemplate = `<!DOCTYPE html>
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
	--dim: #888888;
	--faster: #000000;
	--slower: #888888;
	--failed: #000000;
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
	max-width: 980px;
	margin: 0 auto;
}

header {
	border-bottom: 1px solid var(--border);
	padding-bottom: 24px;
	margin-bottom: 40px;
}

.header-title {
	font-size: 11px;
	font-weight: 400;
	color: var(--dim);
	letter-spacing: 0.08em;
	text-transform: uppercase;
	margin-bottom: 8px;
}

.header-refs {
	font-size: 18px;
	font-weight: 700;
	letter-spacing: -0.02em;
}

.header-refs .sep {
	font-weight: 300;
	color: var(--dim);
	margin: 0 10px;
}

.header-meta {
	margin-top: 8px;
	font-size: 11px;
	color: var(--dim);
	display: flex;
	gap: 24px;
}

.section-label {
	font-size: 10px;
	font-weight: 600;
	letter-spacing: 0.12em;
	text-transform: uppercase;
	color: var(--dim);
	margin-bottom: 12px;
}

.benchmarks {
	margin-bottom: 48px;
}

.benchmark-group {
	margin-bottom: 32px;
	border: 1px solid var(--border);
}

.benchmark-group-header {
	padding: 10px 16px;
	background: var(--black);
	color: var(--white);
	font-size: 11px;
	font-weight: 600;
	letter-spacing: 0.06em;
}

table {
	width: 100%;
	border-collapse: collapse;
}

thead tr {
	border-bottom: 1px solid var(--border);
}

thead th {
	padding: 10px 16px;
	text-align: left;
	font-size: 10px;
	font-weight: 600;
	letter-spacing: 0.1em;
	text-transform: uppercase;
	color: var(--dim);
	white-space: nowrap;
}

thead th:not(:first-child) {
	text-align: right;
}

tbody tr {
	border-bottom: 1px solid var(--border);
}

tbody tr:last-child {
	border-bottom: none;
}

tbody td {
	padding: 12px 16px;
	font-size: 12px;
	vertical-align: middle;
}

tbody td:not(:first-child) {
	text-align: right;
	font-variant-numeric: tabular-nums;
}

.ref-tag {
	display: inline-block;
	font-size: 10px;
	font-weight: 600;
	letter-spacing: 0.04em;
	padding: 2px 6px;
	background: var(--black);
	color: var(--white);
}

.ref-tag.head {
	background: var(--white);
	color: var(--black);
	border: 1px solid var(--black);
}

.val-primary {
	font-weight: 700;
}

.val-secondary {
	color: var(--dim);
}

.failures-badge {
	display: inline-block;
	font-size: 10px;
	font-weight: 600;
	padding: 1px 5px;
	background: var(--black);
	color: var(--white);
	margin-left: 6px;
}

.all-failed {
	color: var(--dim);
	font-style: italic;
	font-size: 11px;
}

.summary {
	border: 1px solid var(--border);
}

.summary-header {
	padding: 10px 16px;
	background: var(--black);
	color: var(--white);
	font-size: 11px;
	font-weight: 600;
	letter-spacing: 0.06em;
}

.summary-row {
	display: flex;
	align-items: baseline;
	justify-content: space-between;
	padding: 14px 16px;
	border-bottom: 1px solid var(--border);
}

.summary-row:last-child {
	border-bottom: none;
}

.summary-label {
	font-size: 11px;
	color: var(--dim);
	flex-shrink: 0;
	margin-right: 16px;
}

.summary-label .ref-tag {
	margin: 0 3px;
}

.summary-value {
	font-size: 14px;
	font-weight: 800;
	letter-spacing: -0.02em;
	white-space: nowrap;
}

.summary-value.faster {
	color: var(--black);
}

.summary-value.slower {
	color: var(--dim);
}

footer {
	margin-top: 40px;
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
	<div class="header-title">Benchmark Report</div>
	<div class="header-refs">
		<span>{{.BaseRef}}</span>
		<span class="sep">vs</span>
		<span>{{.HeadRef}}</span>
	</div>
	<div class="header-meta">
		<span>{{.TotalBenchmarks}} benchmark{{if ne .TotalBenchmarks 1}}s{{end}}</span>
		<span>generated {{.GeneratedAt}}</span>
	</div>
</header>

<div class="benchmarks">
	<div class="section-label">Results</div>
	{{range .Groups}}
	<div class="benchmark-group">
		<div class="benchmark-group-header">{{.Label}}</div>
		<table>
			<thead>
				<tr>
					<th>ref</th>
					<th>mean</th>
					<th>σ</th>
					<th>min</th>
					<th>max</th>
					<th>runs</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td><span class="ref-tag">{{$.BaseRef}}</span></td>
					{{if .Base.AllFailed}}
					<td colspan="5"><span class="all-failed">all {{.Base.Runs}} runs failed</span></td>
					{{else}}
					<td><span class="val-primary">{{.Base.MeanFmt}}</span> <span class="val-secondary">ms</span></td>
					<td><span class="val-secondary">{{.Base.StddevFmt}} ms</span></td>
					<td><span class="val-secondary">{{.Base.MinFmt}} ms</span></td>
					<td><span class="val-secondary">{{.Base.MaxFmt}} ms</span></td>
					<td><span class="val-secondary">{{.Base.SuccessRuns}}{{if .Base.HasFailures}}<span class="failures-badge">{{.Base.Failures}} failed</span>{{end}}</span></td>
					{{end}}
				</tr>
				<tr>
					<td><span class="ref-tag head">{{$.HeadRef}}</span></td>
					{{if .Head.AllFailed}}
					<td colspan="5"><span class="all-failed">all {{.Head.Runs}} runs failed</span></td>
					{{else if .Head.NoData}}
					<td colspan="5"><span class="all-failed">no data</span></td>
					{{else}}
					<td><span class="val-primary">{{.Head.MeanFmt}}</span> <span class="val-secondary">ms</span></td>
					<td><span class="val-secondary">{{.Head.StddevFmt}} ms</span></td>
					<td><span class="val-secondary">{{.Head.MinFmt}} ms</span></td>
					<td><span class="val-secondary">{{.Head.MaxFmt}} ms</span></td>
					<td><span class="val-secondary">{{.Head.SuccessRuns}}{{if .Head.HasFailures}}<span class="failures-badge">{{.Head.Failures}} failed</span>{{end}}</span></td>
					{{end}}
				</tr>
			</tbody>
		</table>
	</div>
	{{end}}
</div>

{{if .Summaries}}
<div class="summary">
	<div class="summary-header">Summary</div>
	{{range .Summaries}}
	<div class="summary-row">
		<div class="summary-label">
			<span class="ref-tag{{if .HeadFaster}} head{{end}}">{{.FasterRef}}</span>
			&nbsp;({{.Label}}) is faster than&nbsp;
			<span class="ref-tag{{if not .HeadFaster}}{{else}} head{{end}}">{{.SlowerRef}}</span>
		</div>
		<div class="summary-value{{if .HeadFaster}} faster{{else}} slower{{end}}">
			{{.RatioFmt}}×&nbsp;<span style="font-weight:400;font-size:11px;color:#888">± {{.UncertaintyFmt}}</span>
		</div>
	</div>
	{{end}}
</div>
{{end}}

<footer>bal-bench &mdash; benchmark report</footer>

</body>
</html>`

type htmlReportData struct {
	BaseRef         string
	HeadRef         string
	TotalBenchmarks int
	GeneratedAt     string
	Groups          []htmlGroup
	Summaries       []htmlSummary
}

type htmlGroup struct {
	Label string
	Base  htmlStats
	Head  htmlStats
}

type htmlStats struct {
	MeanFmt     string
	StddevFmt   string
	MinFmt      string
	MaxFmt      string
	Runs        int
	Failures    int
	SuccessRuns int
	HasFailures bool
	AllFailed   bool
	NoData      bool
}

type htmlSummary struct {
	Label          string
	FasterRef      string
	SlowerRef      string
	HeadFaster     bool
	RatioFmt       string
	UncertaintyFmt string
}

func toHTMLStats(s fileStats) htmlStats {
	return htmlStats{
		MeanFmt:     fmt.Sprintf("%.1f", s.mean),
		StddevFmt:   fmt.Sprintf("%.1f", s.stddev),
		MinFmt:      fmt.Sprintf("%.1f", s.min),
		MaxFmt:      fmt.Sprintf("%.1f", s.max),
		Runs:        s.runs,
		Failures:    s.failures,
		SuccessRuns: s.runs - s.failures,
		HasFailures: s.failures > 0,
		AllFailed:   s.failures == s.runs,
	}
}

func buildHTMLReport(r *benchmarkResult) htmlReportData {
	headByLabel := make(map[string]fileStats, len(r.head))
	for _, s := range r.head {
		headByLabel[s.label] = s
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

		groups = append(groups, htmlGroup{
			Label: base.label,
			Base:  toHTMLStats(base),
			Head:  headStats,
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
		headFaster := true
		if ratio < 1.0 {
			ratio = 1.0 / ratio
			ratioErr = ratioErr / (ratio * ratio)
			fasterRef, slowerRef = r.baseRef, r.headRef
			headFaster = false
		}

		summaries = append(summaries, htmlSummary{
			Label:          base.label,
			FasterRef:      fasterRef,
			SlowerRef:      slowerRef,
			HeadFaster:     headFaster,
			RatioFmt:       fmt.Sprintf("%.2f", ratio),
			UncertaintyFmt: fmt.Sprintf("%.2f", math.Abs(ratioErr)),
		})
	}

	return htmlReportData{
		BaseRef:         r.baseRef,
		HeadRef:         r.headRef,
		TotalBenchmarks: len(groups),
		GeneratedAt:     time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
		Groups:          groups,
		Summaries:       summaries,
	}
}

func writeHTMLReport(r *benchmarkResult, outputPath string) error {
	funcMap := template.FuncMap{
		"not": func(b bool) bool { return !b },
	}

	tmpl, err := template.New("report").Funcs(funcMap).Parse(htmlTemplate)
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
