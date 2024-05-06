package main

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/shamaton/msgpack/v2"
)

type CommandPayload struct {
	m   string
	k   string
	v   string
	ttl string
}
type ResponsePayload struct {
	c string
	v string
}

func main() {
	listener, err := net.Listen("tcp", "localhost:9696")
	if err != nil {
		fmt.Println("Could not set up tcp socket:", err)
		return
	}
	defer listener.Close()

	cache := sync.Map{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, &cache)
	}

}
func handleConnection(conn net.Conn, cache *sync.Map) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Could not set up tcp socket:", err)
			return
		}
		fmt.Println(fmt.Sprintf("Received data: %b", buffer[0]))
		fmt.Println(fmt.Sprintf("Received data: %b", buffer[1]))
		fmt.Println(fmt.Sprintf("Received data: %b", buffer[2]))
		fmt.Println(fmt.Sprintf("Buffer contains %b, %d", buffer, n))
		cmd, err := marshallCommand(buffer, n)
		//if err != nil {
		//	fmt.Println("Failed to marshallCommand", err)
		//	return
		//}
		executeCommand(cmd, conn, cache)
	}

}

func executeCommand(cmd CommandPayload, conn net.Conn, cache *sync.Map) {
	if cmd.m == "1" {
		fmt.Println("Got get command")
		val, _ := cache.Load(cmd.k)
		sVal, ok := val.(string)
		fmt.Println(sVal)
		if !ok {
			return
		}
		writeResponse(conn, "1", sVal)
		// get
	} else if cmd.m == "2" {
		fmt.Println("Got set command")
		cache.Store(cmd.k, cmd.v)
		writeResponse(conn, "1", "")
		// set
	} else if cmd.m == "3" {
		fmt.Println("Get clear command")
		//clear
	}
}

func writeResponse(conn net.Conn, code string, val string) (int, error) {
	v, err := msgpack.Marshal(ResponsePayload{c: code, v: val})
	if err != nil {
		return 1, err
	}
	return conn.Write(v)
}

func marshallCommand(buffer []byte, n int) (cmd CommandPayload, err error) {
	resp := CommandPayload{}
	msgpack.Unmarshal(buffer[:n], &resp)
	fmt.Println(fmt.Sprintf("Command contains %s", resp.m))
	fmt.Println(fmt.Sprintf("Command contains %s", resp.v))
	fmt.Println(fmt.Sprintf("Command contains %s", resp.k))
	fmt.Println(fmt.Sprintf("Command contains %s", resp.ttl))
	if resp.m == "1" {
		// get
		return resp, nil
	} else if resp.m == "2" {
		// set
		return resp, nil
	} else if resp.m == "3" {
		//clear
		return resp, nil
	}
	return CommandPayload{}, errors.New("Not a valid command")
}
