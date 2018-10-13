package program

import "sync"

type Program struct {
}

type Programs struct {
	rw       sync.RWMutex
	programs []Program
}

type Parameters struct {
	start int
	end   int
}

type Result struct {
	Id     string `json:"id"`
	Result string `json:"result"`
}
