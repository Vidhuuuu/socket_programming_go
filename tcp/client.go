package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func StartClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to connect to %v: %v\n", addr, err)
	}
	defer conn.Close()

	fmt.Printf("Connected to server(%v). Type messages here\n", addr)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "bye" {
			fmt.Println("Disconnecting...")
			break
		}
		fmt.Fprintln(conn, msg)
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Failed to capture response from server")
		}
		fmt.Printf("Server: %v", response)
	}
}
