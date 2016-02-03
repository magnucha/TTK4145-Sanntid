package network

import (
	"net"
	"fmt"
	"config"
	"encoding/json"
	"log"
)

var TCP_connections = make([]config.NetworkMessage)

func Network_Init(ch_outgoing_msg <-chan config.Message, ch_incoming_msg chan<- config.Message) {
	
	ch_TCP_transmit := make(chan config.NetworkMessage)
	ch_TCP_received := make(chan config.NetworkMessage, 5)
	ch_UDP_transmit := make(chan config.NetworkMessage)
	ch_UDP_received := make(chan config.NetworkMessage, 5)
	
	UDP_broadcast_socket = UDP_Create_Send_Socket(config.UDP_BROADCAST_ADDR)
	UDP_receive_socket = UDP_Create_Receive_Socket(config.UDP_BROADCAST_ADDR[15:])
	
	//Set up routines for creating new connections
	go UDP_Receive(UDP_receive_socket, ch_UDP_received)
	go TCP_Listen_And_Store_Conn()
	go Connect_TCP_On_UDP_Message(ch_UDP_received)
	
	
	go TCP_Encode_And_Forward_Transmission(ch_TCP_transmit, ch_outgoing_msg)
	go TCP_Decode_And_Forward_Reception(ch_TCP_received, ch_incoming_msg)
	
	UDP_Broadcast_Presence(UDP_broadcast_socket)
}

func Connect_TCP_On_UDP_Message(ch_UDP_received <-chan config.NetworkMessage) {
	for {
		msg := <- ch_UDP_received
		if msg == config.UDP_PRESENCE_MSG {
			append(TCP_connections, TCP_Connect(msg.raddr[:15]))
		}
	}
}

func TCP_Encode_And_Forward_Transmission(ch_transmit chan<- config.NetworkMessage, ch_outgoing_msg <-chan config.Message) {
	for {
		msg := <- ch_transmit
		msg.elevs_in_network_count = len(TCP_connections)
		
		json_msg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("TCP_Encode_And_Forward_Transmission: json error:", err)
		}
		ch_transmit <- NetworkMessage{addr: "", data: json_msg, length: len(json_msg))
	}
}

func TCP_Decode_And_Forward_Reception(ch_received <-chan config.NetworkMessage, ch_incoming_msg chan<- config.Message) {
	for {
		received_data := <- ch_received
		
		var incoming config.Message 
		err := json.Unmarshal(received_data.data, &incoming)
		if err != nil {
			log.Printf("TCP_Decode_And_Forward_Reception: json error:", err)
		}
		
		ch_incoming_msg <- incoming
	}
}





