package program

type Program struct {
	Id string
}

type Result struct {
	JobId      string `json:"job_id"`
	ProgramId  string `json:"program_id"`
	Parameters string `json:"parameters"`
	Result     string `json:"result"`
}
