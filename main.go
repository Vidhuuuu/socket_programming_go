package main

import (
	tcpServer "github.com/Vidhuuuu/socket_programming_go/tcp/tcp_server"
)

const addr = "127.0.0.1:5000"

func main() {
	tcpServer.TCPStartServer(addr)
}
