package config

const NUM_FLOORS = 4
const NUM_BUTTONS = 3

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