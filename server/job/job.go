package job

import (
	"fmt"
	"github.com/geekSiddharth/inout/server/node"
	"github.com/geekSiddharth/inout/server/program"
	"github.com/rsms/gotalk"
	"go/types"
	"log"
	"sync"
	"time"
)

type Job struct {
	program       program.Program
	jobParameter  program.Parameters
	jobResult     program.Result
	node          *node.Node
	isScheduled   bool
	isSent        bool
	inProgress    bool
	isComputed    bool
	scheduledTime time.Time // when was the job scheduled
	sentTime      time.Time // when was the job sent
	computedTime  time.Time // when was the result returned
}

type SendJob struct {
	Wasm      string `json:"wasm"`      //path of the wasm
	Parameter string `json:"parameter"` // input parameter
}
type Jobs struct {
	rw   sync.RWMutex
	jobs []Job
}

func (jobs Jobs) GetTopJob() (Job, error) {
	jobs.rw.RLock()
	defer jobs.rw.RUnlock()
	for _, job := range jobs.jobs {
		if !job.isScheduled {
			job.scheduledTime = time.Now()
			return job, nil
		}
	}
	return Job{}, types.Error{}
}

func (jobs *Jobs) AddJob(job Job) {
	jobs.rw.RLock()
	defer jobs.rw.RUnlock()
	jobs.jobs = append(jobs.jobs, job)
}

func (job Job) SetNode(node *node.Node) (Job) {
	job.scheduledTime = time.Now()
	job.isScheduled = true
	node.IsNew = false

	job.node = node
	//job.SendIt()
	return job
}

func jobToSendJob(job Job) (SendJob) {
	sendJob := SendJob{}
	sendJob.Parameter = "so wow parameters"
	sendJob.Wasm = "so wow wasm"
	return sendJob
}

type NAA struct {
	Res string `json:"res"`
}

func SendingIt(sock gotalk.Sock, sendjob SendJob, naa NAA) {
	err := sock.Request("receive-job", &sendjob, nil)
	if err != nil {
		log.Fatal("Unable to send job: ")
	} else {
		fmt.Println("Job Sen!! t")
	}

}

func (job Job) SendIt() (Job) {
	sendJob := jobToSendJob(job)
	sokk := *job.node.Sock
	naa := NAA{}
	go SendingIt(sokk, sendJob, naa)
	job.isSent = true
	job.sentTime = time.Now()
	return job
}

func MakeRandomJob() (Job) {
	job := Job{}
	return job
}
