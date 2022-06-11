FROM golang:1.15-alpine3.14 AS builder
WORKDIR /go/src/app
COPY . .
RUN go build -o main server.go
CMD ["./main"]