package node

import (
	"fmt"
	"github.com/rsms/gotalk"
	"sync"
	"time"
)

type Node struct {
	Sock              *gotalk.Sock
	ResourceAvailable Resource
	ResourceUsed      Resource
	IsNew             bool
	TTL               time.Duration
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
func GetNode(sock *gotalk.Sock) Node {
	node := Node{}
	node.Sock = sock
	node.IsNew = true
	node.ResourceAvailable = Resource{0, 0}
	node.ResourceUsed = Resource{0, 0}
	node.TTL = 5 * time.Second
	return node
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

func (nodes *Nodes) GetAvailableNodes() []*Node {
	nodes.rw.RLock()
	defer nodes.rw.RUnlock()
	availableNodes := make([]*Node, 0)
	fmt.Printf("Number of total nodes %d ", len(nodes.Nodes))
	for _, node := range nodes.Nodes {
		//if node.ResourceAvailable.Cores-node.ResourceUsed.Cores > 0 {
		if node.IsNew {
			availableNodes = append(availableNodes, &node)
		}
	}
	return availableNodes

}
