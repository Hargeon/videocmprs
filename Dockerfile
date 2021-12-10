# syntax=docker/dockerfile:1

FROM golang:1.16

COPY . /go/src/app

WORKDIR /go/src/app

RUN go build -o videocmprs cmd/videocmprs/main.go

EXPOSE $PORT

CMD ["./videocmprs"]
