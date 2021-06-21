package main

import (
	"container/list"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

type s_client struct {
	num  int
	name string
	*net.UDPAddr
}

var g_client = list.New()

//var g_clients map[uint]list //uint 是ssrc，list是请求的链表
var limitChan = make(chan bool, 1000)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	//返回服务器的端口地址，可以有多个服务器来进行负载均衡
	io.WriteString(w, "127.0.0.1:8080")
}

// UDP goroutine concurrency to read UDP maybe not parallelism,maybe in onethread maybe in multi thread,maybe yes,maybe no
func udpProcess(conn *net.UDPConn) {
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Println("Failed To Read UDP Msg, Error: " + err.Error())
	}

	//var flag = 0x80
	if data[0] == 0x80 { //RTP 协议
		if n < 13 {
			fmt.Println("client error")
		} else {
			//转发
			conn.WriteToUDP([]byte("test"), remoteAddr)
		}
		//求取RTP ssrc
	}

	//Call Endc Resp Read
	str := string(data[0:3])
	if str == "Call" { //单点找人

	} else if str == "Pull" { //拉取流
		fmt.Println("pull")
		//ssrc := 1234
		//l := list.New()
		//l.PushBack(s_client{1,"name",remoteAddr})
		//g_clients[ssrc] = l
		g_client.PushBack(s_client{1, "name", remoteAddr})
	} else if str == "Resp" { //

	} else if str == "Chat" {

	} else if str == "Endc" {

	} else {

	}

	//str := string(data[:n])
	//fmt.Println("Reveive From Client, Data: " + str)
	<-limitChan
}

func udpServer(address string) {
	fmt.Println("server start at udp 8080")
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()

	if err != nil {
		fmt.Println("Read From Connect Failed, Err :" + err.Error())
		os.Exit(1)
	}

	for {
		limitChan <- true
		go udpProcess(conn)
	}

}

func httpServer(address string) {
	http.HandleFunc("/push", HelloServer)
	fmt.Println("server http start at 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func main11() {

	//g_clients :=make(map[uint]list)
	//client[32]="qianbo"
	address := "0.0.0.0:8080"
	go httpServer(address)

	udpServer(address)
}
