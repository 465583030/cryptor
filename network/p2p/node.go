package p2p

import (
	"fmt"
	"net"

	"github.com/thee-engineer/cryptor/network"
)

// Node ..
type Node struct {
	NodeConfig // Static configuration generated at node creation

	conn *net.UDPConn // UDP listening connection

	addr *net.UDPAddr  // Node UDP address
	quit chan struct{} // Stops node from running when it receives
	errc chan error    // Channell for transmiting errors

	addp chan *Peer    // Add peer request channel
	pops chan peerFunc // Peer count and peer list operations
	popd chan struct{} // Peer operation done

	peers map[string]*Peer  // Memory map with key/value peer pairs
	token map[string][]byte // List of tokens used in requests
}

// Function for peer list and peer count
type peerFunc func(map[string]*Peer)

// NewNode ...
func NewNode(ip string, port int, quit chan struct{}) *Node {
	return &Node{
		addr:  network.IPPToUDP(ip, port),
		quit:  quit,
		addp:  make(chan *Peer),
		errc:  make(chan error),
		pops:  make(chan peerFunc),
		popd:  make(chan struct{}),
		peers: make(map[string]*Peer),
	}
}

// Start ...
func (n *Node) Start() {

	go n.listen()

	for {
		select {
		case err := <-n.errc:
			fmt.Println("err:", err) // DEBUG
		case <-n.quit:
			return
		case peer := <-n.addp:
			n.peers[peer.addr.String()] = peer
		case operation := <-n.pops:
			operation(n.peers)
			n.popd <- struct{}{}
		}
	}
}

// Stop ...
func (n *Node) Stop() {
	close(n.quit)
}

// AddPeer ...
func (n *Node) AddPeer(peer *Peer) {
	select {
	case <-n.quit:
	case n.addp <- peer:
	}
}

// Peers ...
func (n *Node) Peers() []*Peer {
	var peerList []*Peer

	select {
	case n.pops <- func(peers map[string]*Peer) {
		for _, p := range peers {
			peerList = append(peerList, p)
		}
	}:
		<-n.popd
	case <-n.quit:

	}

	return peerList
}

// PeerCount ...
func (n *Node) PeerCount() int {
	var count int

	select {
	case n.pops <- func(peerList map[string]*Peer) { count = len(peerList) }:
		<-n.popd
	case <-n.quit:
	}

	return count
}

func (n *Node) listen() {
	fmt.Println("listening")
	conn, err := net.ListenUDP("udp", n.addr)
	if err != nil {
		n.errc <- err
		return
	}
	defer conn.Close()

	var buffer [1024]byte

	for {
		r, addr, err := conn.ReadFromUDP(buffer[:])
		if err != nil {
			n.errc <- err
			return
		}
		if r > 0 {
			// DEBUG
			fmt.Println(addr.String(), "|", r, "|", string(buffer[:r]))
			go n.dial(addr)
		}
	}
}

// WIP
func (n *Node) dial(addr *net.UDPAddr) {
	fmt.Println("dial")
	conn, err := net.DialUDP("udp", n.addr, addr)
	if err != nil {
		n.errc <- err
		return
	}
	defer conn.Close()

	conn.Write([]byte("hello world"))
}
