package server

import (
	"github.com/geekSiddharth/inout/server/job"
	"github.com/geekSiddharth/inout/server/node"
	"github.com/geekSiddharth/inout/server/program"
	"github.com/rsms/gotalk"
	"sync"
)

type Server struct {
	Programs   map[string]program.Programs
	RWPrograms sync.RWMutex
	Jobs       map[string]job.Job
	RWJobs     sync.RWMutex
	Socks      map[*gotalk.Sock]node.Node
	RWSocks    sync.RWMutex
}

func (server *Server) Init() {
	server.Programs = make(map[string]program.Programs)
	server.Jobs = make(map[string]job.Job)
	server.Socks = make(map[*gotalk.Sock]node.Node)
}

//func (server Server) add(sock gotalk.Sock) {
//	sock.se
//}
