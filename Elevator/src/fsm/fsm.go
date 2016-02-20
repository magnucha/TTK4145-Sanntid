package fsm

import (
	"config"
	"queue"
	"hardware"
	"time"
)

ch_open_door := make(chan bool)

func Event_Reached_Floor(floor int, ch_outgoing_msg chan<- config.Message) {
	config.Local_elev.Last_floor = floor
	hardware.Elev_Set_Floor_Indicator(floor)
	if queue.Check_Order(floor) {
		hardware.Elev_Set_Motor_Direction(config.DIR_STOP)
		queue.Delete_Order(floor, ch_outgoing_msg)
		ch_open_door <- true
		config.Local_elev.Is_idle = true
	}
}

func Event_Order_Received(button config.ButtonStruct) {
	//Add a check for cost function when implemented
	queue.Add_Order(button)
	
	if config.Local_elev.Door_open {
		if queue.Check_Order(button.Floor) {
			queue.Delete_Order(button.Floor)
			ch_open_door <- true
		}
	} else if config.Local_elev.Is_idle {
		dir = queue.Choose_New_Direction()
		config.Local_elev.Direction = dir
		hardware.Elev_Set_Motor_Direction(dir)
		if dir == config.DIR_STOP {
			ch_open_door <- true
		} else {
			config.Local_elev.Is_idle = false
		}
	}
}

func Event_Door_Closed() {
	config.Local_elev.Door_open = false
	hardware.Elev_Set_Door_Open_Lamp(0)
	dir := queue.Choose_New_Direction()
	config.Local_elev.Direction = dir
	hardware.Elev_Set_Motor_Direction(dir)
}

func Open_Door(ch_open <-chan bool) {
	const duration = 2 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-ch_open:
			timer.Reset(duration)
			config.Local_elev.Door_open = true
			hardware.Elev_Set_Door_Open_Lamp(1)
		case <-timer.C:
			timer.Stop()
			Event_Door_Closed()
		}
	}
}
