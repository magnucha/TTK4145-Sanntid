package cost

import (
	"config"
	"queue"
	"fsm"
	"math"
	"strings"
)

/*
Calculation logic:
	- If moving toward the floor:
		"distance" + "distance moving past destination"*2 (+ stopping costs)
	- If moving away from the floor:
		"distance moving away" + "distance to order after turning" (+ stopping costs)
		== "distance moving past destination"*2 - "distance" (+ stopping costs)
*/
func Calculate(addr string, button config.ButtonStruct) int {
	const COST_MOVE_ONE_FLOOR = 1;
	const COST_DOOR_OPEN = 2;
	const COST_STOP = 3;
	elev := config.Active_elevs[addr];

	if elev.fsm.Is_Moving_Toward(button.Floor) {
		cost := abs(elev.Last_floor - button.Floor) * COST_MOVE_ONE_FLOOR;
		for f := elev.Last_floor; f != button.Floor; f += elev.Direction {
			cost += COST_STOP;
		}
	} else {
		cost := -abs(elev.Last_floor - button.Floor) * COST_MOVE_ONE_FLOOR;
	}

	//Calculate cost of activities if we pass the floor without stopping (order in opposite direction)
	if (button.ButtonType == config.BUTTON_CALL_UP && elev.Direction == config.DIR_DOWN) || (button.ButtonType == config.BUTTON_CALL_DOWN && elev.Direction == config.DIR_UP) {
		furthest_floor := -1
		for f := button.Floor; f >= 0 && f < config.NUM_FLOORS; f += elev.Direction {
			for b := config.BUTTON_CALL_UP; b <= config.BUTTON_COMMAND; b++ {
				if Queue.Get_Order(config.ButtonType{b,f}).Addr == addr {
					furthest_floor = f;
					cost += COST_STOP;
					break;
				}
			}
		}
		cost += abs(furthest_floor - button.Floor) * COST_MOVE_ONE_FLOOR * 2;
	}

	return cost;
}

func Get_Optimal_Elev(button config.ButtonStruct) string {
	optimal := ""
	lowest := 100000
	for addr,_ := range(config.Active_elevs) {
		cost := Calculate(addr,button);
		if cost < lowest {
			optimal = addr;
			lowest = cost;
		} else if cost == lowest {
			if addr < optimal {
				optimal = addr;
			}
		}
	}

	return optimal;
}