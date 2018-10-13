package job

import (
	"fmt"
	"github.com/geekSiddharth/inout/server/program"
	"github.com/rsms/gotalk"
	"log"
	"time"
)

type Job struct {
	Id            string
	ProgramId     string
	JobParameter  program.Parameters
	JobResult     program.Result
	Sock          *gotalk.Sock
	IsScheduled   bool
	IsSent        bool
	InProgress    bool
	IsComputed    bool
	ScheduledTime time.Time // when was the job scheduled
	SentTime      time.Time // when was the job sent
	ComputedTime  time.Time // when was the result returned
}

type SendJob struct {
	Wasm      string `json:"wasm"`      //path of the wasm
	Parameter string `json:"parameter"` // input parameter
}

//type Jobs struct {
//	rw   sync.RWMutex
//	Jobs map[string]
//}

//func (jobs Jobs) GetTopJob() (string, error) {
//	jobs.rw.RLock()
//	defer jobs.rw.RUnlock()
//	for job, _id := range jobs.Jobs {
//		if !job {
//			job.scheduledTime = time.Now()
//			return _id, nil
//		}
//	}
//	return "", types.Error{}
//}
//
//func (jobs *Jobs) AddJob(job Job) {
//	jobs.rw.RLock()
//	defer jobs.rw.RUnlock()
//	jobs.jobs = append(jobs.jobs, job)
//}
//
//func (job Job) SetNode(node *node.Node) (Job) {
//	job.scheduledTime = time.Now()
//	job.isScheduled = true
//	node.IsNew = false
//
//	job.node = node
//	//job.SendIt()
//	return job
//}

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
	err := sock.Request("receive-job", &sendjob, naa)
	fmt.Println("After sending Request")
	if err != nil {
		log.Fatal("Unable to send job: ")
	} else {
		fmt.Println("Job Sen!! t")
	}

}

//func (job Job) SendIt() (Job) {
//	sendJob := jobToSendJob(job)
//	sokk := *job.node.Sock
//	naa := NAA{}
//	fmt.Println("Before calling sending it")
//	go SendingIt(sokk, sendJob, naa)
//	job.isSent = true
//	job.sentTime = time.Now()
//	return job
//}
//
//func MakeRandomJob() (Job) {
//	job := Job{}
//	return job
//}
