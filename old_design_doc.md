# Server

## Data Structure
- Program
- Job
- Node
- Program List
- Nodes List
- Jobs Queue

## Functions 
- Program To Jobs 
- Task Scheduler (Global): Job to Node matching 
- Add Node
- Remove Node
- Receive Info From Node
    - Receives any message
    - forwards it to apt function
- Send Job To Node
    - a function that keeps on checking Job array and it sends the scheduled tasks to the scheduled node and make necessary changes to the job
- Result Aggregator
    - Runs for each program
    - 
    
###  Handlers at teh server side functions
- receive-job:  (json)
    - It should send the following: 
        - IsOkay string `json:"is_okay"`
    - It should accept:   (json)
    	- JobId      string `json:"job_id"`
    	- ProgramId  string `json:"program_id"`
    	- Wasm       string `json:"wasm"`       //path of the wasm
    	- Parameters string `json:"parameters"` // input parameter
        
### Handlers at teh server side 
- resource-available: returns a json with the following fields
    - accepts  (json)
        - cores
        - mem
    - returns 
        - Okay
- resource-used: accepts a json with the fields same as above and send "Okay"
- job-complete: called after a task is done
    - accepts (json)
        - JobId      string `json:"job_id"`
        - ProgramId  string `json:"program_id"`
        - Parameters string `json:"parameters"`
        - Result     string `json:"result"`
    - returns
        - "Okay"


### Task config
 - file_name:  `<program_id>_input` should contains the input parameter
 - each line should contain the parameter
 - the parameter is passed as it is to the receive-job  handler on the browser
### Result Config
   I will result of the input taken from `<program_id>_input` to this file
  - file_name: `<program_id>_output`. 
  - Each row is like: <parameter>\t<Result>