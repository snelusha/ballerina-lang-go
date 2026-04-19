import ballerina/io;

public function main() {
    string|error content = io:fileReadString("hello.txt");
    if (content is string) {
        io:println(content);
    } else {
        io:println("error while reading file");
    }
}
