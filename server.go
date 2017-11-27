package main

import (
	"net/rpc"
	"net"
	"time"
	"crypto/tls"
	"log"
	"io/ioutil"
	"crypto/x509"
	"crypto/rand"
)

type Timer int

var acceptTime time.Time


func (r *Timer) ClockSynchronize(S string, ntpTimes *([]time.Time)) error { // 接收者 + 函数名 + 参数 + 返回类型
	*ntpTimes = append(*ntpTimes, acceptTime) //append t2
	time.Sleep(5*time.Second)
	*ntpTimes = append(*ntpTimes, time.Now()) //append t3
	return nil
}

func serve(conn net.Conn) {
	defer conn.Close()
	rpc.ServeConn(conn)
	log.Println("serve: connection closed")
}

func main() {
	/*Load server cert and key*/
	cert, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Println("server: load certificate & key: ", err)
		return
	}


	/*Build client cert pool*/
	caCrt, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		log.Println("server: load ca.crt: ", err)
		return
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCrt)


	/*Establish LISTENER*/
	config := tls.Config {
		ClientAuth: tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs: certPool,
	}
	config.Rand = rand.Reader
	rpc.Register(new(Timer))
	rpc.HandleHTTP()
	listener, err := tls.Listen("tcp", "114.212.85.197:10255", &config)
	if err != nil {
		log.Println("server: listen: ", err)
		return
	}
	log.Println("server: listening")
	defer listener.Close()

	
	/*Build connection and serve RPC*/
	for {
		conn, err := listener.Accept()
		acceptTime = time.Now()
		if err != nil {
			log.Println("server: accept: ", err)
			continue
		}
		log.Println("server: accept from: ", conn.RemoteAddr())
		go serve(conn)
	}
}