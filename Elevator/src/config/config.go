package config

import "time"

const NUM_FLOORS = 4
const NUM_BUTTONS = 3

const QUEUE_FILENAME = "q.txt"

var UDP_BACKUP_MSG = "The cake is a lie!"
const UDP_BROADCAST_ADDR = "129.241.187.255"
const UDP_BROADCAST_PORT = ":20003"
const UDP_ALIVE_PORT = ":20103"

var Laddr = ""

const MESSAGE_PREFIX = "Ey Billy!"

const TIMEOUT_REMOTE = 2 * time.Second
const TIMEOUT_LOCAL = time.Second
const TIMEOUT_ORDER = 10 * time.Second
const TIMEOUT_UDP = 500*time.Millisecond
const TIMEOUT_HARDWARE = 8*time.Second

type ButtonType int
const (
	BUTTON_CALL_UP = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)

type ButtonStruct struct {
	Button_type ButtonType
	Floor       int
}

type MotorDir int
const (
	DIR_DOWN = iota - 1
	DIR_STOP
	DIR_UP
)

type ElevState struct {
	Is_idle    bool
	Door_open  bool
	Direction  MotorDir
	Last_floor int
	Timer      *time.Timer `json:"-"` 
}

type MessageType int
const (
	ACK = iota
	STATE_UPDATE
	ADD_ORDER
	DELETE_ORDER
)

type Message struct {
	Raddr       string `json:"-"`
	Msg_type    MessageType
	State       ElevState
	Button      ButtonStruct
}

type NetworkMessage struct {
	Raddr  string
	Data   []byte
	Length int
}

var Active_elevs = make(map[string]*ElevState)
var Local_elev *ElevState
