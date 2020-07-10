FROM golang:1.13 as builder

ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on

ADD ./ /go/src/github.com/A-Ethan/golang-demo
RUN go get github.com/go-sql-driver/mysql

WORKDIR /go/src/github.com/A-Ethan/golang-demo

RUN CGO_ENABLED=0 GOOS=linux go build main.go -a -installsuffix cgo .

EXPOSE 80

CMD ["./main"]
