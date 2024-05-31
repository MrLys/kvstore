package internal

import (
	kvt "github.com/mrlys/kvstore-types"
	"github.com/shamaton/msgpack/v2"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	var p ProvisioningServer
	p.clients.Store("1024", "aoeuhtns")
	auth := kvt.AuthProvisionPayload{I: "1024", S: "aoeuhtns"}
	packedMsg, err := msgpack.Marshal(auth)
	if err != nil || p.authenticate(packedMsg) != nil {
		t.Fatalf("error!")
	}
}
func TestAuthenticateWrongId(t *testing.T) {
	var p ProvisioningServer
	p.clients.Store("1024", "aoeuhtns")
	auth := kvt.AuthProvisionPayload{I: "1023", S: "aoeuhtns"}
	packedMsg, err := msgpack.Marshal(auth)
	if err != nil || p.authenticate(packedMsg) == nil {
		t.Fatalf("error!")
	}
}
func TestAuthenticateWrongSecret(t *testing.T) {
	var p ProvisioningServer
	p.clients.Store("1024", "aoeuhtns")
	auth := kvt.AuthProvisionPayload{I: "1024", S: "Not the correct secret"}
	packedMsg, err := msgpack.Marshal(auth)
	if err != nil || p.authenticate(packedMsg) == nil {
		t.Fatalf("error!")
	}
}
