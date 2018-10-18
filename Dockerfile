FROM golang:1.11

COPY . /go/src/materialized-view
WORKDIR /go/src/materialized-view

ENV GO111MODULE=on

RUN go build

EXPOSE 8080

CMD ./materialized-view