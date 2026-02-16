import ballerina/io;

public function main() returns error? {
    string path = "test.txt";
    string|error s = io:fileReadString(path);
    io:println(s);
}
