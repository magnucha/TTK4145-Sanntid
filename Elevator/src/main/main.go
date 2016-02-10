package main

import (
	"network"
	"config"
	"time"
)

func Message_Server(ch_incoming_msg <-chan config.Message) {
	for {
		msg <- ch_incoming_msg
		switch msg.Msg_type {
			case config.ACK:
				Increment_Ack_Counter(msg)	//Not yet implemented
			case config.STATE_UPDATE:
				config.active_elevs[msg.Raddr] = msg.State
			case config.ADD_ORDER:
				ch_new_order <- msg
			case config.DELETE_ORDER:
				ch_del_order <- msg
		}
	}
}

func main(){
	ch_incoming_msg := make(chan config.Message)
	ch_outgoing_msg := make(chan config.Message)
	ch_new_order := make(chan config.Message)
	ch_del_order := make(chan config.Message)
	
	network.Network_Init(ch_outgoing_msg, ch_incoming_msg)
	
	state := config.ElevState{Is_idle: false, Direction: 1, Last_floor: 1, Destination_floor: 1}
	test := config.Message{Msg_type: 0, State: state, Pressed_button_type: 0, Pressed_button_floor: 3, Elevs_in_network_count: len(network.TCP_connections)}
	for {
		ch_outgoing_msg <- test
		time.Sleep(2*time.Second)
	}
}
