FROM golang
WORKDIR /app
ADD . /app
RUN go build &&\
    go test -v 

ENTRYPOINT [ "/app/cf_ddns" ]
    