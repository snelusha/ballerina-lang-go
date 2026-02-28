import ballerina/io;

// Aggregates responses from multiple API endpoints, categorizes them by status
// and flags, and processes mixed result types.
type ApiResponse record {|
    string endpoint;
    int status;
    int flags;
    any result;
|};

public function main() {
    // Response flags
    int cached = 1 << 0;
    int partial = 1 << 1;
    int paginated = 1 << 2;

    map<string> statusDesc = {"200": "OK", "206": "Partial"}; // status descriptions

    ApiResponse res1 = {endpoint: "/users", status: 200, flags: cached, result: 150};
    ApiResponse res2 = {endpoint: "/orders", status: 200, flags: cached + paginated, result: "paginated"};
    ApiResponse res3 = {endpoint: "/products", status: 206, flags: partial, result: 42};
    ApiResponse[] responses = [res1, res2, res3];
    [string, int] healthCheck = ["/health", 200];

    io:println("Health Check: ", healthCheck[0], " (", healthCheck[1], ")");
    io:println("\nAPI Responses: ", responses.length());
    io:println("------------------");

    foreach ApiResponse res in responses {
        io:println("\nEndpoint: ", res.endpoint);
        if res.status == 200 {
            io:println("  Status: ", res.status, " ", statusDesc["200"]);
        }
        if res.status == 206 {
            io:println("  Status: ", res.status, " ", statusDesc["206"]);
        }

        if res.flags >= paginated {
            io:println("  [Paginated]");
        }
        if res.flags == cached {
            io:println("  [Cached]");
        }
        if res.flags == partial {
            io:println("  [Partial]");
        }

        any result = res.result;
        if result is int {
            int count = <int>result;
            io:println("  Count: ", count);
            if count > 100 {
                io:println("  Volume: High");
            } else {
                io:println("  Volume: Normal");
            }
        }
        if result is string {
            io:println("  Info: ", result);
        }
    }
}
