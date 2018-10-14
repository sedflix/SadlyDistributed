from flask import Flask
import threading
import math

app = Flask(__name__)

prime = 100000007

bucket_size = 1000

is_prime = True

tasks_available = set()

tasks_pending = set()

i = 0
while i*bucket_size < math.sqrt(prime):
    tasks_available.add((i*bucket_size, (i+1)*bucket_size))
    i += 1

@app.route("/get-task")
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
        return "<end>" # end all?


@app.route("/result/<yay>")
def respondToInput(yay): # runs in own thread
    global is_prime
    result = yay.split("\t")[1]
    if result == "true":
        is_prime = False

app.run(host="localhost", port="5005")
