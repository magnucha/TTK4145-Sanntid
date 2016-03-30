package queue

import (
	"config"
	"log"
	"math"
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

func Add_Order(button config.ButtonStruct, addr string) {
	Queue[button.Floor][button.Button_type].Active = true
	Queue[button.Floor][button.Button_type].Addr = addr
}

func Get_Order(button config.ButtonStruct) config.Order {
	return Queue[button.Floor][button.Button_type]
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
			if Queue[floor][button].Active && Queue[floor][button].Addr == config.Laddr {
				return true
			}
		}
	}
	return false
}

func Is_Empty() bool {
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active && Queue[floor][button].Addr == config.Laddr {
				return false
			}
		}
	}
	return true
}

func Is_Order_Below(floor int) bool {
	for floor := floor - 1; floor >= 0; floor-- {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active && Queue[floor][button].Addr == config.Laddr {
				return true
			}
		}
	}
	return false
}

//-----------------------------------------------------
func Calculate(addr string, button config.ButtonStruct) int {
	const COST_MOVE_ONE_FLOOR = 1
	const COST_DOOR_OPEN = 2
	const COST_STOP = 3
	var cost int
	elev := config.Active_elevs[addr]

	if elev.Is_idle {
		return int(math.Abs(float64(elev.Last_floor - button.Floor)))
	}

	if (elev.Direction == config.DIR_UP && button.Floor > elev.Last_floor) || (elev.Direction == config.DIR_DOWN && button.Floor < elev.Last_floor) {
		cost = int(math.Abs(float64(elev.Last_floor-button.Floor)) * COST_MOVE_ONE_FLOOR)
		for f := elev.Last_floor; f != button.Floor; f += int(elev.Direction) {
			cost += COST_STOP
		}
	} else {
		cost = int(-math.Abs(float64(elev.Last_floor-button.Floor)) * COST_MOVE_ONE_FLOOR)
	}

	//Calculate cost of activities if we pass the floor without stopping (order in opposite direction)
	if (button.Button_type == config.BUTTON_CALL_UP && elev.Direction == config.DIR_DOWN) || (button.Button_type == config.BUTTON_CALL_DOWN && elev.Direction == config.DIR_UP) {
		furthest_floor := -1
		for f := button.Floor; f >= 0 && f < config.NUM_FLOORS; f += int(elev.Direction) {
			for b := config.BUTTON_CALL_UP; b <= config.BUTTON_COMMAND; b++ {
				if Get_Order(config.ButtonStruct{config.ButtonType(b), f}).Addr == addr {
					furthest_floor = f
					cost += COST_STOP
					break
				}
			}
		}
		cost += int(math.Abs(float64(furthest_floor-button.Floor)) * COST_MOVE_ONE_FLOOR * 2)
	}
	return cost
}

func Get_Optimal_Elev(button config.ButtonStruct) string {
	optimal := ""
	lowest := 100000
	for addr, _ := range config.Active_elevs {
		cost := Calculate(addr, button)
		log.Printf("COST: Lowest=%d, Cost=%d, Addr=%s\n", lowest, cost, addr)
		if cost < lowest {
			optimal = addr
			lowest = cost
		} else if cost == lowest {
			if addr < optimal {
				optimal = addr
			}
		}
	}
	return optimal
}
