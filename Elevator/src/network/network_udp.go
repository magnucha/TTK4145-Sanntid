package network

import (
	"config"
	"encoding/json"
	"log"
	"net"
)

func Network_Init(ch_outgoing_msg <-chan config.Message, ch_incoming_msg chan<- config.Message) {
	ch_UDP_transmit := make(chan []byte)
	ch_UDP_received := make(chan config.NetworkMessage, 5)

	UDP_broadcast_socket := UDP_Create_Send_Socket(config.UDP_BROADCAST_ADDR + config.UDP_BROADCAST_PORT)
	UDP_listen_socket := UDP_Create_Listen_Socket(config.UDP_BROADCAST_PORT)
	Store_Local_Addr()

	//We choose to begin receiving UDP after broadcast to avoid creating a connection to ourselves
	go UDP_Send(UDP_broadcast_socket, ch_UDP_transmit)
	UDP_Broadcast_Presence(UDP_broadcast_socket, ch_UDP_transmit)
	go UDP_Receive(UDP_listen_socket, ch_UDP_received)

	go Encode_And_Forward_Transmission(ch_UDP_transmit, ch_outgoing_msg)
	go Decode_And_Forward_Reception(ch_UDP_transmit, ch_UDP_received, ch_incoming_msg)
}

func Store_Local_Addr() {
	baddr, _ := net.ResolveUDPAddr("udp4", config.UDP_BROADCAST_ADDR+config.UDP_BROADCAST_PORT)
	tempConn, _ := net.DialUDP("udp4", nil, baddr)
	tempAddr := tempConn.LocalAddr()
	laddr, _ := net.ResolveUDPAddr("udp4", tempAddr.String())
	config.Laddr = laddr.IP.String()
	Add_Active_Elev(config.Laddr)
	config.Local_elev = config.Active_elevs[config.Laddr]
	defer tempConn.Close()
}

func Add_Active_Elev(raddr string) {
	already_active := false
	for addr, _ := range config.Active_elevs {
		if addr == raddr {
			already_active = true
		}
	}
	if !already_active {
		config.Active_elevs[raddr] = &config.ElevState{Is_idle: true, Door_open: false, Direction: config.DIR_STOP, Last_floor: -1}
	}
}

func Encode_And_Forward_Transmission(ch_transmit chan<- []byte, ch_outgoing_msg <-chan config.Message) {
	for {
		msg := <-ch_outgoing_msg
		msg.Elevs_in_network_count = len(config.Active_elevs)
		json_msg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("UDP_Encode_And_Forward_Transmission: json error:", err)
		}
		ch_transmit <- append([]byte(config.MESSAGE_PREFIX), json_msg...)
	}
}

func Decode_And_Forward_Reception(ch_transmit chan<- []byte, ch_received <-chan config.NetworkMessage, ch_incoming_msg chan<- config.Message) {
	for {
		received := <-ch_received
		if string(received.Data)[:len(config.UDP_PRESENCE_MSG)] == config.UDP_PRESENCE_MSG {
			Add_Active_Elev(received.Raddr)
		} else {
			var msg config.Message
			err := json.Unmarshal(received.Data[len(config.MESSAGE_PREFIX):received.Length], &msg)
			if err != nil {
				log.Printf("TCP_Decode_And_Forward_Reception: json error:", err)
			}
			msg.Raddr = received.Raddr
			ch_incoming_msg <- msg
		}
	}
}
