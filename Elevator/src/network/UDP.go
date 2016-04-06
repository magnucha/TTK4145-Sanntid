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

func UDP_Create_Listen_Socket(port string) *net.UDPConn {
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

func UDP_Broadcast_Presence(conn *net.UDPConn, ch_transmit chan []byte) {
	for i:=0; i<3; i++ {
		ch_transmit <- []byte(config.UDP_PRESENCE_MSG)
		time.Sleep(time.Millisecond)
	}
}

func UDP_Send(conn *net.UDPConn, ch_transmit <-chan []byte){
	for {
		msg := <- ch_transmit
		_, err := conn.Write(msg)
		if err != nil{
			log.Printf("UDP write error: %s", err.Error())
		}
	}
}

func UDP_Receive(conn *net.UDPConn, ch_received chan<- config.NetworkMessage) {
	msg := make([]byte, 1024)
	for {
		length, raddr, err := conn.ReadFromUDP(msg)
		if err != nil{
			log.Printf("UDP read error: %s", err.Error())
			return
		}
		received_msg := config.NetworkMessage{Raddr: raddr.IP.String(), Data: msg, Length: length}
		ch_received <- received_msg
	}
}

func UDP_Alive_Spam(conn *net.UDPConn) {
	time.Sleep(2*time.Second)
	for {
		msg := config.UDP_PRESENCE_MSG
		_, err := conn.Write([]byte(msg))
		if err != nil{
			log.Printf("UDP alive error: %s", err.Error())
			return
		}
		time.Sleep(200*time.Millisecond)
	}
}