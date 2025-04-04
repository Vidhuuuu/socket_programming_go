package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	done := make(chan struct{})
	inputChan := make(chan string)

	go func() {
		reader := bufio.NewReader(conn)
		for {
			response, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Println("Server closed the connection.")
				} else {
					log.Println("Failed to read:", err)
				}
				close(done)
				return
			}
			fmt.Print("Echo: ", response)
		}
	}()

	go func() {
		stdinReader := bufio.NewReader(os.Stdin)
		for {
			s, err := stdinReader.ReadString('\n')
			if err != nil {
				close(inputChan)
				return
			}
			inputChan <- s
		}
	}()

	for {
		select {
		case <-done:
			return
		case msg, ok := <-inputChan:
			if !ok {
				return
			}
			msg = msg[:len(msg)-1]
			if _, err := fmt.Fprintln(conn, msg); err != nil {
				log.Println("Failed to send to server:", err)
				break
			}
		}
	}
}
