FROM golang:1.14

# ARG SERVICE

# ENV SERVICE "$SERVICE"
# RUN echo $SERVICE

# RUN export PATH=$PATH:$(which go)

RUN apt-get update -qqy

# WORKDIR /go/src/github.com/chrispaynes/gRPCrud
# COPY . .
# RUN ls -lahg

VOLUME /go/src/github.com/chrispaynes/vorChall

WORKDIR /go/src/github.com/chrispaynes/vorChall/api

ENV GO111MODULE=on

EXPOSE 3000 3001

COPY . .

RUN go build -o /go/bin/api

# USER nobody

ENTRYPOINT ["/go/bin/api"]
