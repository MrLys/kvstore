package main

import (
	"fmt"
	"net"

	"github.com/shamaton/msgpack/v2"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:9696")
	if err != nil {
		fmt.Println("Could not set up tcp socket:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}

}
func handleConnection(conn net.Conn) {
	type Struct struct {
		String string
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Could not set up tcp socket:", err)
			return
		}
		resp := Struct{}
		msgpack.Unmarshal(buffer[:n], &resp)
		fmt.Printf("Received %s:\n", resp.String)
	}

}
