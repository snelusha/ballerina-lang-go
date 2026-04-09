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
	"flag"
	"fmt"
	"os"
)

type config struct {
	baseRef string
	headRef string
	target  string
	warmup  int
	runs    int
	web     bool
}

func parseConfig() (*config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <base-ref> <head-ref> <target>\n\n", os.Args[0])
		fs.PrintDefaults()
	}

	warmup := fs.Int("warmup", 0, "number of warmup runs (passed to hyperfine)")
	runs := fs.Int("runs", 10, "number of benchmark runs (passed to hyperfine)")
	web := fs.Bool("web", false, "export a HTML report (via hyperfine JSON export)")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	args := fs.Args()
	if len(args) < 3 {
		fs.Usage()
		return nil, fmt.Errorf("missing required arguments")
	}

	return &config{
		baseRef: args[0],
		headRef: args[1],
		target:  args[2],
		warmup:  *warmup,
		runs:    *runs,
		web:     *web,
	}, nil
}

func (c *config) validate() error {
	if c.baseRef == "" {
		return fmt.Errorf("baseRef is required")
	}
	if c.headRef == "" {
		return fmt.Errorf("headRef is required")
	}
	if c.baseRef == c.headRef {
		return fmt.Errorf("baseRef and headRef cannot be the same")
	}
	if c.target == "" {
		return fmt.Errorf("target is required")
	}
	if _, err := os.Stat(c.target); os.IsNotExist(err) {
		return fmt.Errorf("target does not exist: %s", c.target)
	}
	if c.warmup < 0 {
		return fmt.Errorf("warmup must be non-negative")
	}
	if c.runs <= 0 {
		return fmt.Errorf("runs must be greater than zero")
	}
	return nil
}
