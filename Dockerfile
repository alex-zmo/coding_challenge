FROM golang:1.14.0

RUN apt-get update
RUN apt-get install npm -y
RUN npm install
WORKDIR /go/src/github.com/gmo-personal/coding_challenge/
ENTRYPOINT ["go", "run", "main.go"]



