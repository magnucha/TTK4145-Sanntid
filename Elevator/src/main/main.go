package main

import (
	"config"
	"hardware"
	"log"
	"network"
	"queue"
	"time"
	"fsm"
)

var ch_incoming_msg = make(chan config.Message)
var ch_outgoing_msg = make(chan config.Message)
var ch_new_order = make(chan config.Message)
var ch_del_order = make(chan config.Message)
var ch_button_pressed = make(chan config.ButtonStruct)
var ch_floor_poll = make(chan int)
var ch_new_elev = make(chan string)

//var heis config.ElevState

func main() {
	network.Network_Init(ch_outgoing_msg, ch_incoming_msg, ch_new_elev)
	time.Sleep(time.Millisecond)
	if !hardware.Elev_Init() {
		log.Fatal("Unable to initialize elevator hardware!")
	}
	go Message_Server()
	go Channel_Server()
	go hardware.Read_Buttons(ch_button_pressed)
	go hardware.Set_Lights()
	go hardware.Floor_Poller(ch_floor_poll)
	go State_Spammer()
	fsm.FSM_Init(ch_outgoing_msg)
	
	log.Printf("Elev addr: %s", config.Laddr)

	for {
		time.Sleep(5 * time.Second)
	}

}

func Message_Server() {
	for {
		msg := <-ch_incoming_msg
		switch msg.Msg_type {
		//case config.ACK:
		//	Increment_Ack_Counter(msg)	//Not yet implemented
		case config.STATE_UPDATE:
			log.Printf("State spammer received: State: %d ", msg.State.Last_floor)
			*config.Active_elevs[msg.Raddr] = msg.State
			config.Active_elevs[msg.Raddr].Timer.Reset(config.TIMEOUT) //Reset local timer at state spammer
		case config.ADD_ORDER:
			fsm.Event_Order_Received(msg.Button)
		case config.DELETE_ORDER:
			queue.Delete_Order(msg.Button.Floor, ch_outgoing_msg)
		}
	}
}

func Channel_Server() {
	for {
		select {
		case button := <-ch_button_pressed:
			ch_outgoing_msg <- config.Message{Msg_type: config.ADD_ORDER, Button: button}
			fsm.Event_Order_Received(button)
		case floor := <-ch_floor_poll:
			fsm.Event_Reached_Floor(floor, ch_outgoing_msg)
		case raddr := <-ch_new_elev:
			SetActive(raddr)
		}
	}
}

func State_Spammer(){
	for{
		ch_outgoing_msg <- config.Message{Msg_type: config.STATE_UPDATE, State: *config.Active_elevs[config.Laddr]}
		time.Sleep(500*time.Millisecond)
		//config.Active_elevs[config.Laddr].Timer.Reset(config.TIMEOUT)
	}
}


func SetActive(raddr string) {
	already_active := false
	for addr, _ := range config.Active_elevs {
		if addr == raddr {
			already_active = true
		}
	}
	killer := func(){ //Poppiloppi-kode
		delete(config.Active_elevs, raddr)
		//Redistribute orders
	}
	if !already_active {
		config.Active_elevs[raddr] = &config.ElevState{Is_idle: true, Door_open: false, Direction: config.DIR_STOP, Last_floor: -1, Timer: time.AfterFunc(config.TIMEOUT, killer)}
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