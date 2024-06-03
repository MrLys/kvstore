package main

import (
	"testing"

	kvt "github.com/mrlys/kvstore-types"
	"github.com/shamaton/msgpack/v2"
)

func TestMarshallCommand(t *testing.T) {
	// buffer []byte, n int, state map[string]string
	for key, _ := range CommandMap {

		resp := kvt.CommandPayload{
			I:   "",
			K:   "",
			V:   "",
			M:   key,
			Ttl: "0",
		}
		payload, err := msgpack.Marshal(&resp)
		if err != nil {
			t.Fatalf("failed to marshall command")
		}
		state := map[string]string{}
		state["debug"] = "false"
		cmdPayload, err := marshallCommand(payload, state)
		if err != nil {
			t.Fatalf("failed to unmarshall command")
		}
		if cmdPayload.I != "" {
			t.Fatalf("failed to unmarshall command")
		}

	}
}

func TestMarshallCommandFail(t *testing.T) {
	// buffer []byte, n int, state map[string]string
	resp := kvt.CommandPayload{
		I:   "",
		K:   "",
		V:   "",
		M:   0x5,
		Ttl: "0",
	}
	payload, err := msgpack.Marshal(&resp)
	if err != nil {
		t.Fatalf("failed to marshall command")
	}
	state := map[string]string{}
	state["debug"] = "false"
	_, err = marshallCommand(payload, state)
	if err == nil {
		t.Fatalf("failed to unmarshall command")
	}
}
