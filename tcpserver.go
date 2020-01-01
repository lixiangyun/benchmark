package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

func ServerProc(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, BODY_LENGTH)

	for {
		readcnt, err := conn.Read(buf[0:])
		if err != nil {
			log.Println(err.Error())
			return
		}
		gStat.Add(readcnt, 0)

		var sendcnt int
		for {
			cnt, err := conn.Write(buf[sendcnt:readcnt])
			if err != nil {
				log.Println(err.Error())
				return
			}
			sendcnt += cnt
			if sendcnt >= readcnt {
				break
			}
		}
	}
}

func TcpServer(tlscfg *tls.Config) {

	addr := fmt.Sprintf("%s:%d", ADDRESS, PORT)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if tlscfg != nil {
		listen = tls.NewListener(listen, tlscfg)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go ServerProc(conn)
	}
}
