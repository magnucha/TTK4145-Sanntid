package main

import (
	"network"
	"config"
	"log"
	"time"
)

var ch_incoming_msg = make(chan config.Message)
var ch_outgoing_msg = make(chan config.Message)
var ch_new_order = make(chan config.Message)
var ch_del_order = make(chan config.Message)

func main(){
	network.Network_Init(ch_outgoing_msg, ch_incoming_msg)
	go Message_Server()
	
	state := config.ElevState{Is_idle: true, Direction: 0, Last_floor: -1, Destination_floor: -1}
	for {
		ch_outgoing_msg <- config.Message{Msg_type: config.ADD_ORDER, State: state, Button_type: config.BUTTON_CALL_UP, Floor: 1}
		time.Sleep(2*time.Second)
	}
}

func Message_Server() {
	for {
		msg := <- ch_incoming_msg
		switch msg.Msg_type {
			//case config.ACK:
			//	Increment_Ack_Counter(msg)	//Not yet implemented
			case config.STATE_UPDATE:
				config.Active_elevs[msg.Raddr] = msg.State
			case config.ADD_ORDER:
				//ch_new_order <- msg
				log.Printf("Floor: %d, Button: %d, Destination: %d Elevs: %d", msg.Floor, msg.Button_type, msg.State.Destination_floor, msg.Elevs_in_network_count)
			case config.DELETE_ORDER:
				ch_del_order <- msg
		}
	}
}
