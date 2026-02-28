import ballerina/io;

// Shows the delivery schedule of pending orders
type Order record {|
    string id;
    string customer;
    int amount;
    int areaCode;
|};

public function main() {
    Order order1 = {id: "ORD001", customer: "Alice", amount: 150, areaCode: 1};
    Order order2 = {id: "ORD002", customer: "Bob", amount: 75, areaCode: 2};
    Order order3 = {id: "ORD003", customer: "Carol", amount: 200, areaCode: 1};

    Order[] orders = [order1, order2, order3];
    map<int> deliveryFee = {"Downtown": 5, "Suburb": 10};

    io:println("=== Delivery Schedule ===");
    foreach int code in 1 ... 2 {
        if code == 1 {
            io:println("\nDowntown (fee: $", deliveryFee["Downtown"], ")");
        }
        if code == 2 {
            io:println("\nSuburb (fee: $", deliveryFee["Suburb"], ")");
        }
        foreach Order ord in orders {
            if ord.areaCode == code {
                io:println("  ", ord.id, " | ", ord.customer, " | $", ord.amount);
            }
        }
    }
}
