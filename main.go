package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const addr = "127.0.0.1:5000"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	defer listener.Close()

	log.Printf("Listening on %s...\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("New connection from %s\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("Received: %s\n", text)
		fmt.Fprintf(conn, "You said: %s\n", text)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v\n", err)
	}
}
