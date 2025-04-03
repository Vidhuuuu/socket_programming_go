package main

import (
	"flag"
	"fmt"
	tcp "github.com/Vidhuuuu/socket_programming_go/tcp"
)

const addr = "127.0.0.1:5000"

func main() {
	mode := flag.String("mode", "", "Start as '[s]erver' or '[c]lient'")
	flag.Parse()

	if *mode == "" {
		fmt.Println("Usage: go run main.go -mode=[s|c]")
		return
	}

	switch *mode {
	case "s":
		tcp.StartServer(addr)
	case "c":
		tcp.StartClient(addr)
	default:
		fmt.Printf("Unknown mode: %v\n", *mode)
		fmt.Println("Usage: go run main.go -mode=[s|c]")
	}
}
