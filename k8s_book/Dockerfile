FROM golang:1.16-alpine

RUN apk add --no-cache bash

WORKDIR /app
COPY go.mod ./

RUN go mod download

COPY *.go ./
RUN go build -o /k8s_example

CMD ["/k8s_example"]