import ballerina/io;

# This is a comprehensive test of markdown documentation features.
# 
# This function demonstrates basic documentation lines, parameter docs,
# return docs, inline code, code blocks, and references.
#
# Here's an inline code reference: `getValue`
# And a reference to a type: type `int`
# And a reference to a function: function `main`
#
# Here's a code block:
# ```ballerina
# int x = 10;
# ```
#
# And a double backtick code block:
# ``code content``
#
# + value - The input `int` value to process
# + return - Returns the processed `int` value
function getValue(int value) returns int {
    return value;
}

# This function has deprecation documentation.
# # Deprecated
# This function is deprecated. Use `getValue` instead.
#
# + oldValue - The old value parameter
# + return - Returns the old value
function getOldValue(int oldValue) returns int {
    return oldValue;
}

# Function with code block at the start of documentation.
# ```ballerina
# function example() {
#     // code example
# }
# ```
#
# This is regular documentation after the code block.
#
# + name - The name parameter
# + age - The age parameter  
# + return - Returns a greeting string
function greet(string name, int age) returns string {
    return "Hello";
}

# Function with multiple reference types.
# References: type `string`, service `MyService`, variable `myVar`,
# annotation `MyAnnotation`, module `mymodule`, function `myFunc`,
# parameter `myParam`, const `MY_CONST`
#
# + data - The data to process
# + return - Returns processed data
function processData(string data) returns string {
    return data;
}

# Function with inline code in documentation.
# Use `io:println()` to print output.
# The `value` parameter should be a positive `int`.
#
# + value - A positive `int` value
# + return - Returns the squared value
function square(int value) returns int {
    return value * value;
}

public function main() {
    io:println("Hello, Ballerina!");
    io:println("Value: ", getValue(10));
    io:println("Old Value: ", getOldValue(5));
    io:println("Greeting: ", greet("Alice", 30));
    io:println("Processed: ", processData("test"));
    io:println("Square: ", square(5));
}
