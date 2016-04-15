package network

import (
	"net"
	"config"
	"log"
	"time"
)

func UDPCreateSendSocket(addr string) *net.UDPConn{
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

func UDPCreateListenSocket(port string) *net.UDPConn {
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

func UDPSend(conn *net.UDPConn, ch_transmit <-chan []byte){
	for {
		msg := <- ch_transmit
		_, err := conn.Write(msg)
		if err != nil{
			log.Printf("UDP write error: %s", err.Error())
		}
	}
}

func UDPReceive(conn *net.UDPConn, ch_received chan<- config.NetworkMessage) {
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

func UDPAliveSpam(conn *net.UDPConn) {
	time.Sleep(2*time.Second)
	for {
		msg := config.UDP_BACKUP_MSG
		_, err := conn.Write([]byte(msg))
		if err != nil{
			log.Printf("UDP alive error: %s", err.Error())
			return
		}
		time.Sleep(200*time.Millisecond)
	}
}