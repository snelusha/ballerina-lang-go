// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"os"
	"time"
)

type hyperfineExport struct {
	Results []hyperfineResult `json:"results"`
}

type hyperfineResult struct {
	Command string  `json:"command"`
	Mean    float64 `json:"mean"`
	Stddev  float64 `json:"stddev"`
	Median  float64 `json:"median"`
}

func readHyperfineExport(path string) (*hyperfineExport, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read hyperfine json %q: %w", path, err)
	}
	var exp hyperfineExport
	if err := json.Unmarshal(b, &exp); err != nil {
		return nil, fmt.Errorf("parse hyperfine json %q: %w", path, err)
	}
	if len(exp.Results) == 0 {
		return nil, fmt.Errorf("hyperfine json %q contains no results", path)
	}
	return &exp, nil
}

type webReport struct {
	Title     string
	BaseRef   string
	HeadRef   string
	Target    string
	Generated time.Time
	Runs      []webRun
}

type webRun struct {
	Label  string
	Export *hyperfineExport
}

type webRow struct {
	Label            string
	Base             *hyperfineResult
	Head             *hyperfineResult
	BaseMeanMS       string
	BaseStddevMS     string
	HeadMeanMS       string
	HeadStddevMS     string
	DeltaAvailable   bool
	DeltaRatio       string
	DeltaUncertainty string
	DeltaWinnerRef   string
}

func writeHTMLReport(outPath string, report webReport) error {
	rows := make([]webRow, 0, len(report.Runs))
	for _, r := range report.Runs {
		var base, head *hyperfineResult
		if r.Export != nil && len(r.Export.Results) >= 2 {
			base = &r.Export.Results[0]
			head = &r.Export.Results[1]
		}
		rows = append(rows, webRow{
			Label: r.Label,
			Base:  base,
			Head:  head,
		})
		last := &rows[len(rows)-1]
		if base != nil {
			last.BaseMeanMS = fmt.Sprintf("%.3f", base.Mean*1000.0)
			last.BaseStddevMS = fmt.Sprintf("%.3f", base.Stddev*1000.0)
		}
		if head != nil {
			last.HeadMeanMS = fmt.Sprintf("%.3f", head.Mean*1000.0)
			last.HeadStddevMS = fmt.Sprintf("%.3f", head.Stddev*1000.0)
		}
		last.DeltaAvailable, last.DeltaRatio, last.DeltaUncertainty, last.DeltaWinnerRef = formatDelta(base, head, report.BaseRef, report.HeadRef)
	}

	tpl := template.Must(template.New("report").Parse(htmlTemplate))
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create html report %q: %w", outPath, err)
	}
	defer func() { _ = f.Close() }()

	data := struct {
		Report webReport
		Rows   []webRow
	}{
		Report: report,
		Rows:   rows,
	}

	if err := tpl.Execute(f, data); err != nil {
		return fmt.Errorf("render html report %q: %w", outPath, err)
	}
	return nil
}

func formatDelta(base, head *hyperfineResult, baseRef, headRef string) (bool, string, string, string) {
	if base == nil || head == nil {
		return false, "", "", ""
	}
	if base.Mean <= 0 || head.Mean <= 0 {
		return false, "", "", ""
	}

	rawRatio := base.Mean / head.Mean
	rawRatioErr := rawRatio * math.Sqrt(
		math.Pow(base.Stddev/base.Mean, 2)+math.Pow(head.Stddev/head.Mean, 2),
	)

	winnerRef := headRef
	ratio := rawRatio
	ratioErr := rawRatioErr
	if rawRatio < 1.0 {
		winnerRef = baseRef
		ratio = 1.0 / rawRatio
		// If q = 1/r then dq = dr / r^2.
		ratioErr = rawRatioErr / (rawRatio * rawRatio)
	}
	return true, fmt.Sprintf("%.2f", ratio), fmt.Sprintf("%.2f", math.Abs(ratioErr)), winnerRef
}

const htmlTemplate = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>{{ .Report.Title }}</title>
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
		max-width: 1080px;
		margin: 0 auto;
	}

	header {
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
		margin-top: 8px;
		font-size: 11px;
		color: var(--dim);
		display: flex;
		gap: 20px;
		flex-wrap: wrap;
	}

	.section-label {
		font-size: 10px;
		font-weight: 600;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		color: var(--dim);
		margin-bottom: 16px;
	}

	.table-wrap {
		border: 1px solid var(--border);
		overflow-x: auto;
	}

	table {
		width: 100%;
		border-collapse: collapse;
		min-width: 860px;
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
		padding: 8px 16px 8px;
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

	.val-strong {
		font-size: 16px;
		font-weight: 700;
		color: var(--black);
	}

	.delta {
		font-size: 16px;
		font-weight: 700;
		color: var(--black);
	}

	.delta-err {
		font-size: 14px;
		font-weight: 400;
		color: var(--black);
	}

	.delta-winner {
		display: block;
		font-size: 12px;
		color: var(--black);
		letter-spacing: 0.04em;
		margin-top: 4px;
	}

	.not-available {
		color: var(--dim);
		font-style: italic;
		font-size: 11px;
	}
	</style>
</head>
<body>
	<header>
		<div class="header-eyebrow">{{ .Report.Title }}</div>
		<div class="header-refs">
			<span>{{ .Report.BaseRef }}</span>
			<span class="sep">vs</span>
			<span>{{ .Report.HeadRef }}</span>
		</div>
		<div class="header-meta">
			<span>generated {{ .Report.Generated.Format "2006-01-02 15:04:05" }}</span>
		</div>
	</header>

	<div class="section-label">Results</div>

	<div class="table-wrap">
		<table>
			<thead>
				<tr class="row-groups">
					<th class="left" rowspan="2">case</th>
					<th class="sep" colspan="2">{{ .Report.BaseRef }}</th>
					<th class="sep" colspan="2">{{ .Report.HeadRef }}</th>
					<th class="sep right" rowspan="2">delta</th>
				</tr>
				<tr class="row-subs">
					<th class="sep">mean (ms)</th>
					<th class="sep">stddev (ms)</th>
					<th class="sep">mean (ms)</th>
					<th class="sep">stddev (ms)</th>
				</tr>
			</thead>
			<tbody>
				{{ range .Rows }}
				<tr>
					<td class="left">{{ .Label }}</td>
					<td class="sep">{{ if .Base }}<span class="val-strong">{{ .BaseMeanMS }}</span>{{ else }}<span class="not-available">n/a</span>{{ end }}</td>
					<td class="sep">{{ if .Base }}<span class="val-strong">{{ .BaseStddevMS }}</span>{{ else }}<span class="not-available">n/a</span>{{ end }}</td>
					<td class="sep">{{ if .Head }}<span class="val-strong">{{ .HeadMeanMS }}</span>{{ else }}<span class="not-available">n/a</span>{{ end }}</td>
					<td class="sep">{{ if .Head }}<span class="val-strong">{{ .HeadStddevMS }}</span>{{ else }}<span class="not-available">n/a</span>{{ end }}</td>
					<td class="sep">
						{{ if .DeltaAvailable }}
						<span class="delta">{{ .DeltaRatio }}×</span>
						<span class="delta-err"> ± {{ .DeltaUncertainty }}</span>
						<span class="delta-winner">{{ .DeltaWinnerRef }}</span>
						{{ else }}
						<span class="not-available">n/a</span>
						{{ end }}
					</td>
				</tr>
				{{ end }}
			</tbody>
		</table>
	</div>
</body>
</html>`
