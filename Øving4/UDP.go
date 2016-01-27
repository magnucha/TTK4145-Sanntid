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

//How to get the remote IP from addr?
func UDP_Receive(conn *net.UDPConn) string {
	msg := make([]byte, 1024)
	conn.ReadFromUDP(msg)

	return strings.Trim(msg,"\x00")
}
