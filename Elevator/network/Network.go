package network

import (
	"net"
	"fmt"
	"config"
)

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
	
	ch_TCP_transmit := make(chan config.NetworkMessage)
	ch_TCP_received := make(chan config.NetworkMessage, 5) //Er dette en passende størrelse for bufferet?
	ch_UDP_transmit := make(chan config.NetworkMessage)
	ch_UDP_received := make(chan config.NetworkMessage, 5)
	ch_TCP_new_connection := make(chan net.TCPConn)
	
	UDP_receive_socket = UDP_Create_Receive_Socket(config.UDP_broadcast_addr[15:])
	
	go UDP_Receive(UDP_receive_socket)
	go TCP_Listen_And_Store_Conn(ch_TCP_new_connection), ch_UDP_received)
	
	
	
}
