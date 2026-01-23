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

	"ballerina-lang-go/bir"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input.bir> <output.bir>\n", os.Args[0])
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file %s: %v\n", inputPath, err)
		os.Exit(1)
	}
	defer inputFile.Close()

	// Load BIR package from binary
	fmt.Printf("Loading BIR package from %s...\n", inputPath)
	birPackage, err := bir.LoadBIRPackageFromReader(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading BIR package: %v\n", err)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error type: %T\n", err)
			fmt.Fprintf(os.Stderr, "Error string: %s\n", err.Error())
		}
		os.Exit(1)
	}

	if birPackage == nil {
		fmt.Fprintf(os.Stderr, "Error: loaded BIR package is nil\n")
		os.Exit(1)
	}

	if birPackage.PackageID == nil {
		fmt.Fprintf(os.Stderr, "Error: BIR package has nil PackageID\n")
		os.Exit(1)
	}

	pkgName := birPackage.PackageID.Name.Value()
	fmt.Printf("Successfully loaded BIR package: %s\n", pkgName)

	// Create writer and serialize
	fmt.Printf("Serializing BIR package...\n")
	writer := bir.NewBIRBinaryWriter(birPackage)
	serialized, err := writer.Serialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error serializing BIR package: %v\n", err)
		os.Exit(1)
	}

	// Write to output file
	fmt.Printf("Writing serialized BIR to %s...\n", outputPath)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file %s: %v\n", outputPath, err)
		os.Exit(1)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(serialized)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully serialized BIR package to %s (%d bytes)\n", outputPath, len(serialized))
}
