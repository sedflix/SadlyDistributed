package main

import (
	"fmt"
	"github.com/geekSiddharth/inout/server/job"
	"github.com/geekSiddharth/inout/server/node"
	"github.com/geekSiddharth/inout/server/program"
	"github.com/geekSiddharth/inout/server/server"
	"github.com/rsms/gotalk"
	"net/http"
	"time"
)

var (
	serverThis server.Server
)

func scheduler() {
	for {

		serverThis.RWSocks.RLock()

		// go through all the sockets
		for sockThis, node_ := range serverThis.Socks {
			if node_.IsNew == true {
				fmt.Printf("node avail: %s \n", node_.Sock.Addr())
				// do error checking
				toBeScheduledJob, err := serverThis.Jobs.GetTopJob()
				if err == nil {
					fmt.Printf("Job found\n")
					toBeScheduledJob = toBeScheduledJob.SetNode(&node_)
					node_.IsNew = false
					serverThis.Socks[sockThis] = node_
					toBeScheduledJob = toBeScheduledJob.SendIt()
				} else {
					fmt.Println("Jobs not found")
				}
			} else {
				fmt.Printf("node busy: %s \n", node_.Sock.Addr())
			}
		}
		time.Sleep(1 * time.Second)
		serverThis.RWSocks.RUnlock()
	}

}

func handleJobComplete(s *gotalk.Sock, r program.Result) (string, error) {
	fmt.Println("in job complete handler")
	serverThis.RWSocks.Lock()

	node_ := serverThis.Socks[s]
	node_.IsNew = true
	serverThis.Socks[s] = node_

	serverThis.RWSocks.Unlock()
	fmt.Printf("Job Completed: %+v\n", r)
	return "Okay", nil
}

func onAcceptConnection(sock *gotalk.Sock) {
	fmt.Println("Accepted: ", sock.Addr())

	serverThis.RWSocks.Lock()
	defer serverThis.RWSocks.Unlock()

	serverThis.Socks[sock] = node.GetNode(sock)
	// TODO: Add locks here
	sock.CloseHandler = func(s *gotalk.Sock, _ int) {
		serverThis.RWSocks.Lock()
		defer serverThis.RWSocks.Unlock()
		delete(serverThis.Socks, s)
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

	// adding lots of of jobs for test

	go func() {
		for i := 0; i < 10; i++ {
			serverThis.Jobs.AddJob(job.MakeRandomJob())
		}
	}()

	// starting scheduler in background
	go scheduler()

	// Handle Result
	gotalk.Handle("job-complete", handleJobComplete)

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
	err := http.ListenAndServe("0.0.0.0:1233", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
