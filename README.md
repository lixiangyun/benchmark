# benchmark
- Tcp/HTTP tls benchmark tools 

## ![](https://travis-ci.com/lixiangyun/benchmark.svg?branch=master)

## build
```
go build .
```

## usage
```
Usage of benchmark:
  -address string
        service/client address. (default "127.0.0.1")
  -body int
        transport body length (KB). (default 1)
  -ca string
        tls using ca file for verify other cert.
  -cert string
        tls using cert file.
  -debug
        display debug infomation.
  -head string
        set request head. as[key1=value1,key2=value2]
  -help
        usage help.
  -key string
        tls using key file.
  -limit int
        limit times per second to send. as[0,1000]
  -par int
        the parallel numbers to connect. (default 1)
  -port int
        service/client port. (default 8001)
  -protocal string
        which protocol to using (http/http2/tcp). (default "http")
  -role string
        the tools role (server/client). (default "server")
  -runtime int
        total run time (second). (default 60)
  -stat int
        stat display intval (second). (default 5)
  -tls
        enable/disable tls mode.
  -url string
        set request url. as[/abc/123] (default "/")
```

## HTTPS benchmark test
### https server mode
```
benchmark.exe -role server -tls -protocal http -address 127.0.0.1 -port 8080
```

### https client mode
```
benchmark.exe -role client -tls -protocal http -par 10 -url /test1 -head key1=value1 -address 127.0.0.1 -port 8080
```

### tcp under tls mode
```
benchmark.exe -role server -tls -protocal tcp -address 127.0.0.1 -port 8080
```

### tcp under tls mode
```
benchmark.exe -role client -tls -protocal tcp -address 127.0.0.1 -port 8080
```
