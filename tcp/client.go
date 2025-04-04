package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		reader := bufio.NewReader(conn)
		for {
			response, err := reader.ReadString('\n')
			if err != nil {
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
			log.Println("Server disconnected.")
			return
		case <-stop:
			log.Println("Interrupt received. Closing connection and exiting.")
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
