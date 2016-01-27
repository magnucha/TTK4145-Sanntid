package main

import (
	"net"
	"fmt"
	"bufio"
)

type TCPMessage struct {
	remote_addr string 	//The remote address we are sending to /receiving from
	data []byte			
	length int			//Length of received data, nil when sending
}

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

func TCP_Listen() net.Conn {
	tcp_port, _ := net.ResolveTCPAddr("tcp", ":20003")
	tcp_listener, _ := net.ListenTCP("tcp", tcp_port)
	conn,_ := tcp_listener.Accept()
	return conn
}

func TCP_Send(conn *net.TCPConn, msg string) {
	conn.Write([]byte(msg + string('\x00')))
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
