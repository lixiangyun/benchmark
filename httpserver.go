package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
	// "golang.org/x/net/http2"
)

type DemoHttp struct{}

func (*DemoHttp) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	if DEBUG {
		log.Printf("RemoteAddr:%s, Url:%s, Header:%v\n", req.RemoteAddr, req.URL.String(), req.Header)
	}

	gStat.Add(len(body), 0)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(body))
}

func HttpServer(tlscfg *tls.Config) {

	addr := fmt.Sprintf("%s:%d", ADDRESS, PORT)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("http listen failed!", err.Error())
		return
	}

	log.Printf("Http Proxy Listen %s\r\n", addr)

	if tlscfg != nil {
		lis = tls.NewListener(lis, tlscfg)
	}

	svc := &http.Server{
		Handler:      &DemoHttp{},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig:    tlscfg}

	/*
		if protc == PROTO_HTTP {
			err = svc.Serve(lis)
		} else {
			if DEBUG {
				http2.VerboseLogs = true
			}
			http2.ConfigureServer(svc, &http2.Server{})
			err = svc.ServeTLS(lis, "", "")
		}*/

	err = svc.Serve(lis)
	if err != nil {
		log.Println(err.Error())
	}
}
