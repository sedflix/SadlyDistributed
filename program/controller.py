import threading
import math

prime = 100000007

bucket_size = 1000

is_prime = True

tasks_available = set()

tasks_pending = set()

i = 0
while i*bucket_size < math.sqrt(prime):
    tasks_available.add((i*bucket_size, (i+1)*bucket_size))
    i += 1

def getTask():
    global is_prime
    if is_prime:
        if tasks_available:
            task = tasks_available.pop()
            tasks_pending.add(task)
            return "%d %d" % task
        else:
            return "no-task-available"
    else:
        return "no-task-available"


def respondToInput(): # runs in own thread
    global is_prime
    infile = open("infifo", "r")
    outfile = open("outfifo", "w")
    for line in infile:
        line = line[:-1] # remove newline
        print(is_prime)
        if line == "get-task":
            task = getTask()
            outfile.write(task + "\n")
            outfile.flush()
        elif line[0:8] == "result: ":
            result = line[8:]
            if result == "true":
                is_prime = False

inputThread = threading.Thread(target = respondToInput)
inputThread.daemon = True
inputThread.start()
input()
