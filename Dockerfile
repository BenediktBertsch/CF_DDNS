FROM golang:alpine
WORKDIR /app/src
ADD . /app
RUN apk add build-base &&\
    go build &&\
    go test -v &&\
    /app/src/ddns