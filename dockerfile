FROM golang:latest 

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 

RUN go mod init WikiBot
RUN go mod tidy
RUN go build -o WikiBot . 

CMD ["/app/WikiBot"]