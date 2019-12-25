FROM golang:latest
MAINTAINER lixiangyun linimbus@126.com

WORKDIR /gopath/
ENV GOPATH=/gopath/
ENV GOOS=linux
ENV CGO_ENABLED=0

COPY ./httpserver /gopath/src/github.com/lixiangyun/benchmark/httpserver
COPY ./httpclient /gopath/src/github.com/lixiangyun/benchmark/httpclient
COPY ./tcp /gopath/src/github.com/lixiangyun/benchmark/tcp

WORKDIR /gopath/src/github.com/lixiangyun/benchmark/tcp
RUN go build .

WORKDIR /gopath/src/github.com/lixiangyun/benchmark/httpserver
RUN go build .

WORKDIR /gopath/src/github.com/lixiangyun/benchmark/httpclient
RUN go build .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /
COPY --from=0 /gopath/src/github.com/lixiangyun/benchmark/tcp/tcp ./tcp
COPY --from=0 /gopath/src/github.com/lixiangyun/benchmark/httpserver/httpserver ./httpserver
COPY --from=0 /gopath/src/github.com/lixiangyun/benchmark/httpclient/httpclient ./httpclient

RUN chmod +x *

EXPOSE 8080

CMD ["httpserver","-p",":8080"]
