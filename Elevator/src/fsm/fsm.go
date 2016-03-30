package fsm

import (
	"config"
	"queue"
	"hardware"
	"time"
)

var ch_open_door = make(chan bool)
var ch_outgoing chan<- config.Message

func FSM_Init(ch_outgoing_msg chan<- config.Message) {
	ch_outgoing = ch_outgoing_msg
	go Open_Door(ch_open_door)
}

func Event_Reached_Floor(floor int, ch_outgoing_msg chan<- config.Message) {
	config.Local_elev.Last_floor = floor
	hardware.Elev_Set_Floor_Indicator(floor)
	if queue.Should_Stop_On_Floor(floor) {
		hardware.Elev_Set_Motor_Direction(config.DIR_STOP)
		queue.Delete_Order(floor, ch_outgoing, true)
		ch_open_door <- true
		config.Local_elev.Is_idle = true
	}
}

func Event_Order_Received(button config.ButtonStruct) {
	//Add a check for cost function when implemented
	queue.Add_Order(button)
	
	if config.Local_elev.Door_open {
		if queue.Should_Stop_On_Floor(config.Local_elev.Last_floor) {
			queue.Delete_Order(config.Local_elev.Last_floor, ch_outgoing, true)
			ch_open_door <- true
		}
	} else if config.Local_elev.Is_idle {
		dir := Choose_New_Direction()
		config.Local_elev.Direction = dir
		hardware.Elev_Set_Motor_Direction(dir)
		if queue.Should_Stop_On_Floor(config.Local_elev.Last_floor) {
			ch_open_door <- true
			queue.Delete_Order(config.Local_elev.Last_floor, ch_outgoing, true)
		} else {
			config.Local_elev.Is_idle = false
		}
	}
}

func Event_Door_Closed() {
	config.Local_elev.Door_open = false
	hardware.Elev_Set_Door_Open_Lamp(0)
	dir := Choose_New_Direction()
	config.Local_elev.Direction = dir
	hardware.Elev_Set_Motor_Direction(dir)
}

func (elev *ElevState) Is_Moving_Toward(floor int) bool {
	return (elev.Direction == config.DIR_UP && floor > elev.Last_floor) || (elev.Direction == config.DIR_DOWN && floor < elev.Last_floor);
}

func Open_Door(ch_open <-chan bool) {
	const duration = 2 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		time.Sleep(100*time.Millisecond)
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
