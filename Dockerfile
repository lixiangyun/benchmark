FROM golang:latest
MAINTAINER lixiangyun linimbus@126.com

WORKDIR /gopath/
ENV GOPATH=/gopath/
ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /gopath/src/github.com/lixiangyun/
RUN git clone https://github.com/lixiangyun/benchmark.git

WORKDIR /gopath/src/github.com/lixiangyun/benchmark/tcp
RUN go build .

WORKDIR /gopath/src/github.com/lixiangyun/benchmark/httpserver
RUN go build .

WORKDIR /gopath/src/github.com/lixiangyun/benchmark/httpclient
RUN go build .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /opt/
COPY --from=0 /gopath/src/github.com/lixiangyun/benchmark/tcp/tcp ./tcp
COPY --from=0 /gopath/src/github.com/lixiangyun/benchmark/httpserver/httpserver ./httpserver
COPY --from=0 /gopath/src/github.com/lixiangyun/benchmark/httpclient/httpclient ./httpclient

RUN chmod +x *

EXPOSE 8080

ENTRYPOINT ["./httpserver","-p",":8080"]
