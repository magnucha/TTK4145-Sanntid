package config

const NUM_FLOORS = 4
const NUM_BUTTONS = 3
const NUM_MAX_ELEVATORS = 4

//Declare button types
type button_type int
const (
	BUTTON_CALL_UP = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

//Declare motor directions
type motor_dir int
const (
	DIR_DOWN = iota - 1
	DIR_STOP
	DIR_UP
)

struct elev_state {
	var ID int
	var direction motor_dir
	var lastFloor int
	var destinationFloor int
	
}

type NetworkMessage struct {
	raddr string 	//The remote address we are receiving from, on form IP:port. 
	data []byte			
	length int			//Length of received data, don't care when transmitting
}

TCP_connections := make([]conn.TCPConn,1)

type Order struct{
	active bool; //Is this button pressed?
	addr string; //Which elevator executes this order
}



/*
Are bi-direction channels safe?

Channel description:
	chan_getState int: Send the ID of the elevator
	chan_setState elev_state
	chan_sendState elev_state

*/
