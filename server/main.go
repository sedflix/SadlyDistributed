package main

import (
	"fmt"
	"github.com/geekSiddharth/inout/server/job"
	"github.com/geekSiddharth/inout/server/node"
	"github.com/geekSiddharth/inout/server/program"
	"github.com/geekSiddharth/inout/server/server"
	"github.com/hpcloud/tail"
	"github.com/rsms/gotalk"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	serverThis server.Server
	// TODO: number of result that can be queued
	resultChan = make(chan program.Result, 10)
)

type SendJob struct {
	JobId      string `json:"job_id"`
	ProgramId  string `json:"program_id"`
	Wasm       string `json:"wasm"`       //path of the wasm
	Parameters string `json:"parameters"` // input parameter
}

type JobReceiveResponse struct {
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
				for _, job_ := range serverThis.Jobs {
					if job_.IsScheduled == false {
						//schedule it
						sendJob := SendJob{
							JobId:      job_.Id,
							ProgramId:  job_.ProgramId,
							Parameters: job_.Parameters,
							Wasm:       "/primes/bigprimes.wasm",
						}
						go func() {
							jobRecieveResponse := &JobReceiveResponse{}
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
	// why not write result to the job

	// delete job
	delete(serverThis.Jobs, r.JobId)
	serverThis.RWJobs.Unlock()

	// write to file
	resultChan <- r

	fmt.Printf("Job Completed: %+v\n", r)
	return "Okay", nil
}

func resultChanFeeder() {
	for r := range resultChan {
		f, err := os.OpenFile(r.ProgramId+"_output", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Printf("UNABLE to write: %s \n", r.ProgramId)
		}

		defer f.Close()
		fmt.Fprintf(f, "%s\t%s\n", r.Parameters, r.Result)
		//fmt.Println(f, "get-task")
	}
}

func programJobCreator(programID string) {

	t, err := tail.TailFile(programID+"_input", tail.Config{
		Follow: true,
		ReOpen: true,
	})

	if err != nil {
		fmt.Println("Unable to open the file for program id: %d", programID)
		log.Fatal(err)
		return
	}
	for line := range t.Lines {
		text := strings.Trim(line.Text, "\n\t\r")

		switch text {
		case "<end>":
			// End creating news jobs
			return
		case "<end_all>":
			// jon through jobs and delete all the jobs of this kind
			// TODO: KILLL MEEEEEE
			return
		default:
			newJob := job.Job{
				ProgramId:    programID,
				Parameters:   text,
				CreationTime: time.Now(),
			}

			serverThis.RWJobs.Lock()
			_id := strconv.FormatInt(int64(len(serverThis.Jobs)+1), 10)
			newJob.Id = _id
			serverThis.Jobs[_id] = newJob
			serverThis.RWJobs.Unlock()
		}
	}

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

	// starting scheduler in background
	go scheduler()

	// writes result to a file
	go resultChanFeeder()

	// TODO: Make it dynamic -> Program
	go programJobCreator("1")

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
