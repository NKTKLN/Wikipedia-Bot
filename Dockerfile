FROM golang:latest 

WORKDIR /usr/src/app
COPY ./ ./

RUN go mod init bot
RUN go mod tidy
RUN go build main.go 
