package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	uuid "github.com/google/uuid"
	kvt "github.com/mrlys/kvstore-types"
	"github.com/shamaton/msgpack/v2"
	"ljos.app/msgpack-tcp/utils"
)

var clients = sync.Map{}
var authenticatedClients = sync.Map{}
var CommandMap = map[byte]string{
	0x1: "set",
	0x2: "get",
	0x3: "clear",
	0x4: "auth",
}

func main() {
	clients.Store("e36c06a4-0cff-44a8-af81-24947fde2c5a", "q+6FHhAWInTXstGU5R6VWTt9aq+dtgbpy7h+fyRefpK6suzdkWCxXZ8G1+suDLOp1ngYU8VUh5bN+htvYIcR6u4+vdh7gKOEaVM4BvyhfnUfzqxZYDuLrA6xZUnKaPAq4dsJH/gdM5BNeYEMwtGfKhFykmZKuEx8Y367TBcmqQU=")
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
	state := map[string]string{}
	// read first command. It should be auth request
	payload, err := utils.ReadFromConnection(buffer, conn)
	requestId := uuid.NewString()
	state["requestId"] = requestId
	err = handleAuthenticationRequest(payload, state)
	if err != nil {
		log("Failed to handleRequest", state, err)
		return
	}
	writeResponse(conn, "1", "OK")

	for {
		payload, err := utils.ReadFromConnection(buffer, conn)
		requestId := uuid.NewString()
		state["requestId"] = requestId
		if err != nil {
			log("Could not set up tcp socket:", state, err)
			return
		}
		err = handleRequest(payload, conn, cache, state)
		if err != nil {
			log("Failed to handleRequest", state, err)
			return
		}
	}
}

func handleAuthenticationRequest(payload []byte, state map[string]string) error {
	start := time.Now()
	elapsedTime := time.Since(start)
	log("handleConnection took:  %s", state, elapsedTime)
	cmd, err := marshallCommand(payload, state)
	if err != nil {
		log("Failed to marshallCommand", state, err)
		return err
	}
	if cmd.M != 0x4 {
		msg, err := json.Marshal(cmd)
		if err != nil {
			return errors.New("Invalid command, authentication request expected")
		}
		log("Invalid command, expected 0x4, got %s", state, string(msg))
		return errors.New("Invalid command, authentication request expected")
	}
	return nil
}

func handleRequest(payload []byte, conn net.Conn, cache *sync.Map, state map[string]string) error {
	start := time.Now()
	cmd, err := marshallCommand(payload, state)
	if err != nil {
		log("Failed to marshallCommand", state, err)
		return err
	}
	err = executeCommand(cmd, conn, cache, state)
	if err != nil {
		log("Failed to executeCommand", state, err)
		return err
	}
	elapsedTime := time.Since(start)
	log("handleConnection took:  %s", state, elapsedTime)
	return nil
}

func log(msg string, state map[string]string, args ...interface{}) {
	if state["debug"] != "false" {
		fmt.Printf("(%s) "+msg+"\n", state["requestId"], args)
	}
}

func authenticate(addr string, key string, state map[string]string) error {
	res := strings.Split(key, ":")
	if len(res) != 2 {
		log("Invalid key(1)", state, res, key)
		return errors.New("Invalid key")
	}
	val, ok := clients.Load(res[0])
	if ok && val == res[1] {
		fmt.Println("Authorized")
		authMap, ok := authenticatedClients.Load(addr)
		if !ok {
			authMap = map[string]string{}
		}
		authMap.(map[string]string)[res[0]] = uuid.New().String()
		authenticatedClients.Store(addr, authMap)
		return nil
	}
	log("Invalid key(2)", state, res, key)
	return errors.New("Invalid key")
}

func authenticateWithId(addr string, key string, state map[string]string) (string, error) {
	res := strings.Split(key, ":")
	if len(res) != 2 {
		log("Invalid key(3)", state, res, key)
		return "", errors.New("Invalid key")
	}
	authMap, ok := authenticatedClients.Load(addr)
	if !ok {
		return "", errors.New("Not authenticated")
	}
	_, ok = authMap.(map[string]string)[res[0]]
	if !ok {
		return "", errors.New("Not authenticated")
	}
	newId := uuid.NewString()
	authMap.(map[string]string)[res[0]] = newId
	authenticatedClients.Store(addr, authMap)
	return newId, nil
}

func executeCommand(cmd kvt.CommandPayload, conn net.Conn, cache *sync.Map, state map[string]string) error {
	if cmd.M == 0x4 {
		log("Got %s command", state, CommandMap[cmd.M])
		err := authenticate(conn.RemoteAddr().String(), cmd.I, state)
		if err != nil {
			writeResponse(conn, "0", "Unauthorized")
			return err
		}
		writeResponse(conn, "1", "OK")
	} else if cmd.M == 0x2 {
		log("Got %s command", state, CommandMap[cmd.M])
		val, _ := cache.Load(cmd.K)
		sVal, ok := val.(string)
		if !ok {
			return nil
		}
		writeResponse(conn, "1", sVal)
		// get
	} else if cmd.M == 0x1 {
		log("Got %s command", state, CommandMap[cmd.M])
		cache.Store(cmd.K, cmd.V)
		writeResponse(conn, "1", "")
		// set
	} else if cmd.M == 0x3 {
		log("Got %s command", state, CommandMap[cmd.M])
		cache.Delete(cmd.K)
		writeResponse(conn, "1", "")
		//clear
	}
	return nil
}

func writeResponse(conn net.Conn, code string, val string) (int, error) {
	v, err := msgpack.Marshal(kvt.ResponsePayload{C: code, V: val, I: ""})
	if err != nil {
		return 1, err
	}
	return conn.Write(v)
}

func marshallCommand(payload []byte, state map[string]string) (cmd kvt.CommandPayload, err error) {
	resp := kvt.CommandPayload{}
	msgpack.Unmarshal(payload, &resp)
	log("cmd", state, resp.M, resp.V, resp.K)
	_, ok := CommandMap[resp.M]
	if !ok {
		return kvt.CommandPayload{}, errors.New("Not a valid command")
	}
	return resp, nil
}
