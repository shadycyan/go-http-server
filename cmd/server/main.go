package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		fmt.Println("Client connected:", conn.RemoteAddr())

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading request line:", err)
		return
	}
	fmt.Println("Request line:", strings.TrimSpace(requestLine))

	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		response := "HTTP/1.1 400 Bad Request\r\n\r\n"
		conn.Write([]byte(response))
		return
	}

	method, path := parts[0], parts[1]

	if method == "GET" && path == "/" {
		response := "HTTP/1.1 200 OK\r\n\r\n"
		conn.Write([]byte(response))
		return
	}

	if str, found := strings.CutPrefix(path, "/echo/"); method == "GET" && found {
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(str), str)
		conn.Write([]byte(response))
		return
	}

	headers, err := parseHeaders(reader)
	if err != nil {
		response := "HTTP/1.1 400 Bad Request\r\n\r\n"
		conn.Write([]byte(response))
		return
	}
	fmt.Println("Headers:", headers)

	body, err := parseBody(reader, headers)
	if err != nil {
		response := "HTTP/1.1 400 Bad Request\r\n\r\n"
		conn.Write([]byte(response))
		return
	}
	fmt.Println("Body:", body)

	if method == "GET" && path == "/user-agent" {
		userAgent, _ := headers["User-Agent"]
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
		conn.Write([]byte(response))
		return
	}

	if filename, found := strings.CutPrefix(path, "/files/"); found {
		dir := os.Args[1]
		filePath := filepath.Join(dir, filename)

		var response string
		switch method {
		case "POST":
			if err := os.WriteFile(filePath, []byte(body), 0644); err != nil {
				fmt.Println("Error writing file:", err)
				response = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
			} else {
				response = "HTTP/1.1 201 Created\r\n\r\n"
			}
		case "GET":
			data, err := os.ReadFile(filePath)
			if err != nil {
				response = "HTTP/1.1 404 Not Found\r\n\r\n"
			} else {
				response = fmt.Sprintf(
					"HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s",
					len(data), data,
				)
			}
		default:
			response = "HTTP/1.1 400 Bad Request\r\n\r\n"
		}

		conn.Write([]byte(response))
	}
}

func parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("Error reading headers: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		}
	}

	return headers, nil
}

func parseBody(reader *bufio.Reader, headers map[string]string) (string, error) {
	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return "", nil
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return "", fmt.Errorf("invalid Content-Length: %w", err)
	}

	bodyBytes := make([]byte, contentLength)
	_, err = io.ReadFull(reader, bodyBytes)
	if err != nil {
		return "", fmt.Errorf("error reading body: %w", err)
	}

	return string(bodyBytes), nil
}
