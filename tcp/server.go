package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var activeConns int

func StartServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	log.Printf("Listening on %s...\n", addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					log.Println("Shutting down server...")
					log.Println("Active connections on shutdown:", activeConns)
					return
				}
				log.Println("Failed to accept connection:", err)
				continue
			}
			activeConns += 1
			go func (c net.Conn) {
				defer func() {
					log.Println("Done with", conn.RemoteAddr())
					activeConns -= 1
					c.Close()
				}()
				handleConnection(conn)
			}(conn)
		}
	}()

	<-stop
	listener.Close()
}

func handleConnection(conn net.Conn) {
	log.Println("New connection from" , conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmsg := scanner.Text()
		fmt.Printf("from [%s]: %s\n", conn.RemoteAddr(), cmsg)
		fmt.Fprintln(conn, cmsg)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v\n", err)
	}
}
