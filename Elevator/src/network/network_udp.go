package network

import (
	"config"
	"encoding/json"
	"log"
	"net"
	"time"
)

type ACK_Timer struct {
	cnt int
	timer *time.Timer
}

var message_log = make(map[string]*ACK_Timer)

func Init(ch_outgoing_msg chan config.Message, ch_incoming_msg chan<- config.Message, ch_new_elev chan<- string, ch_main_alive chan<- bool) {
	ch_UDP_transmit := make(chan []byte)
	ch_UDP_received := make(chan config.NetworkMessage, 5)

	UDP_broadcast_socket := UDP_Create_Send_Socket(config.UDP_BROADCAST_ADDR + config.UDP_BROADCAST_PORT)
	UDP_listen_socket := UDP_Create_Listen_Socket(config.UDP_BROADCAST_PORT)
	Store_Local_Addr()
	UDP_alive_socket := UDP_Create_Send_Socket(config.Laddr + config.UDP_ALIVE_PORT)

	//We choose to begin receiving UDP after broadcast to avoid creating a connection to ourselves
	go UDP_Send(UDP_broadcast_socket, ch_UDP_transmit)
	go UDP_Alive_Spam(UDP_alive_socket)
	UDP_Broadcast_Presence(UDP_broadcast_socket, ch_UDP_transmit)
	go UDP_Receive(UDP_listen_socket, ch_UDP_received)

	go Encode_And_Forward_Transmission(ch_UDP_transmit, ch_outgoing_msg)
	go Decode_And_Forward_Reception(ch_new_elev, ch_UDP_received, ch_incoming_msg, ch_main_alive)
}

func Store_Local_Addr() {
	baddr, _ := net.ResolveUDPAddr("udp4", config.UDP_BROADCAST_ADDR+config.UDP_BROADCAST_PORT)
	tempConn, _ := net.DialUDP("udp4", nil, baddr)
	tempAddr := tempConn.LocalAddr()
	laddr, _ := net.ResolveUDPAddr("udp4", tempAddr.String())
	config.Laddr = laddr.IP.String()
	defer tempConn.Close()
}


func Encode_And_Forward_Transmission(ch_transmit chan<- []byte, ch_outgoing_msg chan config.Message) {
	for {
		msg := <-ch_outgoing_msg
		msg.Elevs_in_network_count = len(config.Active_elevs)
		json_msg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("UDP_Encode_And_Forward_Transmission: json error:", err)
			continue
		}

		if (msg.Msg_type == config.ADD_ORDER || msg.Msg_type == config.DELETE_ORDER) && len(config.Active_elevs) > 1 {
			retransmit := func() {
				ch_outgoing_msg <- msg
				delete(message_log, string(json_msg))
				log.Println("Retransmission!")
			}
			message_log[string(json_msg[14:])] = &ACK_Timer{cnt: 0, timer: time.AfterFunc(config.TIMEOUT_UDP, retransmit)}
		}

		ch_transmit <- append([]byte(config.MESSAGE_PREFIX), json_msg...)
	}
}

func Decode_And_Forward_Reception(ch_new_elev chan<- string, ch_received <-chan config.NetworkMessage, ch_incoming_msg chan<- config.Message, ch_main_alive chan<- bool) {
	for {
		received := <-ch_received
		if string(received.Data[:len(config.MESSAGE_PREFIX)]) != config.MESSAGE_PREFIX || received.Raddr == config.Laddr {
			continue
		}

		if string(received.Data)[:len(config.UDP_PRESENCE_MSG)] == config.UDP_PRESENCE_MSG {
			ch_new_elev <- received.Raddr
		} else {
			var msg config.Message
			err := json.Unmarshal(received.Data[len(config.MESSAGE_PREFIX):received.Length], &msg)
			if err != nil {
				log.Printf("UDP_Decode_And_Forward_Reception: json error: %s", err)
				continue
			}
			if (msg.Msg_type == config.ACK) {
				Incremement_ACK_Counter(string(received.Data[len(config.MESSAGE_PREFIX)+14:received.Length]))
			}
			msg.Raddr = received.Raddr
			ch_incoming_msg <- msg
		}
	}
}

func Incremement_ACK_Counter(key string) {
	if message_log[key] != nil {
		message_log[key].cnt++
		if message_log[key].cnt >= len(config.Active_elevs)-1 {
			message_log[key].timer.Stop()
			delete(message_log, key)
		}
	}
}