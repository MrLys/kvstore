package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

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
		v := CommandPayload{m: "uouauoauaua", k: "test key", v: "test value, what else?", ttl: "0"}
		packedMsg, err := msgpack.Marshal(v)
		t := CommandPayload{}
		err = msgpack.Unmarshal(packedMsg, &t)
		if err != nil {
			fmt.Println("Error during Marshal", err)
		}
		fmt.Println(fmt.Sprintf("msg : %s", t.k))
		fmt.Println("packedMsg len", len(packedMsg))
		for i := 0; i < len(packedMsg); i++ {
			fmt.Println(fmt.Sprintf("msg : %b", packedMsg[i]))
		}
		fmt.Println("Error during Marshal", err)
		if err != nil {
			fmt.Println("Error during Marshal", err)
		}
		fmt.Println(fmt.Sprintf("packedMsg: %d", packedMsg))
		_, err = conn.Write([]byte{0x5, 0x0, 0x4})
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
		if resp.c != "1" {
			fmt.Println("Error during call.")
		}
		fmt.Println(resp.v)
	}

}
