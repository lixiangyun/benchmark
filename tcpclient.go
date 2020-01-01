package main

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func ClientSend(conn net.Conn, wait *sync.WaitGroup) {
	defer conn.Close()
	defer wait.Done()
	var sendidx uint64

	buf := make([]byte, BODY_LENGTH)
	num := BODY_LENGTH / 8

	for {

		for i := 0; i < num; i++ {
			sendidx++
			binary.BigEndian.PutUint64(buf[i*8:(i+1)*8], sendidx)
		}

		var sendcnt int
		for {
			cnt, err := conn.Write(buf[:])
			if err != nil {
				log.Println(err.Error())
				return
			}
			sendcnt += cnt
			if sendcnt >= len(buf) {
				break
			}
		}
	}
}

func ClientRecv(conn net.Conn, wait *sync.WaitGroup) {
	defer conn.Close()
	defer wait.Done()

	var sendidx uint64
	buf := make([]byte, BODY_LENGTH)

	var remain int

	for {
		cnt, err := conn.Read(buf[remain:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		remain += cnt

		if remain%8 == 0 {
			num := remain / 8
			for i := 0; i < num; i++ {
				idx := binary.BigEndian.Uint64(buf[i*8 : (i+1)*8])
				if idx != sendidx+1 {
					log.Fatalln("recv err body data!", idx, sendidx)
				}
				sendidx = idx
			}

			gStat.Add(remain, 0)
			remain = 0
		}
	}
}

func ClientConn(addr string, tlscfg *tls.Config, client *sync.WaitGroup) {
	defer client.Done()
	var wait sync.WaitGroup

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if tlscfg != nil {
		conn = tls.Client(conn, tlscfg)
	}

	wait.Add(2)

	go ClientSend(conn, &wait)
	go ClientRecv(conn, &wait)

	go func() {
		time.Sleep(time.Duration(RUNTIME) * time.Second)
		conn.Close()
	}()

	wait.Wait()
}

func TcpClient(tlscfg *tls.Config) {
	addr := fmt.Sprintf("%s:%d", ADDRESS, PORT)
	var wait sync.WaitGroup
	wait.Add(PARALLEL_NUM)
	for i := 0; i < PARALLEL_NUM; i++ {
		go ClientConn(addr, tlscfg, &wait)
	}
	wait.Wait()
}
