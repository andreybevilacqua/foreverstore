package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/andreybevilacqua/foreverstore/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fsOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}
	s := NewFileServer(fsOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

// 08:02:02
func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(2 * time.Second)

	go s2.Start()
	time.Sleep(2 * time.Second)

	// data := bytes.NewReader([]byte("my big data file here!"))
	// s2.Store("coolPicture.jpg", data)
	// time.Sleep(5 * time.Millisecond)

	r, err := s2.Get("coolPicture.jpg")
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
