package node

import (
	"fmt"
	"github.com/rsms/gotalk"
	"sync"
)

type Node struct {
	Sock              *gotalk.Sock
	ResourceAvailable Resource
	ResourceUsed      Resource
	IsNew             bool
	//Job               job.Job // job currently assigned // TODO: make it an array
	rw sync.RWMutex
}

type Nodes struct {
	rw    sync.RWMutex
	Nodes []Node
}

type Resource struct {
	// Resources available in Node
	// Can be set by a user
	Cores int8  `json:"cores"` // number of corers
	Mem   int16 `json:"mem"`   // Memory in MB
}

// adds a sock to a Nodes
func (nodes Nodes) AddNode(sock *gotalk.Sock) *Node {
	node := Node{}
	node.Sock = sock
	node.IsNew = true
	node.ResourceAvailable = Resource{0, 0}
	node.ResourceUsed = Resource{0, 0}

	fmt.Println(sock.Addr())
	// TODO: // Delete handler
	//sock.CloseHandler = func(s *gotalk.Sock, _ int) {
	//	nodes.rw.Lock()
	//	defer nodes.rw.Unlock()
	//
	//}

	// add the node to the Nodes struct
	nodes.rw.Lock()
	defer nodes.rw.Unlock()
	nodes.Nodes = append(nodes.Nodes, node)
	return &node
}

func (node Node) UpdateResourceUsed(r Resource) {
	node.rw.Lock()
	defer node.rw.Unlock()
	node.ResourceUsed = r
	fmt.Printf("used resource: %+v\n", r)
}

func (node Node) UpdateResourceAvailable(r Resource) {
	node.rw.Lock()
	defer node.rw.Unlock()
	node.ResourceAvailable = r
	fmt.Printf("available resource: %+v\n", r)
}
