package main

import (
	"net"
	"fmt"
	"bufio"
	"config"
)

func TCP_Connect(IP string) *net.TCPConn {
	//Get the servers TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp", IP)
	if err != nil {
		fmt.Println("ResolveTCPAddr failed: ", err.Error())
		return nil
	}
	
	//Connect to the TCP server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("DialTCP failed: ", err.Error())
		conn = nil
	}
	return conn
}

func TCP_Listen() {
	tcp_port, _ := net.ResolveTCPAddr("tcp", ":20003")
	tcp_listener, _ := net.ListenTCP("tcp", tcp_port)
	for {
		conn,_ := tcp_listener.Accept()
		ch_TCP_new_connection <- conn
	}
}

func TCP_Transmit(conn *net.TCPConn, ch_transmit <-chan NetworkMessage) {
	for {
		msg <- ch_transmit
		conn.Write([]byte(msg.data + string('\x00')))
	}
}

func TCP_Receive(conn *net.TCPConn) string {
	msg, _ := bufio.NewReader(conn).ReadString(byte('\x00'))
	return msg
}

/*
UDP:
	- Send(*UDPConnection)
	- msg, addr = Receive(*UDPConnection)

TCP:
	- *TCPConnection := TCP_Connect(IP)
	- Send(msg*TCPConnection)
	- msg, addr := Receive(*TCPConnection)

Network:
	- Array_Serialize(list int)
	- Array_Deserialize(list int)
	- *TCP_connection := Connect_to_new_elevator(IP)
	- IP string := Listen_for_new_elevator()
*/
