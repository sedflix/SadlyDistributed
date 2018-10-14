FROM golang:1.11

## Install Python
#RUN apt-get install python2.7
## Install node.js
#RUN apt-get install nodejs
## Install CMake (optional, only needed for tests and building Binaryen)
#RUN apt-get install cmake
## Install Java (optional, only needed for Closure Compiler minification)
#RUN apt-get install default-jre
#
## Get the emsdk repo
#RUN git clone https://github.com/juj/emsdk.git
#
## Enter that directory
#WORKDIR ./emsdk
#RUN git pull
## Download and install the latest SDK tools.
#RUN ./emsdk install latest
## Make the "latest" SDK "active" for the current user. (writes ~/.emscripten file)
#RUN ./emsdk activate latest
## Activate PATH and other environment variables in the current terminal
#RUN source ./emsdk_env.

RUN apt-get install python3
RUN go get github.com/hpcloud/tail
RUN go get github.com/rsms/gotalk
WORKDIR $GOPATH/src/github.com/geekSiddharth/inout/

COPY . $GOPATH/src/github.com/geekSiddharth/inout/

CMD ["make", "wasm"]
