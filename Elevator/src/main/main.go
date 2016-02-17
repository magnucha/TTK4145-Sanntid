package main

import (
	"network"
	"config"
	"log"
	"time"
	"queue"
	"hardware"
)

var ch_incoming_msg = make(chan config.Message)
var ch_outgoing_msg = make(chan config.Message)
var ch_new_order = make(chan config.Message)
var ch_del_order = make(chan config.Message)
var ch_button_pressed = make(chan config.ButtonStruct)
var ch_floor_poll = make(chan int)
//var heis config.ElevState

func main(){
	network.Network_Init(ch_outgoing_msg, ch_incoming_msg)
	time.Sleep(time.Millisecond)
	if (!hardware.Elev_Init()) {
		log.Fatal("Unable to initialize elevator hardware!")
	}
	go Message_Server()
	go Channel_Server()
	go hardware.Read_Buttons(ch_button_pressed)
	go hardware.Set_Lights()
	go hardware.Floor_Poller(ch_floor_poll)

	log.Printf("Elev addr: %s", config.Laddr)





}

func Message_Server() {
	for {
		msg := <- ch_incoming_msg
		switch msg.Msg_type {
			//case config.ACK:
			//	Increment_Ack_Counter(msg)	//Not yet implemented
			case config.STATE_UPDATE:
				*config.Active_elevs[msg.Raddr] = msg.State
			case config.ADD_ORDER:
				//ch_new_order <- msg
				log.Printf("Floor: %d, Button: %d, Destination: %d Elevs: %d", msg.Button.Floor, msg.Button.Button_type, msg.State.Destination_floor, msg.Elevs_in_network_count)
			case config.DELETE_ORDER:
				queue.Delete_Order(msg.Button.Floor)
		}
	}
}

func Channel_Server(){
	for{
		select{
			case button := <- ch_button_pressed:
				queue.Add_Order(button)
				log.Printf("Pressed floor: %d, button %d", button.Floor, button.Button_type)
				ch_outgoing_msg <- config.Message{Msg_type: config.ADD_ORDER, Button: button}
			case floor := <- ch_floor_poll:
				if(queue.Check_Order(floor)){
					queue.Delete_Order(floor)
					ch_outgoing_msg <- config.Message{Msg_type: config.DELETE_ORDER, Button: config.ButtonStruct{Floor: floor}}	//Fiksien popo
					hardware.Stop_On_Floor()
				}
		}
	}
}


/*
Channel server:
	- Receive button pressed
		- Add to queue
		- Broadcast new order
	- Receive from floor poller
		- If order in same dir
			- Stop in floor
			- Delete pick-up order
			- Delete if command order on floor
			- Continue to destination
	- Completed order
		- Broadcast DELETE_ORDER
*/