package main

import "github.com/codecrafters-io/redis-starter-go/app/transmission"

func main() {
	transmission.MakeServer(9092)
}
