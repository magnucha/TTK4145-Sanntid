package config

import (
	"hardware"
)

const NUM_FLOORS = 4
const NUM_BUTTONS = 3
const NUM_MAX_ELEVATORS = 4

const UDP_PRESENCE_MSG = "Pella"
const UDP_BROADCAST_ADDR = "255.255.255.255:20003"
const TCP_PORT = ":30003"

//Declare button types
type ButtonType int
const (
	BUTTON_CALL_UP = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

//Declare motor directions
type MotorDir int
const (
	DIR_DOWN = iota - 1
	DIR_STOP
	DIR_UP
)

type ElevState struct {
	is_idle bool
	direction MotorDir
	last_floor int
	destination_floor int
}

type MessageType int
const (
	STATE_UPDATE = iota
	ADD_ORDER
	DELETE_ORDER
)

type Message struct { 					//The data to be sent through a NetworkMessage
	msg_type MessageType
	state ElevState		
	pressed_button_type ButtonType		//-1 if N/A
	pressed_button_floor int			//-1 if N/A
	elevs_in_network_count int			//Used by receiver to check if sender and receiver "see" the same network, to make sure all necessary connections are made
}

type NetworkMessage struct {
	raddr string 						//The remote address we are receiving from, on form IP:port. 
	data []byte			
	length int							//Length of received data, don't care when transmitting
}

type Order struct{
	active bool; 						//Is there an order here?
	addr string; 						//Which elevator executes this order, blank for local elevator
}

/*
Coding convention:
	Functions: This_Is_A_Func()
	Variables: this_is_a_var
	Typenames: ThisIsATypename
*/
