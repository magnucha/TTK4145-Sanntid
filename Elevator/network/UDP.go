package network

import (
	"net"
	"config"
	"log"
	"time"
)

func UDP_Create_Send_Socket(addr string) *net.UDPConn{
	UDPaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil{
		log.Printf(err)
	}

	connection, err := net.DialUDP("udp", nil, UDPaddr)
	if err != nil{
		log.Printf("DialUDP error: %s",err)
	}
	return connection
}

func UDP_Create_Receive_Socket(port string) *net.UDPConn {
	UDPaddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil{
		log.Printf("ResolveUDPAddr error: %s", err)
	}

	connection, err := net.ListenUDP("udp", UDPaddr)
	if err != nil{
		log.Printf("ListenUDP error: %s", err)
	}
	return connection
}

func UDP_Broadcast_Presence(conn *net.UDPConn) {
	for i=0; i<10; i++ {
		UDP_Send(conn, config.UDP_PRESENCE_MSG)
		time.Sleep(100*time.Millisecond)
	}
}

func UDP_Send(conn *net.UDPConn, msg string){
	_, err := conn.Write([]byte(msg))
	if err != nil{
		fmt.Printf("UDP write error: %s", err)
	}
}

func UDP_Receive(conn *net.UDPConn, ch_received chan<- config.NetworkMessage) {
	for {
		msg := make([]byte, 1024)
		length, raddr, _ := conn.ReadFromUDP(msg)
		received_msg := config.NetworkMessage{raddr: raddr.IP.String(), data: msg, length: length}
		ch_received <-received_msg
	}
}
