package main

import "ballerina-lang-go/bir"

func main() {
	err := bir.Parse("")
	if err != nil {
		panic(err)
	}
}
