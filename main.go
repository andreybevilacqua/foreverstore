package main

import (
	"log"
	"time"

	"github.com/andreybevilacqua/foreverstore/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// todo: OnPeer: ,
	}
	fsOpts := FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         p2p.NewTCPTransport(tcpTransportOpts),
	}
	s := NewFileServer(fsOpts)
	go func() {
		time.Sleep(time.Second * 3)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
