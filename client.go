package main

import (
	"fmt"
	"net/rpc"
	"time"
	"crypto/tls"
	"log"
	"crypto/x509"
	"io/ioutil"
)

func communicationDelay(t1 time.Time, t2 time.Time, t3 time.Time, t4 time.Time) time.Duration {
	return ( t2.Sub(t1) + t3.Sub(t4) ) / 2
}

func main() {
	/*Load client cert and key*/
	cert, err := tls.LoadX509KeyPair("certs/client.crt", "certs/client.key")
	if err != nil {
		log.Println("client: load certificate & key: ", err)
		return
	}


	/*Build client certPool*/
	caCrt, err := ioutil.ReadFile("certs/ca.crt")
	if err != nil {
		log.Println("client: load ca.crt: ", err)
		return
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCrt)


	/*Build connection*/
	config := tls.Config {
		Certificates: []tls.Certificate{cert},
		RootCAs: certPool,
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", "114.212.85.197:10255", &config)
	if err != nil {
		log.Println("Client dial: ", err)
		return
	}
	defer conn.Close()


	/*Remote Procedure Call*/
	var ntpTimes []time.Time
	startTime := time.Now()
	err = rpc.NewClient(conn).Call("Timer.ClockSynchronize", "Hello server", &ntpTimes)
	if err != nil {
		log.Println("client: Fall to call Timer RPC: ", err)
	}
	ntpTimes = append(ntpTimes, startTime, time.Now()) //append T1, T4


	/*Result*/
	t1, t2, t3, t4 := ntpTimes[2], ntpTimes[0], ntpTimes[1], ntpTimes[3]
	//fmt.Println("t1: ",t1, "\nt2: ", t2, "\nt3: ", t3, "\nt4: ", t4)
	timeOffset := communicationDelay(t1, t2, t3, t4)
	fmt.Println("Time offset: ", timeOffset)
	localTime := time.Now()
	serverTime := localTime.Add(timeOffset)
	fmt.Println("Local time:", localTime, "\nServer time: ", serverTime)
}

