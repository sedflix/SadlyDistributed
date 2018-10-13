package server

import (
	"github.com/geekSiddharth/inout/server/job"
	"github.com/geekSiddharth/inout/server/node"
	"github.com/geekSiddharth/inout/server/program"
	"github.com/rsms/gotalk"
)

type Server struct {
	Programs program.Programs
	Nodes    node.Nodes
	Jobs     job.Jobs
	Socks    map[*gotalk.Sock]*node.Node
}

func (server *Server) Init() {
	server.Programs = program.Programs{}
	server.Nodes = node.Nodes{}
	server.Jobs = job.Jobs{}
	server.Socks = make(map[*gotalk.Sock]*node.Node)
}

//func (server Server) add(sock gotalk.Sock) {
//	sock.se
//}
