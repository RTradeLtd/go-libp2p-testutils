package testutils

import (
	"context"
	"testing"

	"github.com/multiformats/go-multiaddr"
)

func Test_NewPrivateKey(t *testing.T) {
	if pk := NewPrivateKey(t); pk == nil {
		t.Fatal("should not be nil")
	}
}

func Test_NewSecret(t *testing.T) {
	if secret := NewSecret(t); secret == nil {
		t.Fatal("should not be nil")
	}
}

func Test_NewPeerstore(t *testing.T) {
	if peerstore := NewPeerstore(t); peerstore == nil {
		t.Fatal("should not be nil")
	}
}

func Test_NewDatastore(t *testing.T) {
	if datastore := NewDatastore(t); datastore == nil {
		t.Fatal("should not be nil")
	}
}

func Test_NewMultiaddr(t *testing.T) {
	if addr := NewMultiaddr(t); addr == nil {
		t.Fatal("should not be nil")
	}
}

func Test_NewKeystore(t *testing.T) {
	if keystore := NewKeystore(t); keystore == nil {
		t.Fatal("should not be nil")
	}
}

func Test_NewLibp2pHostAndDHT(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := NewLogger(t)
	ds := NewDatastore(t)
	ps := NewPeerstore(t)
	pk := NewPrivateKey(t)
	addrs := []multiaddr.Multiaddr{NewMultiaddr(t)}
	host, dht := NewLibp2pHostAndDHT(
		ctx, t, logger.Desugar(), ds, ps, pk, addrs, NewSecret(t),
	)
	dht.Close()
	host.Close()
}
