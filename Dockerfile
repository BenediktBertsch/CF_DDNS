FROM golang
WORKDIR /app/src
ADD . /app
RUN go build &&\
    go test -v 

ENTRYPOINT [ "/app/src/ddns" ]