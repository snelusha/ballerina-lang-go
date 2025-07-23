// Copyright 2025 Sithija Nelusha
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

// IndexOutOfBoundsError represents an error when an index is out of bounds.
// This corresponds to Java's IndexOutOfBoundsException.
type IndexOutOfBoundsError struct {
	message string
}

// Error implements the error interface.
func (e IndexOutOfBoundsError) Error() string {
	return e.message
}

// NewIndexOutOfBoundsError creates a new IndexOutOfBoundsError with the given message.
func NewIndexOutOfBoundsError(message string) *IndexOutOfBoundsError {
	return &IndexOutOfBoundsError{
		message: message,
	}
}

// IllegalArgumentError represents an error when an argument is invalid.
// This corresponds to Java's IllegalArgumentException.
type IllegalArgumentError struct {
	message string
}

// Error implements the error interface.
func (e IllegalArgumentError) Error() string {
	return e.message
}

// NewIllegalArgumentError creates a new IllegalArgumentError with the given message.
func NewIllegalArgumentError(message string) *IllegalArgumentError {
	return &IllegalArgumentError{
		message: message,
	}
}
