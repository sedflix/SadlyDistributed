package job

import (
	"github.com/geekSiddharth/hackinout/server/node"
	"github.com/geekSiddharth/hackinout/server/program"
	"sync"
	"time"
)

type Job struct {
	program       program.Program
	jobParameter  program.Parameters
	jobResult     program.Result
	node          node.Node
	isScheduled   bool
	isSent        bool
	inProgress    bool
	isComputed    bool
	scheduledTime time.Time // when was the job scheduled
	sentTime      time.Time // when was the job sent
	computedTime  time.Time // when was the result returned
}

type Jobs struct {
	rw   sync.RWMutex
	jobs []Job
}
