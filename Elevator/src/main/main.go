package main

import (
	"config"
	"fsm"
	"hardware"
	"log"
	"network"
	"queue"
	"time"
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
		time.Sleep(3 * time.Second)
		/*for _, elev := range config.Active_elevs {
			log.Printf("Direction: %d, Num elevs: %d", elev.Last_floor, len(config.Active_elevs))
		}*/
	}

}

func Message_Server() {
	for {
		msg := <-ch_incoming_msg
		switch msg.Msg_type {
		//case config.ACK:
		//	Increment_Ack_Counter(msg)	//Not yet implemented
		case config.STATE_UPDATE:
			already_active := false
			for addr, _ := range config.Active_elevs {
				if msg.Raddr == addr {
					already_active = true
					break
				}
			}
			if !already_active {
				ch_new_elev <- msg.Raddr
				time.Sleep(10 * time.Microsecond)
			}
			State_Copy(config.Active_elevs[msg.Raddr], &msg.State)
			config.Active_elevs[msg.Raddr].Timer.Reset(config.TIMEOUT)
		case config.ADD_ORDER:
			fsm.Event_Order_Received(msg.Button)
		case config.DELETE_ORDER:
			queue.Delete_Order(msg.Button.Floor, ch_outgoing_msg, false)
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
			Set_Active(raddr)
		}
	}
}

func State_Spammer() {
	for {
		time.Sleep(500 * time.Millisecond)
		ch_outgoing_msg <- config.Message{Msg_type: config.STATE_UPDATE, State: *config.Active_elevs[config.Laddr]}
	}
}

func Set_Active(raddr string) {
	for addr, _ := range config.Active_elevs {
		if addr == raddr {
			return
		}
	}
	killer := func() {
		delete(config.Active_elevs, raddr)
		//Redistribute orders
	}
	config.Active_elevs[raddr] = &config.ElevState{Is_idle: true, Door_open: false, Direction: config.DIR_STOP, Last_floor: -1, Timer: time.AfterFunc(config.TIMEOUT, killer)}
}

func State_Copy(a *config.ElevState, b *config.ElevState) {
	a.Is_idle = b.Is_idle
	a.Door_open = b.Door_open
	a.Direction = b.Direction
	a.Last_floor = b.Last_floor
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


Bestillingsutdeling:
	- GetCost(Order, ElevState) int	//Returnerer cost til en bestilling for den gitte heisen
	- AssignElev(Order)				//Finner beste heis, og setter Queue.Addr
	- ReassignOrders(addr string) 	//Kaller AssignElev på alle bestillinger med addr=addr








*/
