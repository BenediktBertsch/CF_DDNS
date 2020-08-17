FROM golang
WORKDIR /app
ADD . /app
RUN go build &&\
    go test -v &&\
    ls -la

ENTRYPOINT [ "/app/cf_ddns" ]