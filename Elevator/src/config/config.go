package config

const NUM_FLOORS = 4
const NUM_BUTTONS = 3
const NUM_MAX_ELEVATORS = 4

const UDP_PRESENCE_MSG = "Pella"
const UDP_BROADCAST_ADDR = "129.241.187.255"
const UDP_BROADCAST_PORT = ":20003"

var Laddr = ""

type ButtonType int
const (
	BUTTON_CALL_UP = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

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
	ACK = iota
	STATE_UPDATE
	ADD_ORDER
	DELETE_ORDER
)

type Message struct { 					//The data to be sent through a NetworkMessage
	Raddr string `json:"-"`
	Msg_type MessageType
	State ElevState		
	Button_type ButtonType		
	Floor int			
	Elevs_in_network_count int			//Used by receiver to check if sender and receiver "see" the same network, to make sure all necessary connections are made
}

type NetworkMessage struct {
	Raddr string						//The remote address we are receiving from, on form IP:port.
	Data []byte			
	Length int							//Length of received data, don't care when transmitting
}

type Order struct{
	Active bool; 						//Is there an order here?
	Addr string; 						//Which elevator executes this order, blank for local elevator
}

var Active_elevs = make(map[string]ElevState)