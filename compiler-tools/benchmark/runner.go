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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	baseInterpreterName = "bal-base"
	headInterpreterName = "bal-head"
)

type benchmark struct {
	config
	workRoot string
}

func (b *benchmark) run() error {
	workRoot, err := os.MkdirTemp("", "bal-bench-*")
	if err != nil {
		return fmt.Errorf("create work directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(workRoot) }()
	b.workRoot = workRoot

	baseTree, err := b.checkoutWorktree(b.baseRef)
	if err != nil {
		return err
	}
	defer b.removeWorktree(baseTree)

	headTree, err := b.checkoutWorktree(b.headRef)
	if err != nil {
		return err
	}
	defer b.removeWorktree(headTree)

	if err := b.buildInterpreter(b.baseRef, baseTree, baseInterpreterName); err != nil {
		return err
	}
	if err := b.buildInterpreter(b.headRef, headTree, headInterpreterName); err != nil {
		return err
	}

	tgt, err := resolveTarget(b.target)
	if err != nil {
		return fmt.Errorf("resolve target: %w", err)
	}

	if !b.web {
		switch tgt.mode {
		case multipleFilesMode:
			return b.runDirectoryBenchmarks(baseTree, headTree, tgt)
		default:
			cmds, err := b.benchmarkCommands(baseTree, headTree, tgt)
			if err != nil {
				return err
			}
			_, err = b.runHyperfine(cmds, "")
			return err
		}
	}

	return b.runWebReport(baseTree, headTree, tgt)
}

func (b *benchmark) runDirectoryBenchmarks(baseTree, headTree string, tgt *benchmarkTarget) error {
	paths := append([]string(nil), tgt.paths...)
	sort.Strings(paths)

	for _, p := range paths {
		cmds := benchmarkPair(baseTree, headTree, p)
		if _, err := b.runHyperfine(cmds, ""); err != nil {
			return err
		}
	}
	return nil
}

func (b *benchmark) benchmarkCommands(baseTree, headTree string, tgt *benchmarkTarget) ([]string, error) {
	switch tgt.mode {
	case singleFileMode, packageMode:
		return benchmarkPair(baseTree, headTree, tgt.paths[0]), nil
	default:
		return nil, fmt.Errorf("unknown target mode %v", tgt.mode)
	}
}

func (b *benchmark) runHyperfine(cmds []string, jsonOut string) (*hyperfineExport, error) {
	args := b.hyperfineFlags()
	if jsonOut != "" {
		args = append(args, "--export-json", jsonOut)
	}
	args = append(args, cmds...)
	if err := runCmd(".", "hyperfine", args...); err != nil {
		return nil, err
	}
	if jsonOut == "" {
		return nil, nil
	}
	exp, err := readHyperfineExport(jsonOut)
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func benchCommand(interpreter, balPath string) string {
	return fmt.Sprintf("%s run %s", shellQuote(interpreter), shellQuote(balPath))
}

func benchmarkPair(baseTree, headTree, balPath string) []string {
	return []string{
		benchCommand(filepath.Join(baseTree, baseInterpreterName), balPath),
		benchCommand(filepath.Join(headTree, headInterpreterName), balPath),
	}
}

func shellQuote(s string) string {
	if s == "" {
		return "''"
	}
	if !strings.ContainsAny(s, " \t\n'\"$\\`") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}

func (b *benchmark) hyperfineFlags() []string {
	var args []string
	if b.warmup > 0 {
		args = append(args, "--warmup", strconv.Itoa(b.warmup))
	}
	if b.runs > 0 {
		args = append(args, "--runs", strconv.Itoa(b.runs))
	}
	return args
}

func (b *benchmark) checkoutWorktree(ref string) (string, error) {
	if err := ensureRefAvailable(ref); err != nil {
		return "", err
	}

	path := filepath.Join(b.workRoot, "worktree-"+sanitizeRef(ref))
	b.removeWorktree(path)

	if err := runCmd(".", "git", "worktree", "add", "--detach", path, ref); err != nil {
		return "", fmt.Errorf("checkout worktree for ref %q: %w", ref, err)
	}
	return path, nil
}

func ensureRefAvailable(ref string) error {
	if err := runCmdSilent(".", "git", "rev-parse", "--verify", ref+"^{commit}"); err == nil {
		return nil
	}

	if err := runCmd(".", "git", "fetch", "--all", "--tags"); err != nil {
		return fmt.Errorf("failed to fetch tags while resolving ref %q: %w", ref, err)
	}

	if err := runCmdSilent(".", "git", "rev-parse", "--verify", ref+"^{commit}"); err != nil {
		return fmt.Errorf("git ref %q not found locally (even after fetching tags)", ref)
	}
	return nil
}

func sanitizeRef(ref string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':':
			return '-'
		default:
			return r
		}
	}, ref)
}

func sanitizeFilename(s string) string {
	if s == "" {
		return "output"
	}
	s = strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':':
			return '-'
		default:
			return r
		}
	}, s)
	s = strings.TrimSpace(s)
	if s == "" {
		return "output"
	}
	return s
}

func (b *benchmark) runWebReport(baseTree, headTree string, tgt *benchmarkTarget) error {
	outDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}
	jsonDir, err := os.MkdirTemp("", "bal-bench-json-*")
	if err != nil {
		return fmt.Errorf("create temp json directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(jsonDir) }()

	stamp := time.Now().Format("20060102-150405")
	outBase := sanitizeFilename(tgt.label) + "-" + stamp

	var runs []webRun

	switch tgt.mode {
	case multipleFilesMode:
		paths := append([]string(nil), tgt.paths...)
		sort.Strings(paths)

		for _, p := range paths {
			display := p
			if tgt.root != "" {
				if rel, err := filepath.Rel(tgt.root, p); err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
					display = rel
				}
			}

			jsonOut := filepath.Join(jsonDir, outBase+"-"+sanitizeFilename(display)+".json")
			cmds := benchmarkPair(baseTree, headTree, p)
			exp, err := b.runHyperfine(cmds, jsonOut)
			if err != nil {
				return err
			}
			runs = append(runs, webRun{Label: display, Export: exp})
		}
	default:
		jsonOut := filepath.Join(jsonDir, outBase+".json")
		cmds, err := b.benchmarkCommands(baseTree, headTree, tgt)
		if err != nil {
			return err
		}
		exp, err := b.runHyperfine(cmds, jsonOut)
		if err != nil {
			return err
		}
		runs = append(runs, webRun{Label: tgt.label, Export: exp})
	}

	htmlOut := filepath.Join(outDir, outBase+".html")
	report := webReport{
		Title:     "Ballerina benchmark",
		BaseRef:   b.baseRef,
		HeadRef:   b.headRef,
		Target:    tgt.label,
		Generated: time.Now(),
		Runs:      runs,
	}
	if err := writeHTMLReport(htmlOut, report); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "web report: %s\n", htmlOut)
	return nil
}

func (b *benchmark) removeWorktree(path string) {
	_ = runCmdSilent(".", "git", "worktree", "remove", "--force", path)
}

func (b *benchmark) buildInterpreter(ref, repoRoot, output string) error {
	if err := runCmd(repoRoot, "go", "build", "-o", output, "./cli/cmd"); err != nil {
		return fmt.Errorf("build interpreter for ref %q: %w", ref, err)
	}
	return nil
}

func runCmd(dir, name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func runCmdSilent(dir, name string, args ...string) error {
	c := exec.Command(name, args...)
	c.Dir = dir
	c.Stdout = nil
	c.Stderr = nil
	return c.Run()
}
