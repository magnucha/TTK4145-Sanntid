package main

import (
	"network"
	"config"
	"time"
	"log"
)

func main(){
	ch_incoming_msg := make(chan config.Message)
	ch_outgoing_msg := make(chan config.Message)
	network.Network_Init(ch_incoming_msg, ch_outgoing_msg)
	
	state := config.ElevState{Is_idle: false, Direction: 1, Last_floor: 1, Destination_floor: 1}
	test := config.Message{Msg_type: 0, State: state, Pressed_button_type: 0, Pressed_button_floor: 0, Elevs_in_network_count: len(network.TCP_connections)}
	for {
		select {
			rec <- ch_incoming_msg:
			    log.Printf("Incoming: {Floor: %s, ButtonType: %s", rec.Pressed_button_floor, rec.Pressed_button_type)
		}
		ch_outgoing_msg <- test
		time.Delay(2*time.Second)
	}
}
