package fsm

import (
	"config"
	"hardware"
	"log"
	"queue"
	"time"
)

var ch_door_open = make(chan bool)
var ch_outgoing chan<- config.Message
var ch_order_received chan config.ButtonStruct

func Init(ch_outgoing_msg chan<- config.Message, ch_new_order chan config.ButtonStruct) {
	ch_outgoing = ch_outgoing_msg
	ch_order_received = ch_new_order
	go OpenDoor(ch_door_open)
}

func EventReachedFloor(floor int) {
	config.Local_elev.Last_floor = floor
	hardware.SetFloorIndicator(floor)
	if queue.ShouldStopOnFloor(floor) {
		config.Local_elev.Is_idle = true
		hardware.SetMotorDirection(config.DIR_STOP)
		queue.DeleteOrder(floor, ch_outgoing, true)
		ch_door_open <- true
	}
}

func EventOrderReceived() {
	for {
		button := <-ch_order_received
		var target string
		if button.Button_type == config.BUTTON_COMMAND {
			target = config.Laddr
		} else {
			target = queue.GetOptimalElev(button)
		}
		queue.AddOrder(button, target, ch_outgoing, ch_order_received)

		if target == config.Laddr {
			if config.Local_elev.Door_open {
				if queue.ShouldStopOnFloor(config.Local_elev.Last_floor) {
					ch_door_open <- true
					time.Sleep(time.Millisecond)
					queue.DeleteOrder(config.Local_elev.Last_floor, ch_outgoing, true)
				}
			} else if config.Local_elev.Is_idle {
				dir := ChooseNewDirection()
				config.Local_elev.Direction = dir
				hardware.SetMotorDirection(dir)
				if queue.ShouldStopOnFloor(config.Local_elev.Last_floor) {
					ch_door_open <- true
					time.Sleep(time.Millisecond)
					queue.DeleteOrder(config.Local_elev.Last_floor, ch_outgoing, true)
				} else {
					config.Local_elev.Is_idle = false
				}
			}
		}
	}
}

func ChooseNewDirection() config.MotorDir {
	floor := config.Local_elev.Last_floor
	dir := config.Local_elev.Direction
	if queue.IsEmpty() {
		return config.DIR_STOP
	}
	switch dir {
	case config.DIR_UP:
		if queue.IsOrderAbove(floor) {
			return config.DIR_UP
		} else {
			return config.DIR_DOWN
		}
	case config.DIR_DOWN:
		if queue.IsOrderBelow(floor) {
			return config.DIR_DOWN
		} else {
			return config.DIR_UP
		}
	case config.DIR_STOP:
		if queue.IsOrderAbove(floor) {
			return config.DIR_UP
		} else if queue.IsOrderBelow(floor) {
			return config.DIR_DOWN
		} else {
			return config.DIR_STOP
		}
	default:
		log.Fatal("Choose_Direction failed!")
		return 0
	}
}

func OpenDoor(ch_open <-chan bool) {
	const duration = 2 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case <-ch_open:
			timer.Reset(duration)
			config.Local_elev.Door_open = true
			hardware.SetDoorOpenLamp(1)
		case <-timer.C:
			timer.Stop()
			CloseDoor()
		}
	}
}

func CloseDoor() {
	config.Local_elev.Door_open = false
	hardware.SetDoorOpenLamp(0)
	dir := ChooseNewDirection()
	config.Local_elev.Direction = dir
	hardware.SetMotorDirection(dir)
}