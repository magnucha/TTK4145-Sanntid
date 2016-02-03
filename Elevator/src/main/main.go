package main

import (
	"network"
	"config"
	"time"
)

func main(){
	ch_incoming_msg := make(chan config.Message)
	ch_outgoing_msg := make(chan config.Message)
	network.Network_Init(ch_outgoing_msg, ch_incoming_msg)
	
	state := config.ElevState{Is_idle: false, Direction: 1, Last_floor: 1, Destination_floor: 1}
	test := config.Message{Msg_type: 0, State: state, Pressed_button_type: 0, Pressed_button_floor: 3, Elevs_in_network_count: len(network.TCP_connections)}
	for {
		ch_outgoing_msg <- test
		time.Sleep(2*time.Second)
	}
}
