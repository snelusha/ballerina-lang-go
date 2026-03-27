package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type inputMode byte

const (
	modeSingleFile inputMode = iota
	modePackage
	modeDirectory
)

// Target describes what will be passed to each compiler invocation.
type Target struct {
	Mode  inputMode
	Paths []string
	Label string
}

// resolveTarget inspects the given path and returns the appropriate Target.
//
//   - A .bal file      → SingleFile
//   - Dir + Ballerina.toml → Package (compiled as a unit)
//   - Dir without toml → Directory (each .bal file benchmarked individually)
func resolveTarget(input string) (*Target, error) {
	info, err := os.Stat(input)
	if err != nil {
		return nil, fmt.Errorf("input path %q not found: %w", input, err)
	}

	if !info.IsDir() {
		if filepath.Ext(input) != ".bal" {
			return nil, fmt.Errorf("input file %q must have a .bal extension", input)
		}
		return &Target{
			Mode:  modeSingleFile,
			Paths: []string{input},
			Label: filepath.Base(input),
		}, nil
	}

	if _, err := os.Stat(filepath.Join(input, "Ballerina.toml")); err == nil {
		return &Target{
			Mode:  modePackage,
			Paths: []string{input},
			Label: fmt.Sprintf("%s (package)", filepath.Base(input)),
		}, nil
	}

	files, err := collectBalFiles(input)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no .bal files found in directory %q", input)
	}

	return &Target{
		Mode:  modeDirectory,
		Paths: files,
		Label: fmt.Sprintf("%s/", filepath.Base(input)),
	}, nil
}

func collectBalFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".bal" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error scanning directory %q: %w", dir, err)
	}
	return files, nil
}
