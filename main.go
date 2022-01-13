package main

import (
	"fmt"
	"os"

	"github.com/nickjameswebb/weather-alert/pkg/consumer"
	"github.com/nickjameswebb/weather-alert/pkg/producer"
)

func PrintUsage() {
	fmt.Println("Usage: weather-alert [producer|consumer]")
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		PrintUsage()
		os.Exit(1)
	}

	if args[0] == "producer" {
		producer.StartProducer()
	} else if args[0] == "consumer" {
		consumer.StartConsumer()
	} else {
		PrintUsage()
		os.Exit(1)
	}
}
