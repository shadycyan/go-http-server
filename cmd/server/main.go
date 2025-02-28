package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:42069")
	if err != nil {
		fmt.Println("Failed to bind to port 42069")
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Server listening on port 42069")

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading request line:", err)
		return
	}
	fmt.Println("Request line:", strings.TrimSpace(requestLine))

	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		fmt.Println("Malformed request")
		return
	}

	path := parts[1]

	response := "HTTP/1.1 404 Not Found\r\n\r\n"

	if path == "/" {
		response = "HTTP/1.1 200 OK\r\n\r\n"
	}

	conn.Write([]byte(response))
}
