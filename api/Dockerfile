FROM golang:1.14

RUN apt-get update -qqy

WORKDIR /go/src/github.com/chrispaynes/vorChall/api

ENV GO111MODULE=on

EXPOSE 3000 3001

COPY . .

RUN go build -o /go/bin/api -mod=readonly

USER nobody

ENTRYPOINT ["/go/bin/api"]
