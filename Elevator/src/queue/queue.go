package queue

import (
	"config"
	"log"
	"math"
	//"os"
	"encoding/json"
	//"io"
	"io/ioutil"
)

type Order struct{
	Active bool; 						//Is this button pressed?
	Addr string; 						//Which elevator executes this order
}

var Queue = [config.NUM_FLOORS][config.NUM_BUTTONS]Order{}

//Assume everyone enters when we stop --> delete all orders on the floor
func Delete_Order(floor int, ch_outgoing_msg chan<- config.Message, call_is_local bool) {
	for button := config.BUTTON_CALL_UP; button <= config.BUTTON_COMMAND; button++ {
		Queue[floor][button].Active = false
		Queue[floor][button].Addr = ""
	}
	if call_is_local {
		ch_outgoing_msg <- config.Message{Msg_type: config.DELETE_ORDER, Button: config.ButtonStruct{Floor: floor}}
	}
	File_Write(config.QUEUE_FILENAME)
}

func Add_Order(button config.ButtonStruct, addr string) {
	Queue[button.Floor][button.Button_type].Active = true
	Queue[button.Floor][button.Button_type].Addr = addr
	File_Write(config.QUEUE_FILENAME)
}

func Get_Order(button config.ButtonStruct) Order {
	if (button.Floor < 0 || button.Floor > 3) {
		log.Fatal("FATAL: Get_Order floor out of range! Floor is", button.Floor)
	}
	return Queue[button.Floor][button.Button_type]
}

func Should_Stop_On_Floor(floor int) bool {
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

func Reassign_Orders(addr string, ch_new_order chan<- config.ButtonStruct) {
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Addr == addr {
				ch_new_order <- config.ButtonStruct{config.ButtonType(button), floor}
			}
		}
	}
}

func File_Read(filename string) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("ERROR: Could not read queue from file! error: %s", err.Error())
		return
	}
	err = json.Unmarshal(input, &Queue)
	if err != nil {
		//Test contents of Queue! Is fatal needed?
		log.Fatal("FATAL: Could not decode file input! json error: %s", err.Error())
	}
}

func File_Write(filename string) {
	output, err := json.Marshal(Queue)
	if err != nil {
		log.Printf("ERROR: Could not encode queue! json error: %s", err.Error())
	}
	err = ioutil.WriteFile(filename, output, 0666)
	if err != nil {
		log.Printf("ERROR: Coult not write queue to file! error: %s", err.Error())
	}
}
//-----------------------------------------------------
/*
	- Idle
		- Abs avstand
	- På vei mot etasjen, plukker opp på veien
		- Abs avstand + stopp
	- På vei mot etasjen, skal passere
		- Abs avstand + 2*passeravstand + stopp
	- På vei vekk fra etasjen
		- 2*avstand til furthest + abs avstand + stopp
*/
func Calculate(addr string, button config.ButtonStruct) int {
	const COST_MOVE_ONE_FLOOR = 1
	const COST_DOOR_OPEN = 1
	const COST_STOP = 2
	elev := config.Active_elevs[addr]

	cost := int(math.Abs(float64(elev.Last_floor - button.Floor))) * COST_MOVE_ONE_FLOOR

	//Moving towards destination floor
	if (elev.Direction == config.DIR_UP && button.Floor > elev.Last_floor) || (elev.Direction == config.DIR_DOWN && button.Floor < elev.Last_floor) {
		for f := elev.Last_floor; f != button.Floor; f += int(elev.Direction) {
			for b := config.BUTTON_CALL_UP; b <= config.BUTTON_COMMAND; b++ {
				if Get_Order(config.ButtonStruct{config.ButtonType(b), f}).Addr == addr {
					cost += COST_STOP
					break
				}
			}
		}
		//Should pass the floor before serving it
		if (elev.Direction == config.DIR_UP && button.Button_type == config.BUTTON_CALL_DOWN) || (elev.Direction == config.DIR_DOWN && button.Button_type == config.BUTTON_CALL_UP) {
			furthest_floor := button.Floor
			for f := button.Floor + int(elev.Direction); f >= 0 && f < config.NUM_FLOORS; f += int(elev.Direction) {
				for b := config.BUTTON_CALL_UP; b <= config.BUTTON_COMMAND; b++ {
					if Get_Order(config.ButtonStruct{config.ButtonType(b), f}).Addr == addr {
						furthest_floor = f
						cost += COST_STOP
						break
					}
				}
			}
			cost += int(math.Abs(float64(furthest_floor - button.Floor))) * COST_MOVE_ONE_FLOOR * 2
		}
	//Is moving away from destination floor
	} else if !elev.Is_idle {
		furthest_floor := -1
		for f := elev.Last_floor; f >= 0 && f < config.NUM_FLOORS; f += int(elev.Direction) {
			for b := config.BUTTON_CALL_UP; b <= config.BUTTON_COMMAND; b++ {
				if Get_Order(config.ButtonStruct{config.ButtonType(b), f}).Addr == addr {
					furthest_floor = f
					cost += COST_STOP
					break
				}
			}
		}
		cost += int(math.Abs(float64(furthest_floor - elev.Last_floor))) * COST_MOVE_ONE_FLOOR * 2

		for f := elev.Last_floor; f != button.Floor; f -= int(elev.Direction) {
			for b := config.BUTTON_CALL_UP; b <= config.BUTTON_COMMAND; b++ {
				if Get_Order(config.ButtonStruct{config.ButtonType(b), f}).Addr == addr {
					cost += COST_STOP
					break
				}
			}
		}
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