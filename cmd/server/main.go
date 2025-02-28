package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting server on port 42069...")

	l, err := net.Listen("tcp", "0.0.0.0:42069")
	if err != nil {
		fmt.Println("Failed to bind to port 42069")
		os.Exit(1)
	}

	_, err = l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
}
