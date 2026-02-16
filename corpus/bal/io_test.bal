import ballerina/io;

public function main() {
    string|error s = io:fileReadString("test.txt");
}
