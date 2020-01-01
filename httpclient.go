package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var gBody []byte

func NewHttpClient(tlscfg *tls.Config) *http.Client {

	newTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
		TLSClientConfig: tlscfg,
	}

	return &http.Client{
		Transport: newTransport,
		Timeout:   10 * time.Second,
	}
}

type Header struct {
	key   string
	value string
}

func HttpRequest(client *http.Client, path string, header []Header, body []byte) (time.Duration, error) {

	request, err := http.NewRequest("POST", path, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}

	for _, v := range header {
		request.Header.Add(v.key, v.value)
	}

	tmBefore := time.Now()

	if DEBUG {
		headers := fmt.Sprintf("\r\nHeader:\r\n")
		for key, value := range request.Header {
			headers += fmt.Sprintf("\t%s:%v\r\n", key, value)
		}
		log.Printf("Request Method:%s\r\n Host:%s \r\nURL:%s%s\r\n",
			request.Method, request.URL.Host, request.URL.Path, headers)
	}

	rsp, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return 0, err
	}

	if DEBUG {
		headers := fmt.Sprintf("\r\nHeader:\r\n")
		for key, value := range rsp.Header {
			headers += fmt.Sprintf("\t%s:%v\r\n", key, value)
		}
		log.Printf("Response Code:%d%s%s\r\n", rsp.StatusCode, headers, body)
	}

	if rsp.StatusCode != http.StatusOK {
		return 0, errors.New("response " + rsp.Status)
	}

	tmAfter := time.Now()

	return tmAfter.Sub(tmBefore), nil
}

func HttpClient(tlscfg *tls.Config) {

	var path string
	if TLS_ENABLE {
		path = fmt.Sprintf("https://%s:%d%s", ADDRESS, PORT, URL)
	} else {
		path = fmt.Sprintf("http://%s:%d%s", ADDRESS, PORT, URL)
	}

	log.Printf("Request : %s\r\n", path)
	log.Printf("BodyLen : %d\r\n", BODY_LENGTH)
	log.Printf("PARAL   : %d\r\n", PARALLEL_NUM)

	header := make([]Header, 0)
	if HEADER != "" {
		list := strings.Split(HEADER, ",")
		for _, v := range list {
			keyvalue := strings.Split(v, "=")
			if len(keyvalue) == 2 {
				header = append(header, Header{key: keyvalue[0], value: keyvalue[1]})
			}
		}
	}

	for _, v := range header {
		log.Printf("Header  : [%s:%s] \r\n", v.key, v.value)
	}

	gBody := make([]byte, BODY_LENGTH)
	for i := 0; i < BODY_LENGTH; i++ {
		gBody[i] = byte('A')
	}

	log.Println("Http Benchmark Start!")

	var stop bool
	var wait sync.WaitGroup
	wait.Add(PARALLEL_NUM)

	for i := 0; i < PARALLEL_NUM; i++ {
		go func() {
			defer wait.Done()
			client := NewHttpClient(tlscfg)
			for {
				timestamp, err := HttpRequest(client, path, header, gBody)
				if err != nil {
					log.Println(err.Error())
					time.Sleep(5 * time.Second)
					continue
				}
				gStat.Add(len(gBody), uint64(timestamp))
				if stop {
					break
				}
				if LIMITE_RATE != 0 && LIMITE_RATE < 1000 {
					time.Sleep(time.Second / time.Duration(LIMITE_RATE))
				}
			}
		}()
	}
	time.Sleep(time.Duration(RUNTIME) * time.Second)
	stop = true
	wait.Wait()

	log.Println("Http Benchmark Stop!")
}
