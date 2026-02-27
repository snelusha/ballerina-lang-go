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
	Children []FileNode `json:"children,omitempty"`
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

func buildTree(rootPath string) (FileNode, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return FileNode{}, err
	}

	node := FileNode{
		Name: filepath.Base(rootPath),
	}

	if info.IsDir() {
		node.Kind = "dir"
		entries, err := os.ReadDir(rootPath)
		if err != nil {
			return FileNode{}, err
		}

		for _, entry := range entries {
			childPath := filepath.Join(rootPath, entry.Name())
			childNode, err := buildTree(childPath)
			if err != nil {
				return FileNode{}, err
			}
			node.Children = append(node.Children, childNode)
		}
	} else {
		node.Kind = "file"
		node.Language = getLanguage(node.Name)
		content, err := os.ReadFile(rootPath)
		if err != nil {
			return FileNode{}, err
		}
		node.Content = string(content)
	}

	return node, nil
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

		node, err := buildTree(dir)
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
