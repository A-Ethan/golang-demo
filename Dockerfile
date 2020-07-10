FROM golang:1.13 as builder

RUN go get github.com/A-Ethan/golang-demo && go get github.com/disintegration/imaging
RUN go get github.com/go-sql-driver/mysql

WORKDIR /go/src/github.com/A-Ethan/golang-demo

RUN CGO_ENABLED=0 GOOS=linux go build main.go -a -installsuffix cgo .

EXPOSE 80

CMD ["./main"]
