import ballerina/io;

public function main() {
    int a = 10;
    int b = a / 0;
    io:println("Result: ", b);
}
