package testutils

import (
	"context"
	"encoding/hex"
	"math/rand"
	"testing"

	dopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"

	datastore "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	keystore "github.com/ipfs/go-ipfs-keystore"
	"github.com/ipfs/go-ipns"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	host "github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peerstore"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	pnet "github.com/libp2p/go-libp2p-pnet"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/multiformats/go-multiaddr"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

var (
	// EncodedPK is a hex encoded key
	// to be reused across tests
	EncodedPK = "0801124018c93db89bc9614d463003dab59eb9f8028b27835d4b42abe0b707770cbfc6bd9873de48ab48d753e6be17bc50e821e09f50959da17e45448074fdecccf3e7c0"
)

// NewLibp2pHostAndDHT is used to create a new libp2p host
// and an unbootstrapped dht
func NewLibp2pHostAndDHT(
	ctx context.Context,
	t *testing.T,
	logger *zap.Logger,
	ds datastore.Batching,
	ps peerstore.Peerstore,
	pk crypto.PrivKey,
	addrs []multiaddr.Multiaddr,
	secret []byte) (host.Host, *dht.IpfsDHT) {

	var opts []libp2p.Option
	if secret != nil && len(secret) > 0 {
		var key [32]byte
		copy(key[:], secret)
		prot, err := pnet.NewV1ProtectorFromBytes(&key)
		if err != nil {
			t.Fatal(err)
		}
		opts = append(opts, libp2p.PrivateNetwork(prot))
	}
	opts = append(opts,
		libp2p.Identity(pk),
		libp2p.ListenAddrs(addrs...),
		libp2p.Peerstore(ps),
		libp2p.DefaultMuxers,
		libp2p.DefaultTransports,
		libp2p.DefaultSecurity)
	h, err := libp2p.New(ctx, opts...)
	if err != nil {
		t.Fatal(err)
	}

	idht, err := dht.New(ctx, h,
		dopts.Validator(record.NamespacedValidator{
			"pk":   record.PublicKeyValidator{},
			"ipns": ipns.Validator{KeyBook: ps},
		}),
	)
	if err != nil {
		t.Fatal(err)	
	}
	rHost := routedhost.Wrap(h, idht)
	return rHost, idht
}

// NewPrivateKey is used to create a new private key
// for testing purposes
func NewPrivateKey(t *testing.T) crypto.PrivKey {
	pkBytes, err := hex.DecodeString(EncodedPK)
	if err != nil {
		t.Fatal(err)
	}
	pk, err := crypto.UnmarshalPrivateKey(pkBytes)
	if err != nil {
		t.Fatal(err)
	}
	return pk
}

// NewSecret is used to generate a
// secret used to secure private libp2p connections
func NewSecret(t *testing.T) []byte {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		t.Fatal(err)
	}
	return data
}

// NewPeerstore is ued to generate an in-memory peerstore
func NewPeerstore(t *testing.T) peerstore.Peerstore {
	return pstoremem.NewPeerstore()
}

// NewDatastore is used to create a new in memory datastore
func NewDatastore(t *testing.T) datastore.Batching {
	return dssync.MutexWrap(datastore.NewMapDatastore())
}

// NewMultiaddr is used to create a new multiaddress
func NewMultiaddr(t *testing.T) multiaddr.Multiaddr {
	addr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")
	if err != nil {
		t.Fatal(err)
	}
	return addr
}

// NewLogger is used to return a test zap logger
func NewLogger(t *testing.T) *zap.SugaredLogger {
	return zaptest.NewLogger(t).Sugar()
}

// NewKeystore is used to return a new in memory keystore
func NewKeystore(t *testing.T) keystore.Keystore {
	return keystore.NewMemKeystore()
}
