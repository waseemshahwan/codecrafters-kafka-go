package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	fmt.Println("Listening on port 9092")

	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}

	l.Accept()
}
