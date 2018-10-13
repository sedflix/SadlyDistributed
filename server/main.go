package main

import (
	"fmt"
	"github.com/geekSiddharth/hackinout/server/node"
	"github.com/geekSiddharth/hackinout/server/server"
	"github.com/rsms/gotalk"
	"net/http"
)

var (
	serverThis server.Server
)

func onAcceptConnection(sock *gotalk.Sock) {

	fmt.Println("Accepted: ", sock.Addr())
	serverThis.Socks[sock] = serverThis.Nodes.AddNode(sock)
	// TODO: Add locks here
	sock.CloseHandler = func(s *gotalk.Sock, _ int) {
		delete(serverThis.Socks, sock)
		fmt.Println("Closed")
	}

	// add the node to the Nodes struct
}

func main() {
	serverThis = server.Server{}
	serverThis.Init()

	//gotalk.HandleBufferRequest("ech0", func(s *gotalk.Sock, op string, b []byte) ([]byte, error) {
	//	return b, nil
	//})

	// RESOURCE STUFFS
	gotalk.Handle("resource-available",
		func(s *gotalk.Sock, r node.Resource) (string, error) {
			serverThis.Socks[s].UpdateResourceAvailable(r)
			return "Okay", nil
		})

	gotalk.Handle("resource-used",
		func(s *gotalk.Sock, r node.Resource) (string, error) {
			serverThis.Socks[s].UpdateResourceUsed(r)
			return "Okay", nil
		})

	webSocketHandler := gotalk.WebSocketHandler()
	webSocketHandler.OnAccept = onAcceptConnection
	http.Handle("/gotalk/", webSocketHandler)
	http.Handle("/", http.FileServer(http.Dir("./client")))
	err := http.ListenAndServe("0.0.0.0:1235", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
