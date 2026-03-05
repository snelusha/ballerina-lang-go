import ballerina/io;

public function main() {
    string content = io:fileReadString("shello.txt");
    io:println(content);
}
