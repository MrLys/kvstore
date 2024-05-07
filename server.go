package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	kvt "github.com/mrlys/kvstore-types"
	"github.com/shamaton/msgpack/v2"
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
		return
	}
}

func executeCommand(cmd kvt.CommandPayload, conn net.Conn, cache *sync.Map) {
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
		cache.Delete(cmd.K)
		fmt.Println("Get clear command")
		writeResponse(conn, "1", "")
		//clear
	}
}

func writeResponse(conn net.Conn, code string, val string) (int, error) {
	v, err := msgpack.Marshal(kvt.ResponsePayload{C: code, V: val})
	if err != nil {
		return 1, err
	}
	return conn.Write(v)
}

func marshallCommand(buffer []byte, n int) (cmd kvt.CommandPayload, err error) {
	resp := kvt.CommandPayload{}
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
	return kvt.CommandPayload{}, errors.New("Not a valid command")
}
