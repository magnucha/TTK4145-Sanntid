package queue

import (
	"config"
	"log"
)

var Queue = [config.NUM_FLOORS][config.NUM_BUTTONS]config.Order{}

//Assume everyone enters when we stop --> delete all orders on the floor
func Delete_Order(floor int, ch_outgoing_msg chan<- config.Message, call_is_local bool) {
	for button := config.BUTTON_CALL_UP; button <= config.BUTTON_COMMAND; button++ {
		Queue[floor][button].Active = false
		Queue[floor][button].Addr = ""
	}
	if call_is_local {
		ch_outgoing_msg <- config.Message{Msg_type: config.DELETE_ORDER, Button: config.ButtonStruct{Floor: floor}}
	}
}
 	
func Add_Order(button config.ButtonStruct) {
	Queue[button.Floor][button.Button_type].Active = true
	Queue[button.Floor][button.Button_type].Addr = "" //Use cost function when we get one(i.e. Assign_Order_To_Lift())
}

func Get_Order(button config.ButtonStruct) config.Order {
	return Queue[button.Floor][button.Button_type];
}

func Should_Stop_On_Floor(floor int) bool {
	/*
	Stopp hvis:
		1) Øverst eller nederst
		2) Ingen bestillinger lenger i samme retning (Dekker egentlig 1)
		3) Command || Call i samme retning
	*/
	dir := config.Local_elev.Direction
	
	if dir == config.DIR_UP && !Is_Order_Above(floor) {
		return true
	}
	if dir == config.DIR_DOWN && !Is_Order_Below(floor) {
		return true
	}
	
	var button_in_current_dir config.ButtonType
	if dir == config.DIR_UP || floor == 0 {
		button_in_current_dir = config.BUTTON_CALL_UP
	} else {
		button_in_current_dir = config.BUTTON_CALL_DOWN
	}

	pick_up := Queue[floor][button_in_current_dir].Active
	command := Queue[floor][config.BUTTON_COMMAND].Active
	return pick_up || command
}

func Is_Order_Above(floor int) bool {
	for floor := floor + 1; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active {// && Queue[floor][button].Addr == config.Laddr {
				return true
			}
		}
	}
	return false
}

func Is_Empty() bool{
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active {// && Queue[floor][button].Addr == config.Laddr {
				return false
			}
		}
	}
	return true
}

func Is_Order_Below(floor int) bool {
	for floor := floor - 1; floor >= 0; floor-- {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active {// && Queue[floor][button].Addr == config.Laddr {
				return true
			}
		}
	}
	return false
}