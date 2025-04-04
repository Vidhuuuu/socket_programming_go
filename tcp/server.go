package tcp

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	activeConns   = make(map[net.Conn]struct{})
	activeConnsMu sync.Mutex
)

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
					return
				}
				log.Println("Failed to accept connection:", err)
				continue
			}

			activeConnsMu.Lock()
			activeConns[conn] = struct{}{}
			activeConnsMu.Unlock()

			go func (c net.Conn) {
				defer func() {
					log.Println("Done with", conn.RemoteAddr())
					activeConnsMu.Lock()
					delete(activeConns, c)
					activeConnsMu.Unlock()
					c.Close()
				}()
				handleConnection(conn)
			}(conn)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	activeConnsMu.Lock()
	for conn := range activeConns {
		conn.Close()
	}
	numClosed := len(activeConns)
	activeConns   = make(map[net.Conn]struct{})
	activeConnsMu.Unlock()

	listener.Close()
	log.Println("Active connections on shutdown:", numClosed)
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
		return
	}
}
