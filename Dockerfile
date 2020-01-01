FROM golang:latest
MAINTAINER lixiangyun linimbus@126.com

WORKDIR /gopath/
ENV GOPATH=/gopath/
ENV GOOS=linux
ENV CGO_ENABLED=0


RUN go get -u -v github.com/lixiangyun/benchmark
WORKDIR /gopath/src/github.com/lixiangyun/benchmark
RUN go build .


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /bin
COPY --from=0 /gopath/src/github.com/lixiangyun/benchmark/benchmark ./benchmark

RUN chmod +x *

EXPOSE 8080

CMD ["benchmark","-port","8080"]
