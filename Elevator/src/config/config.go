package config

//import (
//	"hardware"
//)

const NUM_FLOORS = 4
const NUM_BUTTONS = 3
const NUM_MAX_ELEVATORS = 4

const UDP_PRESENCE_MSG = "Pella"
const UDP_BROADCAST_ADDR = "129.241.187.255:20003"
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
	Is_idle bool
	Direction MotorDir
	Last_floor int
	Destination_floor int
}

type MessageType int
const (
	STATE_UPDATE = iota
	ADD_ORDER
	DELETE_ORDER
)

type Message struct { 					//The data to be sent through a NetworkMessage
	Msg_type MessageType
	State ElevState		
	Pressed_button_type ButtonType		//-1 if N/A
	Pressed_button_floor int			//-1 if N/A
	Elevs_in_network_count int			//Used by receiver to check if sender and receiver "see" the same network, to make sure all necessary connections are made
}

type NetworkMessage struct {
	Raddr string 						//The remote address we are receiving from, on form IP:port. 
	Data []byte			
	Length int							//Length of received data, don't care when transmitting
}

type Order struct{
	Active bool; 						//Is there an order here?
	Addr string; 						//Which elevator executes this order, blank for local elevator
}

/*
Coding convention:
	Functions: This_Is_A_Func()
	Variables: this_is_a_var
	Typenames: ThisIsATypename
*/
