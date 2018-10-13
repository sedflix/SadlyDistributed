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
    
##  Browser Side functions
### Handler 
- get-total-resources: returns a json with the following fields
    - cores
    - mem
