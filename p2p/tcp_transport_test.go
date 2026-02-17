package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":3000"
	opts := TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddr, listenAddr)

	// server
	assert.Nil(t, tr.ListenAndAccept())
}
