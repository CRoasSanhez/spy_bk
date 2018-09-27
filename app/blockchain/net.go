package blockchain

/*
import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type ConnectionsQueue chan string
type NodeChannel chan *Node
type Node struct {
	*net.TCPConn
	lastSeen int
}

type Nodes map[string]*Node

type Network struct {
	Nodes
	ConnectionsQueue
	Address            string
	ConnectionCallback NodeChannel
	BroadcastQueue     chan Message
	IncomingMessages   chan Message
}

func (n Nodes) AddNode(node *Node) bool {

	key := node.TCPConn.RemoteAddr().String()

	if key != BlockCore.Network.Address && n[key] == nil {

		fmt.Println("Node connected", key)
		n[key] = node

		go HandleNode(node)

		return true
	}
	return false
}

func HandleNode(node *Node) {
	// TODO ... handle node algorithm
}

func SetupNetwork(address, port string) *Network {

	n := new(Network)

	n.BroadcastQueue, n.IncomingMessages = make(chan Message), make(chan Message)
	n.ConnectionsQueue, n.ConnectionCallback = CreateConnectionsQueue()
	n.Nodes = Nodes{}
	n.Address = address //fmt.Sprintf("%s:%s", address, port)

	return n
}

func (n *Network) Run() {

	fmt.Println("Listening in", BlockCore.Address)
	listenCb := StartListening(BlockCore.Address)

	for {
		select {
		case node := <-listenCb:
			Nodes.AddNode(node)

		case node := <-n.ConnectionCallback:
			Nodes.AddNode(node)

		case message := <-n.BroadcastQueue:
			go n.BroadcastMessage(message)
		}
	}
}

func CreateConnectionsQueue() (ConnectionsQueue, NodeChannel) {

	in := make(ConnectionsQueue)
	out := make(NodeChannel)

	go func() {

		for {
			address := <-in

			address = fmt.Sprintf("%s:%s", address, BLOCKCHAIN_PORT)

			if address != Network.Address && Nodes[address] == nil {

				go ConnectToNode(address, 5*time.Second, false, out)
			}
		}
	}()

	return in, out
}

func StartListening(address string) NodeChannel {

	cb := make(NodeChannel)
	addr, err := net.ResolveTCPAddr("tcp4", address)
	networkError(err)

	listener, err := net.ListenTCP("tcp4", addr)
	networkError(err)

	go func(l *net.TCPListener) {

		for {
			connection, err := l.AcceptTCP()
			networkError(err)

			cb <- &Node{connection, int(time.Now().Unix())}
		}

	}(listener)

	return cb
}

func ConnectToNode(dst string, timeout time.Duration, retry bool, cb NodeChannel) {

	addrDst, err := net.ResolveTCPAddr("tcp4", dst)
	networkError(err)

	var con *net.TCPConn = nil
loop:
	for {

		breakChannel := make(chan bool)
		go func() {

			con, err = net.DialTCP("tcp", nil, addrDst)

			if con != nil {

				cb <- &Node{con, int(time.Now().Unix())}
				breakChannel <- true
			}
		}()

		select {
		case <-Timeout(timeout):
			if !retry {
				break loop
			}
		case <-breakChannel:
			break loop
		}

	}
}

func (n *Network) BroadcastMessage(message Message) {

	b, _ := message.MarshalBinary()

	for k, node := range n.Nodes {
		fmt.Println("Broadcasting...", k)
		go func() {
			_, err := node.TCPConn.Write(b)
			if err != nil {
				fmt.Println("Error bcing to", node.TCPConn.RemoteAddr())
			}
		}()
	}
}

func GetIpAddress() []string {

	name, err := os.Hostname()
	if err != nil {

		return nil
	}

	addrs, err := net.LookupHost(name)
	if err != nil {

		return nil
	}

	return addrs
}

func networkError(err error) {

	if err != nil && err != io.EOF {

		log.Println("Blockchain network: ", err)
	}
}
*/
