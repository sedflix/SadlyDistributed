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
	result string `json:"result"`
}
