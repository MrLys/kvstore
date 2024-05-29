package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	kvt "github.com/mrlys/kvstore-types"
	"github.com/shamaton/msgpack/v2"
	"ljos.app/msgpack-tcp/utils"
)

type Client struct {
	ClientId     string
	ClientSecret string
}
type ClientStore struct {
	Clients []Client
}
type ProvisioningServer struct {
	clients sync.Map
}

func (p *ProvisioningServer) start() {
	listener, err := net.Listen("tcp", "localhost:9697")
	if err != nil {
		fmt.Println("Could not set up tcp socket:", err)
		return
	}

	defer listener.Close()
	p.clients = sync.Map{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go p.handleConnection(conn)
	}

}

func (p *ProvisioningServer) handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	payload, err := utils.ReadFromConnection(buffer, conn)
	if err != nil {
		return
	}
	err = p.authenticate(payload)

}

func (p *ProvisioningServer) readProvisionedClients() {
	data, err := os.ReadFile("store.json")
	if err != nil {
		return
	}

	var store ClientStore
	err = json.Unmarshal(data, &store)
	if err != nil {
		fmt.Println("Error unmarshalling json store")
		return
	}
	for i := 0; i < len(store.Clients); i++ {
		client := store.Clients[i]
		p.clients.Store(client.ClientId, client.ClientSecret)
	}
}

func (p *ProvisioningServer) storeClients() {
	clientBytes, err := json.Marshal(&p.clients)
	f, err := os.Create("store.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	f.Write(clientBytes)

}
func (p *ProvisioningServer) authenticate(payload []byte) error {
	var auth kvt.AuthProvisionPayload
	err := msgpack.Unmarshal(payload, &auth)
	if err != nil {
		return err
	}
	val, ok := p.clients.Load(auth.I)
	if ok && val == auth.S {
		return nil
	}
	return errors.New("Invalid key and or secret")
}

func (p *ProvisioningServer) Authenticate(key string) error {
	res := strings.Split(key, ":")
	if len(res) != 2 {
		return errors.New("Invalid key")
	}
	val, ok := p.clients.Load(res[0])
	if ok && val == res[1] {
		fmt.Println("Authorized")
		return nil
	}
	return errors.New("Invalid key")
}
