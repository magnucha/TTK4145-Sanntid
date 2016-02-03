package main

import (
	"net"
	"fmt"
)


func UDP_Create_Send_Socket(addr string) *net.UDPConn{
	UDPaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil{
		fmt.Println(err)
	}

	connection, err := net.DialUDP("udp", nil, UDPaddr)
	if err != nil{
		fmt.Println(err)
	}
	return connection
}

func UDP_Create_Receive_Socket(port string) *net.UDPConn {
	fmt.Println("Creating UDP Socket..")
	UDPaddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil{
		fmt.Println(err)
	}

	connection, err := net.ListenUDP("udp", UDPaddr)
	if err != nil{
		fmt.Println(err)
	}
	return connection
}

func UDP_Send(conn *net.UDPConn, msg string){
	_, err := conn.Write([]byte(msg))
	if err != nil{
		fmt.Println(err)
	}
}

func UDP_Receive(conn *net.UDPConn, ch_received <-chan config.NetworkMessage) {
	for {
		msg := make([]byte, 1024)
		length, raddr, _ := conn.ReadFromUDP(msg)
		received_msg := config.NetworkMessage{raddr: raddr.IP.String(), data: msg, length: length}
		ch_received <-received_msg
	}
}
