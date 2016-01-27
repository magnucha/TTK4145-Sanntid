package main

import (
	//"time"
	"fmt"
	//"net"
)

func main(){
	BROADCAST_UDP := "129.241.187.255:20010"
	UDP_send_connection := UDP_Create_Send_Socket(BROADCAST_UDP)
	
	
//	for i := 1; i < 10; i++	{
		UDP_Send(UDP_send_connection, "129.241.187.159")
//		time.Sleep(time.Second)
//	}
	fmt.Println("Waiting for TCP connection")
	TCP_connection := TCP_Listen()
	fmt.Println("Waiting for TCP message")	
	for{
		fmt.Println(TCP_Receive(TCP_connection))
	}
}
