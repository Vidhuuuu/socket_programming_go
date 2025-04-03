package main

import (
	tcp "github.com/Vidhuuuu/socket_programming_go/tcp"
)

const addr = "127.0.0.1:5000"

func main() {
	tcp.StartServer(addr)
}
