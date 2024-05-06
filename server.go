package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/shamaton/msgpack/v2"
)

type (
	Payload struct {
		String string
	}
	CommandPayload struct {
		M   byte
		K   string
		V   string
		Ttl string
	}
	ResponsePayload struct {
		C string
		V string
	}
)

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
		start := time.Now()
		if err != nil {
			fmt.Println("Could not set up tcp socket:", err)
			return
		}
		cmd, err := marshallCommand(buffer, n)
		if err != nil {
			fmt.Println("Failed to marshallCommand", err)
			return
		}
		executeCommand(cmd, conn, cache)
		elapsedTime := time.Since(start)
		fmt.Printf("handleConnection took:  %s \n", elapsedTime)
	}
}

func executeCommand(cmd CommandPayload, conn net.Conn, cache *sync.Map) {
	if cmd.M == 0x2 {
		fmt.Println("Got get command")
		val, _ := cache.Load(cmd.K)
		sVal, ok := val.(string)
		fmt.Println(sVal)
		if !ok {
			return
		}
		writeResponse(conn, "1", sVal)
		// get
	} else if cmd.M == 0x1 {
		fmt.Println("Got set command")
		cache.Store(cmd.K, cmd.V)
		writeResponse(conn, "1", "")
		// set
	} else if cmd.M == 0x3 {
		fmt.Println("Get clear command")
		//clear
	}
}

func writeResponse(conn net.Conn, code string, val string) (int, error) {
	v, err := msgpack.Marshal(ResponsePayload{C: code, V: val})
	if err != nil {
		return 1, err
	}
	return conn.Write(v)
}

func marshallCommand(buffer []byte, n int) (cmd CommandPayload, err error) {
	resp := CommandPayload{}
	msgpack.Unmarshal(buffer[:n], &resp)
	fmt.Println("cmd", resp.M, resp.V, resp.K)
	if resp.M == 0x2 {
		// get
		return resp, nil
	} else if resp.M == 0x1 {
		// set
		return resp, nil
	} else if resp.M == 0x3 {
		//clear
		return resp, nil
	}
	return CommandPayload{}, errors.New("Not a valid command")
}
