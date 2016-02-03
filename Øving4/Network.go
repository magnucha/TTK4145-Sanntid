package main

import (
	"net"
	"fmt"
	"config"
)

type NetworkMessage struct {
	raddr string 	//The remote address we are sending to /receiving from, on form IP:port
	data []byte			
	length int			//Length of received data, don't care when transmitting
}

/*
Function flow:
	- Create variables and channels
		- Channel:
			- What to transmit
			- What is received
	- go TCP_listen()
	- go UDP_broadcast_precence()
	- 
*/

func network_init() {
	conn_broadcast := TCP_Connect(config.TCP_broadcast_addr)
	
	ch_TCP_transmit := make(chan NetworkMessage)
	ch_TCP_received := make(chan NetworkMessage, 5) //Er dette en passende størrelse for bufferet?
	ch_UDP_received := make(chan NetworkMessage, 5)
	ch_TCP_new_connection := make(chan net.TCPConn)
	
	
	UDP_receive_socket = UDP_Create_Receive_Socket(config.UDP_broadcast_addr[15:])
	go UDP_receive(UDP_receive_socket)
	go TCP_listen(ch_TCP_new_connection), ch_UDP_received)
	
}


