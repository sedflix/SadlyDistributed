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
            return "<end>"
    else:
        return "<end>"


def writeTasks():
    outfile = open("1_input", "w")
    while len(tasks_pending) < 100 and tasks_available:
        task = getTask()
        outfile.write(task + "\n")
        outfile.flush()

def respondToInput(): # runs in own thread
    global is_prime
    infile = open("1_output", "r")
    writeTasks()
    for line in infile:
        line = line[:-1] # remove newline
        if len(line.split("\t")) == 2:
            result = line.split("\t")[1]
            if result == "true":
                is_prime = False
        writeTasks()

inputThread = threading.Thread(target = respondToInput)
inputThread.daemon = True
inputThread.start()
input("press any key to continue")
