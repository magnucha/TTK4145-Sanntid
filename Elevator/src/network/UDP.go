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
		log.Printf(err.Error())
	}

	connection, err := net.DialUDP("udp", nil, UDPaddr)
	if err != nil{
		log.Printf("DialUDP error: %s",err.Error())
	}
	return connection
}

func UDP_Create_Receive_Socket(port string) *net.UDPConn {
	UDPaddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil{
		log.Printf("ResolveUDPAddr error: %s", err.Error())
	}

	connection, err := net.ListenUDP("udp", UDPaddr)
	if err != nil{
		log.Printf("ListenUDP error: %s", err.Error())
	}
	return connection
}

func UDP_Broadcast_Presence(conn *net.UDPConn) {
	for i:=0; i<1; i++ {
		UDP_Send(conn, config.UDP_PRESENCE_MSG)
		time.Sleep(100*time.Millisecond)
	}
}

func UDP_Send(conn *net.UDPConn, msg string){
	_, err := conn.Write([]byte(msg))
	if err != nil{
		log.Printf("UDP write error: %s", err.Error())
	}
}

func UDP_Receive(conn *net.UDPConn, ch_received chan<- config.NetworkMessage) {
	for {
		msg := make([]byte, 1024)
		length, raddr, _ := conn.ReadFromUDP(msg)
		received_msg := config.NetworkMessage{Raddr: raddr.IP.String(), Data: msg, Length: length}
		ch_received <- received_msg
	}
}
