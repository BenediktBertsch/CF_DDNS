FROM golang:alpine

RUN mkdir /app
ADD . /app
WORKDIR /app/src
RUN go build

ENTRYPOINT [ "/app/src/ddns" ]