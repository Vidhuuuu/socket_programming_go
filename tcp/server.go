package tcp

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func StartServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	defer listener.Close()

	log.Printf("Listening on %s...\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		go func (c net.Conn) {
			defer func() {
				log.Println("Done with", conn.RemoteAddr())
				c.Close()
			}()
			handleConnection(conn)
		}(conn)
	}
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
