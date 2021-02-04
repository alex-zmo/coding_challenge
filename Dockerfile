FROM golang:latest

RUN apt-get update
RUN apt-get install npm  -y
RUN npm install
RUN go get github.com/go-sql-driver/mysql
RUN go get github.com/DATA-DOG/go-sqlmock
RUN go get github.com/stretchr/testify/assert
RUN go get golang.org/x/crypto/bcrypt
WORKDIR /go/src/github.com/gmo-personal/coding_challenge/
ENTRYPOINT ["go", "run", "main.go"]



