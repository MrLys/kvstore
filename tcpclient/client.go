package main

import (
	"fmt"
	uuid "github.com/google/uuid"
	kvt "github.com/mrlys/kvstore-types"
	"github.com/shamaton/msgpack/v2"
	"net"
)

func main() {

	go spamServer()
}
func spamServer() {
	conn, err := net.Dial("tcp", "localhost:9696")
	if err != nil {
		fmt.Println("Error connecing to server:", err)
		return
	}
	defer conn.Close()
	buffer := make([]byte, 1024)
	clientId := "e36c06a4-0cff-44a8-af81-24947fde2c5aa"
	clientSecret := "q+6FHhAWInTXstGU5R6VWTt9aq+dtgbpy7h+fyRefpK6suzdkWCxXZ8G1+suDLOp1ngYU8VUh5bN+htvYIcR6u4+vdh7gKOEaVM4BvyhfnUfzqxZYDuLrA6xZUnKaPAq4dsJH/gdM5BNeYEMwtGfKhFykmZKuEx8Y367TBcmqQU="
	v := &kvt.CommandPayload{M: 0x4, K: "", V: "", Ttl: "0", I: toId(clientId, clientSecret)}
	packedMsg, err := msgpack.Marshal(v)
	if err != nil {
		fmt.Println("Error during Marshal", err)
		return
	}
	_, err = conn.Write(packedMsg)
	if err != nil {
		fmt.Println("Error writing to socket", err)
		return
	}
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}
	resp := kvt.ResponsePayload{}
	msgpack.Unmarshal(buffer[:n], &resp)
	if resp.C != "1" {
		fmt.Println("Error during call.")
	}
	fmt.Println(resp.V)
	fmt.Println(resp.I)
	fmt.Println(resp.C)
	fmt.Println("Calling set")
	id := uuid.NewString()
	mykey := uuid.NewString()
	v = &kvt.CommandPayload{M: 0x1, K: mykey, V: id, Ttl: "0", I: toId(clientId, resp.I)}
	packedMsg, err = msgpack.Marshal(v)
	_, err = conn.Write(packedMsg)
	if err != nil {
		fmt.Println("Error writing to socket", err)
	}
	n, err = conn.Read(buffer)
	if err != nil {
		return
	}
	fmt.Println("Calling get")
	v = &kvt.CommandPayload{M: 0x2, K: mykey, V: "", Ttl: "0", I: toId(clientId, resp.I)}
	packedMsg, err = msgpack.Marshal(v)
	if err != nil {
		fmt.Println("Error during `json Marshal`", err)
		return
	}
	_, err = conn.Write(packedMsg)
	if err != nil {
		fmt.Println("Error writing to socket", err)
	}
	n, err = conn.Read(buffer)
	if err != nil {
		return
	}
	resp = kvt.ResponsePayload{}
	msgpack.Unmarshal(buffer[:n], &resp)
	if resp.C != "1" {
		fmt.Println("Error during call.")
	}
	fmt.Println(resp.V)
	fmt.Println(resp.I)
	fmt.Println(resp.C)
	if resp.V != id {
		fmt.Println("Error during get", resp.V, id)
	}
}
func toId(clientId, secret string) string {
	return clientId + ":" + secret
}
