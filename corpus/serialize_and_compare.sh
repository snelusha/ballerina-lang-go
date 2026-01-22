#!/bin/bash

# Script to serialize BIR files from subset1 and compare with original files
# Usage: ./serialize_and_compare.sh
#
# Note: This script only processes actual binary BIR files.
# Git LFS pointer files are automatically skipped.

set -e

# Get the script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BIR_EXAMPLES_DIR="$PROJECT_ROOT/bir/examples"
INPUT_DIR="$SCRIPT_DIR/bir/subset1"
OUTPUT_DIR="$SCRIPT_DIR/tmp-bir"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
TOTAL=0
PASSED=0
FAILED=0
ERRORS=0
SKIPPED=0

# Function to check if file is a Git LFS pointer
is_lfs_pointer() {
    local file="$1"
    # Check if file starts with "version https://git-lfs.github.com/spec/v1"
    head -1 "$file" 2>/dev/null | grep -q "^version https://git-lfs.github.com/spec/v1" && return 0 || return 1
}

# Function to check if file is a valid binary BIR file
is_valid_bir() {
    local file="$1"
    # Check if file starts with BIR magic bytes: 0xba 0x10 0xc0 0xde
    local magic=$(head -c 4 "$file" 2>/dev/null | od -An -tx1 | tr -d ' \n')
    [ "$magic" = "ba10c0de" ] && return 0 || return 1
}

echo "=== BIR Serialization and Comparison Script ==="
echo ""
echo "Input directory: $INPUT_DIR"
echo "Output directory: $OUTPUT_DIR"
echo ""

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Change to project root for go run to work
cd "$PROJECT_ROOT"

# Find all .bir files in subset1 and process them
# Use process substitution to avoid subshell issues with counters
while IFS= read -r input_file; do
    TOTAL=$((TOTAL + 1))
    
    # Get relative path from subset1
    rel_path="${input_file#$INPUT_DIR/}"
    
    # Check if file is a Git LFS pointer
    if is_lfs_pointer "$input_file"; then
        echo -e "${YELLOW}SKIPPED${NC}: $rel_path (Git LFS pointer - need to pull actual file)"
        SKIPPED=$((SKIPPED + 1))
        continue
    fi
    
    # Check if file is a valid binary BIR file
    if ! is_valid_bir "$input_file"; then
        echo -e "${YELLOW}SKIPPED${NC}: $rel_path (not a valid binary BIR file)"
        SKIPPED=$((SKIPPED + 1))
        continue
    fi
    
    # Create output path maintaining directory structure
    output_file="$OUTPUT_DIR/$rel_path"
    output_dir="$(dirname "$output_file")"
    
    # Create output directory if needed
    mkdir -p "$output_dir"
    
    # Run serialization
    echo -n "Processing: $rel_path ... "
    
    if go run "$BIR_EXAMPLES_DIR/main.go" "$input_file" "$output_file" > /dev/null 2>&1; then
        # Compare files
        if cmp -s "$input_file" "$output_file"; then
            echo -e "${GREEN}PASSED${NC}"
            PASSED=$((PASSED + 1))
        else
            echo -e "${RED}FAILED${NC} (files differ)"
            FAILED=$((FAILED + 1))
        fi
    else
        echo -e "${RED}ERROR${NC} (serialization failed)"
        ERRORS=$((ERRORS + 1))
    fi
done < <(find "$INPUT_DIR" -type f -name "*.bir" | sort)

# Print summary
echo ""
echo "=== Summary ==="
echo "Total files: $TOTAL"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo -e "${RED}Errors: $ERRORS${NC}"
if [ $SKIPPED -gt 0 ]; then
    echo -e "${YELLOW}Skipped: $SKIPPED${NC} (LFS pointers or invalid files)"
fi

# Exit with error if any failures
if [ $FAILED -gt 0 ] || [ $ERRORS -gt 0 ]; then
    exit 1
else
    exit 0
fi
