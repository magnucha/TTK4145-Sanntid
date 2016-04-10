package fsm

import (
	"config"
	"hardware"
	"log"
	"queue"
	"time"
)

var ch_open_door = make(chan bool)
var ch_outgoing chan<- config.Message
var ch_order_received chan config.ButtonStruct

func Init(ch_outgoing_msg chan<- config.Message, ch_new_order chan config.ButtonStruct) {
	ch_outgoing = ch_outgoing_msg
	ch_order_received = ch_new_order
	go Open_Door(ch_open_door)
}

func Event_Reached_Floor(floor int) {
	config.Local_elev.Last_floor = floor
	hardware.Elev_Set_Floor_Indicator(floor)
	if queue.Should_Stop_On_Floor(floor) {
		hardware.Elev_Set_Motor_Direction(config.DIR_STOP)
		queue.Delete_Order(floor, ch_outgoing, true)
		ch_open_door <- true
		config.Local_elev.Is_idle = true
	}
}

func Event_Order_Received() {
	for {
		button := <-ch_order_received
		var target string
		if button.Button_type == config.BUTTON_COMMAND {
			target = config.Laddr
		} else {
			target = queue.Get_Optimal_Elev(button)
		}
		queue.Add_Order(button, target, ch_outgoing, ch_order_received)

		if target == config.Laddr {
			if config.Local_elev.Door_open {
				if queue.Should_Stop_On_Floor(config.Local_elev.Last_floor) {
					ch_open_door <- true
					time.Sleep(time.Millisecond)
					queue.Delete_Order(config.Local_elev.Last_floor, ch_outgoing, true)
				}
			} else if config.Local_elev.Is_idle {
				dir := Choose_New_Direction()
				config.Local_elev.Direction = dir
				hardware.Elev_Set_Motor_Direction(dir)
				if queue.Should_Stop_On_Floor(config.Local_elev.Last_floor) {
					ch_open_door <- true
					time.Sleep(time.Millisecond)
					queue.Delete_Order(config.Local_elev.Last_floor, ch_outgoing, true)
				} else {
					config.Local_elev.Is_idle = false
				}
			}
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

func Choose_New_Direction() config.MotorDir {
	floor := config.Local_elev.Last_floor
	dir := config.Local_elev.Direction
	if queue.Is_Empty() {
		return config.DIR_STOP
	}
	switch dir {
	case config.DIR_UP:
		if queue.Is_Order_Above(floor) {
			return config.DIR_UP
		} else {
			return config.DIR_DOWN
		}
	case config.DIR_DOWN:
		if queue.Is_Order_Below(floor) {
			return config.DIR_DOWN
		} else {
			return config.DIR_UP
		}
	case config.DIR_STOP:
		if queue.Is_Order_Above(floor) {
			return config.DIR_UP
		} else if queue.Is_Order_Below(floor) {
			return config.DIR_DOWN
		} else {
			return config.DIR_STOP
		}
	default:
		log.Fatal("Choose_Direction failed!")
		return 0
	}
}

func Open_Door(ch_open <-chan bool) {
	const duration = 2 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		time.Sleep(100 * time.Millisecond)
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
