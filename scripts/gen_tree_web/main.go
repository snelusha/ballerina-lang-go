package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type FileNode struct {
	Kind     string     `json:"kind"`
	Name     string     `json:"name"`
	Language string     `json:"language,omitempty"`
	Content  string     `json:"content,omitempty"`
	Children []FileNode `json:"children"`
}

func getLanguage(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".bal":
		return "ballerina"
	case ".toml":
		return "toml"
	case ".txtar":
		return "txtar"
	default:
		return ""
	}
}

func buildTree(rootPath string) (FileNode, bool, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return FileNode{}, false, err
	}

	node := FileNode{
		Name: filepath.Base(rootPath),
	}

	if info.IsDir() {
		node.Kind = "dir"
		entries, err := os.ReadDir(rootPath)
		if err != nil {
			return FileNode{}, false, err
		}

		for _, entry := range entries {
			childPath := filepath.Join(rootPath, entry.Name())
			childNode, skip, err := buildTree(childPath)
			if err != nil {
				return FileNode{}, false, err
			}
			if skip {
				continue
			}
			node.Children = append(node.Children, childNode)
		}
		if node.Children == nil {
			node.Children = make([]FileNode, 0)
		}
		// If this directory (recursively) does not contain any files,
		// skip it so the UI does not show empty folders.
		if !hasFiles(node) {
			return FileNode{}, true, nil
		}
	} else {
		node.Kind = "file"
		node.Language = getLanguage(node.Name)
		content, err := os.ReadFile(rootPath)
		if err != nil {
			return FileNode{}, false, err
		}
		node.Content = string(content)
	}

	return node, false, nil
}

func hasFiles(node FileNode) bool {
	if node.Kind == "file" {
		return true
	}
	return slices.ContainsFunc(node.Children, hasFiles)
}

func main() {
	var roots []FileNode

	// Add all example projects under web/tree (e.g. fibonacci)
	examplesRoot := "web/tree"
	exampleEntries, err := os.ReadDir(examplesRoot)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error reading examples directory %s: %v\n", examplesRoot, err)
		return
	}
	if err == nil {
		for _, entry := range exampleEntries {
			childPath := filepath.Join(examplesRoot, entry.Name())
			node, skip, err := buildTree(childPath)
			if err != nil {
				fmt.Printf("Error building tree for %s: %v\n", childPath, err)
				continue
			}
			if skip {
				continue
			}
			roots = append(roots, node)
		}
	}

	data, err := json.MarshalIndent(roots, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	outputPath := "web/src/data/bal-tree.json"
	err = os.MkdirAll(filepath.Dir(outputPath), 0o755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	err = os.WriteFile(outputPath, data, 0o644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("Successfully generated %s\n", outputPath)
}

