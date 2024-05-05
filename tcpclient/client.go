package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/shamaton/msgpack/v2"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter greeting:")
	text, _, _ := reader.ReadLine()
	type Struct struct {
		String string
	}
	v := Struct{String: fmt.Sprintf("Greeting form client: %v", text)}
	packedMsg, err := msgpack.Marshal(v)
	conn, err := net.Dial("tcp", "localhost:9696")
	if err != nil {
		fmt.Println("Error connecing to server:", err)
		return
	}
	defer conn.Close()
	_, err = conn.Write([]byte(packedMsg))
	if err != nil {
		fmt.Println("Error connecing to server:", err)
		return
	}

}
