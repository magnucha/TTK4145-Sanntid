package main

import (
	"fmt"
	"time"
	"strings"
)

func main () {
	udpPort := ":20010"
	tcpPort := ":30010"
	
	udp_conn := UDP_Create_Receive_Socket(udpPort)
	raddr := UDP_Receive(udp_conn)

	tcp_conn := TCP_Connect(raddr+tcpPort)
	for {
		fmt.Println("Sending...")
		TCP_Send(tcp_conn, "Klaska laksen")
		time.Sleep(2*time.Second)
	}
}
/*
func main () {
	tcp_conn := TCP_Connect("129.241.187.159:30010")
	for {
		fmt.Println("Sending...")
		TCP_Send(tcp_conn, "Klaska laksen")
		time.Sleep(2*time.Second)
	}

}*/
