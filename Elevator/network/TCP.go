package network

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

func TCP_Listen_And_Store_Conn() {
	tcp_port, _ := net.ResolveTCPAddr("tcp", ":20003")
	tcp_listener, _ := net.ListenTCP("tcp", tcp_port)
	for {
		conn,_ := tcp_listener.Accept()
		append(config.TCP_connections, conn)
	}
}

func TCP_Broadcast(ch_transmit <-chan config.NetworkMessage) {
	for i=0; i<len(config.TCP_connections); i++ {
		TCP_Trasmit(config.TCP_connections[i], ch_transmit)
	}
}

func TCP_Transmit(conn *net.TCPConn, ch_transmit <-chan config.NetworkMessage) {
	for {
		msg <- ch_transmit
		append(msg.data, byte('\x00')
		conn.Write(msg.data)
	}
}

func TCP_Receive(conn *net.TCPConn, ch_received <-chan config.NetworkMessage) {
	for {
		msg, _ := bufio.NewReader(conn).ReadBytes(byte('\x00'))
		received_msg := config.NetworkMessage{raddr: conn.RemoteAddr(), data: msg, length: len(msg)}
		ch_received <- received_msg
	}
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
