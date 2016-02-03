package fsm

import (
	"config"
)

func state_handler() {
	//states[0] will always be local elevator. Append when a new elevator is discovered on the network
	states := make([]elev_state, 1, NUM_MAX_ELEVATORS)
	for {
		select {
			case elevID := <-chan_getState:
				for i := 0; i<len(states); i++ {
					if states[i].ID == elevID {
						chan_state <- states[elevID]
					}
				}
			case newState := <- chan_setState:
				var i int = 0
				for i < len(states) {
					if states[i].ID == newState.ID {
						states[i] = newState
						break;
					}
				}
			case <- chan_getNumElevators:
				chan_getNumElevators <- len(states) //Is this safe, or should another channel be used?
		}
	}
}

func get_optimal_elevator(destinationFloor int) {
	
}
