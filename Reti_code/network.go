package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"
)

type Network struct {
	Address string

	PSpeerNodes []string
	CSpeerNodes []string

	PSBroadcastPeers []*net.TCPConn
	CSBroadcastPeers []*net.TCPConn

	PSBroadcastMessages chan Message
	CSBroadcastMessages chan Message
	IncomingMessages    chan Message
}

func SetupNetwork(node Node) *Network {
	n := new(Network)
	n.Address = node.Addr
	n.PSpeerNodes = node.PSPeerAddr
	n.CSpeerNodes = node.CSPeerAddr
	n.IncomingMessages = make(chan Message, 100)
	n.CSBroadcastMessages, n.PSBroadcastMessages = make(chan Message, 100), make(chan Message, 100)
	go n.Listen()

	// Try to connect other nodes, try three times max.
	tries := 3
	for _, peer := range node.PSPeerAddr {
		var conn *net.TCPConn
		for i := 1; i <= tries; i++ {
			tcpAddr, err := net.ResolveTCPAddr("tcp", peer)
			if err != nil {
				log.Printf("Error resolving address %s: %s\n", peer, err)
				break
			}

			conn, err = net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				log.Printf("Error connecting to node %s (try %d/%d): %s\n", peer, i, tries, err)
				if i == tries {
					break
				}
				// wait 5 second, then try again
				time.Sleep(5 * time.Second)
				continue
			}
			// connect successful
			log.Printf("Successfully connected to node %s\n", peer)
			n.PSBroadcastPeers = append(n.PSBroadcastPeers, conn)
			break
		}
	}
	for _, peer := range node.CSPeerAddr {
		var conn *net.TCPConn
		for i := 1; i <= tries; i++ {
			tcpAddr, err := net.ResolveTCPAddr("tcp", peer)
			if err != nil {
				log.Printf("Error resolving address %s: %s\n", peer, err)
				break
			}

			conn, err = net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				log.Printf("Error connecting to node %s (try %d/%d): %s\n", peer, i, tries, err)
				if i == tries {
					break
				}
				// wait 5 second, then try again
				time.Sleep(5 * time.Second)
				continue
			}
			// connect successful
			log.Printf("Successfully connected to node %s\n", peer)
			n.CSBroadcastPeers = append(n.CSBroadcastPeers, conn)
			break
		}
	}
	return n
}

func (n *Network) Listen() {
	tcpaddr, _ := net.ResolveTCPAddr("", n.Address)
	listener, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		log.Fatalf("Error starting listener: %s\n", err)
	}

	n.Address = listener.Addr().String()
	log.Printf("Listening on %s\n", n.Address)

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				log.Printf("Error accepting connection: %s\n", err)
				continue
			}

			log.Printf("Received connection from %s\n", conn.RemoteAddr().String())

			//n.BroadcastPeers = append(n.BroadcastPeers, conn)

			go n.HandleMessage(conn)
		}
	}()
}

func (n *Network) HandleMessage(conn *net.TCPConn) {
	for {

		buf := make([]byte, 1024)
		num, err := conn.Read(buf)

		if err != nil {
			log.Printf("Error reading from TCP connection: %s\n", err)
			return
		}
		var msg Message
		decoder := gob.NewDecoder(bytes.NewReader(buf[:num]))
		err = decoder.Decode(&msg)
		if err != nil {
			log.Printf("Error decoding message from %s: %s\n", conn.RemoteAddr().String(), err)
			break
		}
		n.IncomingMessages <- msg
	}
}

func (n *Network) Run() {
	for {
		select {
		case message := <-n.PSBroadcastMessages:
			n.BroadcastMessage(message, 0)

		case message := <-n.CSBroadcastMessages:
			n.BroadcastMessage(message, 1)
		}

	}
}

func (n *Network) BroadcastMessage(message Message, psorcs int) {
	var temp bytes.Buffer
	error := gob.NewEncoder(&temp).Encode(message)
	if error != nil {
		log.Printf("Error encoding message: %s\n", error)
	}
	if psorcs == 0 {
		for _, conn := range n.PSBroadcastPeers {
			//log.Printf("Broadcasting to process shard node %s\n", conn.RemoteAddr().String())
			_, err := conn.Write(temp.Bytes())
			if err != nil {
				log.Printf("Error broadcasting to process shard node%s: %s\n", conn.RemoteAddr().String(), err)
			}
		}
	}
	if psorcs == 1 {
		for _, conn := range n.CSBroadcastPeers {
			//log.Printf("Broadcasting to control shard node %s\n", conn.RemoteAddr().String())
			_, err := conn.Write(temp.Bytes())
			if err != nil {
				log.Printf("Error broadcasting to control shard node%s: %s\n", conn.RemoteAddr().String(), err)
			}
		}
	}
}
