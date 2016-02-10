package network

import (
	"net"
	"bufio"
	"config"
	"log"
)

func TCP_Connect(IP string) net.TCPConn {
	//Get the servers TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp", IP+config.TCP_PORT)
	if err != nil {
		log.Printf("ResolveTCPAddr failed: %s", err.Error())
	}
	
	//Connect to the TCP server
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("DialTCP failed: %s", err.Error())
	}
	return *conn
}

func TCP_Listen_And_Store_Conn() {
	tcp_port, _ := net.ResolveTCPAddr("tcp", config.TCP_PORT)
	tcp_listener, _ := net.ListenTCP("tcp", tcp_port)
	for {
		conn,_ := tcp_listener.AcceptTCP()
		TCP_connections[*conn.RemoteAddr().String()] = *conn
		log.Printf("TCP connection made to %s!", conn.RemoteAddr())
	}
}

func TCP_Broadcast(ch_transmit <-chan config.NetworkMessage) {
	for {
		msg := <- ch_transmit
		for i:=0; i<len(TCP_connections); i++ {
			TCP_Transmit(&TCP_connections[i], msg)
			//log.Printf("Sent %s to %s", msg.Data, TCP_connections[i].RemoteAddr().String())
		}
	}
}

func TCP_Transmit(conn *net.TCPConn, msg config.NetworkMessage) {
	log.Printf("Transmitting...")
	//msg.Data = append(msg.Data, byte('\x00'))
	msg.Data = append([]byte("Halla paa lan, big papi!"), byte('\x00'))
	conn.Write(msg.Data)
}

func TCP_Receive(conn *net.TCPConn, ch_received chan<- config.NetworkMessage) {
	for {
		msg, _ := bufio.NewReader(conn).ReadBytes(byte('\x00'))
		received_msg := config.NetworkMessage{Raddr: conn.RemoteAddr().String(), Data: msg, Length: len(msg)}
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
