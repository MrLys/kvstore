package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"

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
	reader := bufio.NewReader(os.Stdin)
	conn, err := net.Dial("tcp", "localhost:9696")
	if err != nil {
		fmt.Println("Error connecing to server:", err)
		return
	}
	defer conn.Close()
	counter := 0
	for {
		buffer := make([]byte, 1024)
		fmt.Println("Enter greeting:")
		text, _, _ := reader.ReadLine()
		fmt.Println(fmt.Sprintf("Greeting form client: %s", text))
		//var m byte
		//if counter%2 == 0 {
		//	m = 0x1
		//} else {
		//	m = 0x0
		//}
		v := &CommandPayload{M: 0x2, K: "test key", V: "test value, what elseee?", Ttl: "0"}
		msg, err := json.Marshal(v)
		if err != nil {
			fmt.Println("Error during `json Marshal`", err)
			return
		}
		fmt.Println("msg", string(msg))
		//data := Payload{String: string(msg)}
		packedMsg, err := msgpack.Marshal(v)
		fmt.Println("packedMsg len", len(packedMsg))
		fmt.Println("Error during Marshal", err)
		if err != nil {
			fmt.Println("Error during Marshal", err)
		}
		fmt.Println(fmt.Sprintf("packedMsg: %d", packedMsg))
		_, err = conn.Write(packedMsg)
		if err != nil {
			fmt.Println("Error writing to socket", err)
		}

		counter++
		if err != nil {
			fmt.Println("Error connecing to server:", err)
			return
		}
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		resp := ResponsePayload{}
		msgpack.Unmarshal(buffer[:n], &resp)
		if resp.C != "1" {
			fmt.Println("Error during call.")
		}
		fmt.Println(resp.V)
	}

}
