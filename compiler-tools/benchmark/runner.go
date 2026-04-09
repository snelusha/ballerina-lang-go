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

	switch tgt.mode {
	case multipleFilesMode:
		return b.runDirectoryBenchmarks(baseTree, headTree, tgt)
	default:
		cmds, err := b.benchmarkCommands(baseTree, headTree, tgt)
		if err != nil {
			return err
		}
		return b.runHyperfine(cmds)
	}
}

func (b *benchmark) runDirectoryBenchmarks(baseTree, headTree string, tgt *benchmarkTarget) error {
	paths := append([]string(nil), tgt.paths...)
	sort.Strings(paths)

	for _, p := range paths {
		cmds := []string{
			benchCommand(filepath.Join(baseTree, baseInterpreterName), p),
			benchCommand(filepath.Join(headTree, headInterpreterName), p),
		}
		if err := b.runHyperfine(cmds); err != nil {
			return err
		}
	}
	return nil
}

func (b *benchmark) benchmarkCommands(baseTree, headTree string, tgt *benchmarkTarget) ([]string, error) {
	switch tgt.mode {
	case singleFileMode, packageMode:
		base := benchCommand(filepath.Join(baseTree, baseInterpreterName), tgt.paths[0])
		head := benchCommand(filepath.Join(headTree, headInterpreterName), tgt.paths[0])
		return []string{base, head}, nil
	default:
		return nil, fmt.Errorf("unknown target mode %v", tgt.mode)
	}
}

func (b *benchmark) runHyperfine(cmds []string) error {
	args := b.hyperfineFlags()
	args = append(args, cmds...)
	return runCmd(".", "hyperfine", args...)
}

func benchCommand(interpreter, balPath string) string {
	return fmt.Sprintf("%s run %s", shellQuote(interpreter), shellQuote(balPath))
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
	path := filepath.Join(b.workRoot, "worktree-"+sanitizeRef(ref))
	b.removeWorktree(path)

	if err := runCmd(".", "git", "worktree", "add", "--detach", path, ref); err != nil {
		return "", fmt.Errorf("checkout worktree for ref %q: %w", ref, err)
	}
	return path, nil
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
