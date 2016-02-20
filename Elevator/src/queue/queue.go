package queue

import (
	"config"
)

var Queue = [config.NUM_FLOORS][config.NUM_BUTTONS]config.Order{}

//Assume everyone enters when we stop --> delete all orders on the floor
func Delete_Order(floor int, ch_outgoing_msg chan<- config.Message) {
	for button := config.BUTTON_CALL_UP; button <= config.BUTTON_COMMAND; button++ {
		Queue[floor][button].Active = false
		Queue[floor][button].Addr = ""
	}
	ch_outgoing_msg <- config.Message{Msg_type: config.DELETE_ORDER, Button: config.ButtonStruct{Floor: floor}}
}

func Add_Order(button config.ButtonStruct) {
	Queue[button.Floor][button.Button_type].Active = true
	Queue[button.Floor][button.Button_type].Addr = "" //Use cost function when we get one(i.e. Assign_Order_To_Lift())
}

func Check_Order(floor int) bool {
	var button config.ButtonType
	if config.Local_elev.Direction == config.DIR_UP {
		button = config.BUTTON_CALL_UP
	} else {
		button = config.BUTTON_CALL_DOWN
	}
	pick_up := Queue[floor][button].Active
	command := Queue[floor][config.BUTTON_COMMAND].Active
	return pick_up || command
}

func Is_Order_Above(floor int) bool {
	for floor := floor+1; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active && Queue[floor][button].Addr == config.Laddr {
				return true
			}
		}
	}
	return false
}

func Is_Order_Below(floor int) bool {
	for floor := floor-1; floor >= 0); floor-- {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active && Queue[floor][button].Addr == config.Laddr {
				return true
			}
		}
	}
	return false
}

func Choose_New_Direction() config.MotorDir {
	floor := config.Local_elev.Last_floor
	dir := config.Local_elev.Direction
	if Check_Order(floor) {
		return config.DIR_STOP
	}
	switch dir {
		case dir == config.DIR_UP:
			if Is_Order_Above(floor) {
				return config.DIR_UP
			} else {
				return config.DIR_DOWN
			}
		case dir == config.DIR_DOWN:
			if Is_Order_Below(floor) {
				return config.DIR_DOWN
			}else {
				return config.DIR_UP
			}
		case dir == config.DIR_STOP:
			if Is_Order_Above(floor) {
				return config.DIR_UP
			} else if Is_Order_Below(floor) {
				return config.DIR_DOWN
			} else {
				return config.DIR_STOP
			}
	}
}
