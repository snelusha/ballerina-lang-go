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
// software distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"os"

	"ballerina-lang-go/bir"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input.bir> <output.bir>")
		fmt.Println("Example: go run main.go example.bir example_output.bir")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	fmt.Println("=== BIR Load and Write Example ===")
	fmt.Println()

	// Step 1: Open and load the BIR file
	fmt.Printf("Loading BIR file from: %s\n", inputPath)

	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	birPackage, err := bir.LoadBIRPackageFromReader(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading BIR package: %v\n", err)
		os.Exit(1)
	}

	// Print some information about the loaded package
	fmt.Printf("Successfully loaded BIR package\n")
	if birPackage.PackageID != nil {
		pkgName := "<unknown>"
		if birPackage.PackageID.Name != nil {
			pkgName = birPackage.PackageID.Name.Value()
		}
		fmt.Printf("  Package: %s\n", pkgName)
	}
	fmt.Printf("  Functions: %d\n", len(birPackage.Functions))
	fmt.Printf("  Constants: %d\n", len(birPackage.Constants))
	fmt.Printf("  TypeDefs: %d\n", len(birPackage.TypeDefs))
	fmt.Printf("  GlobalVars: %d\n", len(birPackage.GlobalVars))
	fmt.Printf("  ImportModules: %d\n", len(birPackage.ImportModules))

	// Step 2: Create binary writer and serialize
	fmt.Printf("\nSerializing BIR package...\n")

	// Create writer (typeEnv is currently unused, can pass nil)
	writer := bir.NewBIRBinaryWriter(birPackage, nil)
	binaryData, err := writer.Serialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error serializing BIR package: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Serialized to %d bytes\n", len(binaryData))

	// Step 3: Write to output file
	fmt.Printf("Writing to: %s\n", outputPath)

	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	written, err := outputFile.Write(binaryData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	if written != len(binaryData) {
		fmt.Fprintf(os.Stderr, "Error: incomplete write: wrote %d of %d bytes\n", written, len(binaryData))
		os.Exit(1)
	}

	fmt.Printf("Successfully wrote %d bytes\n", written)
	fmt.Println()
	fmt.Println("Successfully completed round-trip: load -> write")
}
