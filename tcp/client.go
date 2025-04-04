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
	log.Printf("Connected to server(%v)...\n", addr)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		smsg := scanner.Text()
		if smsg == "bye" {
			fmt.Println("Disconnecting...")
			break
		}
		fmt.Fprintln(conn, smsg)

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Failed to capture response from server")
		}
		fmt.Print("Echo: ", response)
	}
}
