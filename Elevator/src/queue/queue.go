package queue

import (
	"config"
	"log"
	"math"
	"os"
	"encoding/json"
	//"io"
	"io/ioutil"
	"time"
)

type Order struct {
	Active bool
	Addr   string
	Timer  *time.Timer `json:"-"`
}

var Queue = [config.NUM_FLOORS][config.NUM_BUTTONS]Order{}

//Assume everyone enters when we stop --> delete all orders on the floor
func DeleteOrder(floor int, ch_outgoing_msg chan<- config.Message, call_is_local bool) {
	last_button := config.BUTTON_CALL_DOWN
	if call_is_local {
		last_button = config.BUTTON_COMMAND
		ch_outgoing_msg <- config.Message{Msg_type: config.DELETE_ORDER, Button: config.ButtonStruct{Floor: floor}}
	}
	for button := config.BUTTON_CALL_UP; button <= last_button; button++ {
		Queue[floor][button].Active = false
		Queue[floor][button].Addr = ""
		if Queue[floor][button].Timer != nil {
			Queue[floor][button].Timer.Stop()
		}
	}
	FileWrite(config.QUEUE_FILENAME)
}

func AddOrder(button config.ButtonStruct, addr string, ch_new_order chan<- config.ButtonStruct) {
	Queue[button.Floor][button.Button_type].Active = true
	Queue[button.Floor][button.Button_type].Addr = addr
	order_timeout := func() {
		DeleteOrder(button.Floor, nil, false)
		ch_new_order <- button
	}
	Queue[button.Floor][button.Button_type].Timer = time.AfterFunc(config.TIMEOUT_ORDER, order_timeout)
	FileWrite(config.QUEUE_FILENAME)
}

func GetOrder(button config.ButtonStruct) Order {
	if button.Floor < 0 || button.Floor > 3 {
		log.Fatal("FATAL: GetOrder floor out of range! Floor is ", button.Floor)
	}
	return Queue[button.Floor][button.Button_type]
}

func ShouldStopOnFloor(floor int) bool {
	dir := config.Local_elev.Direction

	if dir == config.DIR_UP && !IsOrderAbove(floor) {
		return true
	}
	if dir == config.DIR_DOWN && !IsOrderBelow(floor) {
		return true
	}
	if dir == config.DIR_STOP && GetOrder(config.ButtonStruct{config.BUTTON_CALL_UP, floor}).Active {
		return true
	}

	var relevant_button config.ButtonType
	if dir == config.DIR_UP || floor == 0 {
		relevant_button = config.BUTTON_CALL_UP
	} else {
		relevant_button = config.BUTTON_CALL_DOWN
	}

	pick_up := Queue[floor][relevant_button].Active
	command := Queue[floor][config.BUTTON_COMMAND].Active
	return pick_up || command
}

func IsOrderAbove(floor int) bool {
	for floor := floor + 1; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active && Queue[floor][button].Addr == config.Laddr {
				return true
			}
		}
	}
	return false
}

func IsEmpty() bool {
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active && Queue[floor][button].Addr == config.Laddr {
				return false
			}
		}
	}
	return true
}

func IsOrderBelow(floor int) bool {
	for floor := floor - 1; floor >= 0; floor-- {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Active && Queue[floor][button].Addr == config.Laddr {
				return true
			}
		}
	}
	return false
}

func ReassignOrders(addr string, ch_new_order chan<- config.ButtonStruct) {
	for floor := 0; floor < config.NUM_FLOORS; floor++ {
		for button := 0; button < config.NUM_BUTTONS; button++ {
			if Queue[floor][button].Addr == addr {
				ch_new_order <- config.ButtonStruct{config.ButtonType(button), floor}
			}
		}
	}
}

func FileRead(filename string) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("ERROR: Could not read queue from file! error: %s", err.Error())
		return
	}
	err = json.Unmarshal(input, &Queue)
	if err != nil {
		if _, err := os.Create(config.QUEUE_FILENAME); err != nil {
			log.Fatal("FATAL: Could not create queue file!")
		}
		log.Printf("ERROR: Could not decode file input! Queue has been reset..")
	}
}

func FileWrite(filename string) {
	output, err := json.Marshal(Queue)
	if err != nil {
		log.Printf("ERROR: Could not encode queue!")
	}
	err = ioutil.WriteFile(filename, output, 0666)
	if err != nil {
		log.Printf("ERROR: Could not write queue to file! error: %s", err.Error())
	}
}

func Calculate(addr string, button config.ButtonStruct) int {
	const COST_MOVE_ONE_FLOOR = 1
	const COST_DOOR_OPEN = 1
	const COST_STOP = 2
	elev := config.Active_elevs[addr]

	cost := int(math.Abs(float64(elev.Last_floor-button.Floor))) * COST_MOVE_ONE_FLOOR

	//Moving towards destination floor
	if (elev.Direction == config.DIR_UP && button.Floor > elev.Last_floor) || (elev.Direction == config.DIR_DOWN && button.Floor < elev.Last_floor) {
		for f := elev.Last_floor; f != button.Floor; f += int(elev.Direction) {
			for b := config.BUTTON_CALL_UP; b < config.BUTTON_COMMAND; b++ {
				if GetOrder(config.ButtonStruct{config.ButtonType(b), f}).Addr == addr {
					cost += COST_STOP
					break
				}
			}
		}
		//Should pass the floor before serving it
		if (elev.Direction == config.DIR_UP && button.Button_type == config.BUTTON_CALL_DOWN) || (elev.Direction == config.DIR_DOWN && button.Button_type == config.BUTTON_CALL_UP) {
			furthest_floor := button.Floor
			for f := button.Floor + int(elev.Direction); f >= 0 && f < config.NUM_FLOORS; f += int(elev.Direction) {
				for b := config.BUTTON_CALL_UP; b < config.BUTTON_COMMAND; b++ {
					if GetOrder(config.ButtonStruct{config.ButtonType(b), f}).Addr == addr {
						furthest_floor = f
						cost += COST_STOP
						break
					}
				}
			}
			cost += int(math.Abs(float64(furthest_floor-button.Floor))) * COST_MOVE_ONE_FLOOR * 2
		}
	//Is moving away from destination floor
	} else if !elev.Is_idle {
		furthest_floor := -1
		for f := elev.Last_floor; f >= 0 && f < config.NUM_FLOORS; f += int(elev.Direction) {
			for b := config.BUTTON_CALL_UP; b < config.BUTTON_COMMAND; b++ {
				if GetOrder(config.ButtonStruct{config.ButtonType(b), f}).Addr == addr {
					furthest_floor = f
					cost += COST_STOP
					break
				}
			}
		}
		cost += int(math.Abs(float64(furthest_floor-elev.Last_floor))) * COST_MOVE_ONE_FLOOR * 2

		for f := elev.Last_floor; f != button.Floor; f -= int(elev.Direction) {
			for b := config.BUTTON_CALL_UP; b < config.BUTTON_COMMAND; b++ {
				if GetOrder(config.ButtonStruct{config.ButtonType(b), f}).Addr == addr {
					cost += COST_STOP
					break
				}
			}
		}
	}
	return cost
}

func GetOptimalElev(button config.ButtonStruct) string {
	optimal := ""
	lowest := 100000
	for addr, _ := range config.Active_elevs {
		cost := Calculate(addr, button)
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
