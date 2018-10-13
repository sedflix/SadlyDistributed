package server

import (
	"github.com/geekSiddharth/inout/server/job"
	"github.com/geekSiddharth/inout/server/node"
	"github.com/geekSiddharth/inout/server/program"
	"github.com/rsms/gotalk"
	"sync"
)

type Server struct {
	Programs  program.Programs
	Jobs      job.Jobs
	Socks     map[*gotalk.Sock]node.Node
	RWSocks sync.RWMutex
}

func (server *Server) Init() {
	server.Programs = program.Programs{}
	server.Jobs = job.Jobs{}
	server.Socks = make(map[*gotalk.Sock]node.Node)
}

//func (server Server) add(sock gotalk.Sock) {
//	sock.se
//}
