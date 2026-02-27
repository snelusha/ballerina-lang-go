package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FileNode struct {
	Kind     string     `json:"kind"`
	Name     string     `json:"name"`
	Language string     `json:"language,omitempty"`
	Content  string     `json:"content,omitempty"`
	Children []FileNode `json:"children"`
}

const corpusBalBaseDir = "corpus/bal"

var skipTestsMap = makeSkipTestsMap([]string{
	"subset5/05-record/1-v.bal",
	"subset5/05-record/cyclic-v.bal",
	"subset5/05-record/field-access-1-v.bal",
	"subset5/05-record/field-access-2-v.bal",
	"subset5/05-record/field-access-3-e.bal",
	"subset5/05-record/field-access-4-e.bal",
	"subset5/05-record/inclusion-v.bal",
	"subset5/05-record/inclusion-override-v.bal",
	"subset5/05-record/inclusion-dup-override-v.bal",
	"subset5/05-record/inclusion-rest-v.bal",
	"subset5/06-float/01-e.bal",
	"subset5/06-float/03-e.bal",
	"subset5/06-float/04-v.bal",
	"subset5/06-float/06-v.bal",
	"subset5/06-float/11-e.bal",
	"subset5/06-float/13-e.bal",
	"subset5/06-float/15-e.bal",
	"subset5/06-float/17-e.bal",
	"subset5/06-float/float-value.bal",
})

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
		if isFileSkipped(rootPath) {
			return FileNode{}, true, nil
		}
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

func isFileSkipped(filePath string) bool {
	relPath, err := filepath.Rel(corpusBalBaseDir, filePath)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)
	return skipTestsMap[relPath]
}

func makeSkipTestsMap(paths []string) map[string]bool {
	m := make(map[string]bool, len(paths))
	for _, path := range paths {
		m[filepath.ToSlash(path)] = true
	}
	return m
}

func hasFiles(node FileNode) bool {
	if node.Kind == "file" {
		return true
	}
	for _, child := range node.Children {
		if hasFiles(child) {
			return true
		}
	}
	return false
}

func main() {
	corpusDirs := []string{
		"corpus/bal/subset1",
		"corpus/bal/subset2",
		"corpus/bal/subset3",
		"corpus/bal/subset4",
	}
	var roots []FileNode

	// Prepend hardcoded fibonacci folder
	roots = append(roots, FileNode{
		Kind: "dir",
		Name: "fibonacci",
		Children: []FileNode{
			{
				Kind:     "file",
				Name:     "main.bal",
				Language: "ballerina",
				Content: `import ballerina/io;

public function main() {
    int n = 10;
    int i = 0;
    while (i < n) {
        io:println("F(", i, ") = ", fibonacci(i));
        i += 1;
    }
}

function fibonacci(int n) returns int {
    if (n <= 1) {
        return n;
    }
    int prev = 0;
    int curr = 1;
    int i = 2;
    while (i <= n) {
        int next = prev + curr;
        prev = curr;
        curr = next;
        i += 1;
    }
    return curr;
}`,
			},
			{
				Kind:     "file",
				Name:     "Ballerina.toml",
				Language: "toml",
				Content: `[package]
org = "wso2"
name = "ballerina"
version = "0.1.0"`,
			},
		},
	})

	for _, dir := range corpusDirs {
		// We want the contents of these dirs to be at the top level, or grouped by subset
		// The user said "add everything @[corpus/bal/subset1] @[corpus/bal/subset2]"
		// I'll group them under their subset names for clarity, or just merge them?
		// Let's group them by subset name as the user suggested.

		node, _, err := buildTree(dir)
		if err != nil {
			fmt.Printf("Error building tree for %s: %v\n", dir, err)
			continue
		}
		roots = append(roots, node)
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
