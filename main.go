package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	STAT_INTVAL int

	TLS_ENABLE bool
	TLS_CERT   string
	TLS_CA     string
	TLS_KEY    string

	HEADER string
	URL    string

	RUNTIME int
	ADDRESS string
	PORT    int

	BODY_LENGTH  int
	PARALLEL_NUM int
	LIMITE_RATE  int

	PROTOCAL string
	ROLE     string

	DEBUG bool
	Help  bool
)

var gStat *Stat

func init() {
	flag.IntVar(&STAT_INTVAL, "stat", 5, "stat display intval (second).")

	flag.BoolVar(&TLS_ENABLE, "tls", false, "enable/disable tls mode.")
	flag.StringVar(&TLS_CERT, "cert", "", "tls using cert file.")
	flag.StringVar(&TLS_KEY, "key", "", "tls using key file.")
	flag.StringVar(&TLS_CA, "ca", "", "tls using ca file for verify other cert.")

	flag.StringVar(&PROTOCAL, "protocal", "http", "which protocol to using (http/http2/tcp).")

	flag.StringVar(&ROLE, "role", "server", "the tools role (server/client).")
	flag.IntVar(&PARALLEL_NUM, "par", 1, "the parallel numbers to connect.")
	flag.IntVar(&RUNTIME, "runtime", 60, "total run time (second).")
	flag.IntVar(&BODY_LENGTH, "body", 1, "transport body length (KB).")

	flag.StringVar(&ADDRESS, "address", "127.0.0.1", "service/client address.")
	flag.IntVar(&PORT, "port", 8001, "service/client port.")
	flag.IntVar(&LIMITE_RATE, "limit", 0, "limit times per second to send. as[0,1000]")

	flag.StringVar(&URL, "url", "/", "set request url. as[/abc/123]")
	flag.StringVar(&HEADER, "head", "", "set request head. as[key1=value1,key2=value2]")

	flag.BoolVar(&DEBUG, "debug", false, "display debug infomation.")
	flag.BoolVar(&Help, "help", false, "usage help.")
}

func main() {

	flag.Parse()
	if Help {
		flag.Usage()
		os.Exit(1)
	}

	BODY_LENGTH = BODY_LENGTH * 1024

	gStat = NewStat(STAT_INTVAL)
	gStat.Prefix(os.Args[0])

	defer gStat.Delete()

	var isClient bool
	if 0 == strings.Compare(ROLE, "server") {
		isClient = false
	} else if 0 == strings.Compare(ROLE, "client") {
		isClient = true
	} else {
		flag.Usage()
		os.Exit(1)
	}

	var err error
	var tlscfg *tls.Config

	addr := fmt.Sprintf("%s:%d", ADDRESS, PORT)

	if TLS_ENABLE {
		var cfg *TlsConfig
		if TLS_KEY != "" && TLS_CERT != "" {
			cfg = &TlsConfig{CA: TLS_CA, Cert: TLS_CERT, Key: TLS_KEY}
		}
		if isClient == true {
			tlscfg, err = TlsConfigClient(cfg, addr)
		} else {
			tlscfg, err = TlsConfigServer(cfg)
		}

		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if 0 == strings.Compare(PROTOCAL, "tcp") {
		if isClient {
			TcpClient(tlscfg)
		} else {
			TcpServer(tlscfg)
		}
	} else if 0 == strings.Compare(PROTOCAL, "http") {
		if isClient {
			HttpClient(tlscfg)
		} else {
			HttpServer(tlscfg)
		}
	} else {
		flag.Usage()
		os.Exit(1)
	}
}
