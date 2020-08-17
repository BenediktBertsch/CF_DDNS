FROM golang
WORKDIR /app
ADD . /app
RUN go build &&\
    go test -v &&\
    /app/cf_ddns

ENTRYPOINT [ "/bin/bash" ]