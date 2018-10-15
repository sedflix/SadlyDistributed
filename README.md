# SadlyDistributed

 Making volunteer computing easy af.  
 Zero Setup for volunteers. Easy for developers.  
 
 We compile **your distributed code in almost any language** to **WebAssembly** and then all of it is **executed in a cluster of browsers**.
 Anyone who just opens our website in their browser will share his or her computing power with us.

## Motivation

We are aware of the amazing computing power of our laptops and phones. What if the scientist around the world would be able to use it when its lying idl(which is most of the time)? What if everyone with a computing device, a browser, and a decent internet connection would be able to contribute to such amazing research happening around the world? Today's world is a sucker for computing resources. ML, AR, VR, Blockchain and tons of other cutting edge tech require tons of computing resources.

There are several nice volunteer computing platform like BONIC, GridCoin, etc. But all of them suffer from a common problem: complex and heavy installation process. Even I, as a techie, hesitate to go through all that hassle. You know setting up anything sucks. Therefore, we came up with SadlyDistributed. 


## Installation 
We <3 Docker. Hence, we have provided a Dockerfile that takes care of all the dependencies of our server code.

#### Building Docker file
```
    docker build -t sadlyDistributed .
```

#### Running our Server(with Prime Number Example)
```
    docker run -it -p 8899:8899 sadlyDistributed
```
 
## Mini-Tutorial
**How can you modify your distributed code for our architecture?**  

We have an example program for finding if a number(it can be as large as you can think of) is prime or not.   
Code: [/client/programs/1](/client/programs/1)  

In general distributed computer, same code is replicated on multiple machines and each machine execute the code over a different range of values(or parameters). 
The output from multiple machines is combined to produce the final output. Your job, as developer, is to give us your pre-existing distributed code, 
in almost any language. The code should take `input from stdin` and gives `output to stdout`. 
After that, you need to define how the range of values should be *divided* and how the results from them should be *combined* to produce the final answer.
After you have defined this logic, we scale your code to all the browsers available to us.  

You code interface with our architecture by reading and writing from files. We have two files:
#### `input`
You specify the parameters that your code will take in. Each parameters is present in a new line. Our prime number uses the following format:
```
1 110000
110000 9990000
9990000 99900000
99900000 99900000
<number 1><space><number 2>

```
Note that our distributed code, `bigprimes.go`, is written to understand these range. In other words, 
`go run bigprimes.go 110000 9990000` will tell us weather our hardcoded number is divisible by anything between 110000 and 9990000. 
Our architecture read this file continuously(as new input arrives or so called tail reading) and distributes the parameters over free nodes.

#### `output`
Our architecture writes the stdout(or output) received from various machine to this file.
Now it is upto you to utilize the info in this file for combining the result. Our prime number used the following format:
```
1 110000    false
99900000 99900000   false
110000 9990000  false
<parameter><tab><output>
```
The parameter is same as specified in the input file. And the output is what your code wrote to stdout when given those parameters.
Our architecture updates this file as it receives output from different browsers.

In our prime number example, we use two file in two different language to show how **our approach is language-agnostic**.
We have used GoLang for the code that needs to be distributed. And we have used python for generating inputs and combining the output. Neat, right?
