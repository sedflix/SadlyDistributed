package main

import (
	"fmt"
	"github.com/geekSiddharth/inout/server/job"
	"github.com/geekSiddharth/inout/server/node"
	"github.com/geekSiddharth/inout/server/program"
	"github.com/geekSiddharth/inout/server/server"
	"github.com/rsms/gotalk"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	serverThis server.Server
)

type SendJob struct {
	Id        string `json:"id"`
	Wasm      string `json:"wasm"`      //path of the wasm
	Parameter string `json:"parameter"` // input parameter
}

type JobRecieveResponse struct {
	IsOkay string `json:"is_okay"`
}

func scheduler() {
	for {
		serverThis.RWSocks.RLock()
		// go through all the sockets
		for sockThis, node_ := range serverThis.Socks {
			if node_.IsNew == true {
				fmt.Printf("node avail: %s \n", node_.Sock.Addr())
				// do error checking
				serverThis.RWJobs.Lock()
				for job_id, job_ := range serverThis.Jobs {
					if job_.IsScheduled == false {
						//schedule it
						sendJob := SendJob{
							Id:        job_id,
							Parameter: "parameters for " + job_id,
							Wasm:      "Wasm path for " + job_id,
						}
						go func() {
							jobRecieveResponse := &JobRecieveResponse{}
							err := sockThis.Request("receive-job", &sendJob, jobRecieveResponse)
							if err != nil {
								fmt.Println(err)
							} else {

								if strings.Compare(jobRecieveResponse.IsOkay, "Okay") == 0 {
									job_.IsScheduled = true
									job_.ScheduledTime = time.Now()
									job_.Sock = sockThis

									node_.IsNew = false
									serverThis.Socks[sockThis] = node_
								} else {
									fmt.Printf("Job Receive Response error: %s \n", jobRecieveResponse.IsOkay)
								}

							}
							return
						}()
						break
					}
				}
				serverThis.RWJobs.Unlock()

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

	// Making the socket free
	serverThis.RWSocks.Lock()
	node_ := serverThis.Socks[s]
	node_.IsNew = true
	serverThis.Socks[s] = node_
	serverThis.RWSocks.Unlock()

	// delete the job and write to a file
	serverThis.RWJobs.Lock()
	job_ := serverThis.Jobs[r.Id]
	job_.JobResult = r

	// delete job
	delete(serverThis.Jobs, r.Id)

	//write to a file
	// TODO : Delete from jobs and send it to a channel
	serverThis.RWJobs.Unlock()
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
		var i int64;
		for i = 0; i < 10; i++ {
			_id := strconv.FormatInt(i, 10)
			serverThis.Jobs[_id] = job.Job{
				Id: _id,
			}
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
