package config

const NUM_FLOORS = 4
const NUM_BUTTONS = 3
const NUM_MAX_ELEVATORS = 4

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

type NetworkMessage struct {
	raddr string 	//The remote address we are receiving from, on form IP:port. 
	data []byte			
	length int			//Length of received data, don't care when transmitting
}


/*
Coding convention:
	Functions: This_Is_A_Func()
	Variables: this_is_a_var
	Typenames: ThisIsATypename
*/
