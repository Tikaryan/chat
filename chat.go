package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"
)

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go read(rw)
	go write(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}
func write(rw *bufio.ReadWriter) {
	read := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		str, err := read.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		rw.WriteString(fmt.Sprintf("%s\n", str))
		rw.Flush()
	}
}

func read(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')
		if str == "" {
			return
		}
		if str != "\n" {
			color.Green(str)
		}
	}
}
func main() {
	servPort := flag.Int("s", 8080, "Server port number")
	flag.Parse()

	r := rand.Reader
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	servMultiAdr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *servPort))
	if err != nil {
		panic(err)
	}

	hostNode, err := libp2p.New(context.Background(), libp2p.ListenAddrs(servMultiAdr), libp2p.Identity(prvKey))

	var port string
	for _, la := range hostNode.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}
	fmt.Printf("Open new bash go to client/chat folder then run -> go build chat.go\n")
	fmt.Printf("Now Run ./chat -d //ip4/127.0.0.1/tcp/%v/p2p/%s' on client node console.\n", port, hostNode.ID())
	fmt.Printf("\nWaiting for incoming connection\n\n")

	hostNode.SetStreamHandler("/chat/1.0.0", handleStream)

	<-make(chan struct{})
}
