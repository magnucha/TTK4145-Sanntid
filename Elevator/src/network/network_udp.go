package network

import (
	"config"
	"encoding/json"
	"log"
	"net"
)

func Network_Init(ch_outgoing_msg <-chan config.Message, ch_incoming_msg chan<- config.Message) {
	ch_UDP_transmit := make(chan config.NetworkMessage)
	ch_UDP_received := make(chan config.NetworkMessage, 5)
	
	UDP_broadcast_socket := UDP_Create_Send_Socket(config.UDP_BROADCAST_ADDR)
	UDP_listen_socket := UDP_Create_Listen_Socket(config.UDP_BROADCAST_ADDR[15:])
	
	//We choose to begin receiving UDP after broadcast to avoid creating a connection to ourselves
	go UDP_Send(UDP_broadcast_socket, ch_UDP_transmit)
	UDP_Broadcast_Presence(UDP_broadcast_socket)
	go UDP_Receive(UDP_listen_socket, ch_UDP_received)
	
	go Encode_And_Forward_Transmission(ch_UDP_transmit, ch_outgoing_msg)
	go Decode_And_Forward_Reception(ch_UDP_transmit, ch_UDP_received, ch_incoming_msg)
}

func Store_Local_Addr() {
	baddr, err = net.ResolveUDPAddr("udp4", config.UDP_BROADCAST_ADDR)
	tempConn, err := net.DialUDP("udp4", nil, baddr)
	tempAddr := tempConn.LocalAddr()
	laddr, err := net.ResolveUDPAddr("udp4", tempAddr.String())
	laddr.Port = localListenPort
	config.Laddr = laddr.String()
	defer tempConn.Close()
}

func Accept_UDP_Message(msg config.Message) bool {
	if msg.MessageType >= 0 && msg.Message <= 2 {
		if msg.Button_type >=0 && msg.Button_type <= 2 {
			if msg.Floor >= 0 && msg.Floor <= 4 {
				if msg.Elevs_in_network_count >= 0 && msg.Elevs_in_network_count <= 4 {
					return true
				}
			}
		}
	}
	return false
}

func Add_Active_Elev(raddr string) {
	already_active := false
	for addr, _ := range config.active_elevs {
		if addr == raddr {
			already_active = true
		}
	}
	if !already_active {
		//Initialize to an invalid state
		config.active_elevs[raddr] = config.ElevState{Is_idle: true, Direction: config.DIR_STOP, Last_floor: -1, Destination_floor: -1}
	}
}

func Encode_And_Forward_Transmission(ch_transmit chan<- config.NetworkMessage, ch_outgoing_msg <-chan config.Message) {
	for {
		msg := <- ch_outgoing_msg
		msg.Elevs_in_network_count = len(UDP_connections)
		json_msg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("UDP_Encode_And_Forward_Transmission: json error:", err)
		}
		ch_transmit <- config.NetworkMessage{Data: json_msg, Length: len(json_msg)}
	}
}

func Decode_And_Forward_Reception(ch_transmit chan<- config.NetworkMessage, ch_received <-chan config.NetworkMessage, ch_incoming_msg chan<- config.Message) {
	for {
		received := <- ch_received
				
		if received.Data == config.UDP_PRESENCE_MSG {
			Add_Active_Elev(received.Raddr)
		}
		else {
			var msg config.Message
			err := json.Unmarshal(received.Data, &msg)
			if err != nil {
				log.Printf("TCP_Decode_And_Forward_Reception: json error:", err)
			}
			if Accept_UDP_Message(incoming) {
				msg.Raddr = received.Raddr
				ch_incoming_msg <- msg
			}
		}
	}
}


