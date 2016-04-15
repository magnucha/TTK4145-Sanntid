package hardware

import (
	"config"
	"queue"
	"time"
	"log"
)

func ReadButtons(ch_button_polling chan<- config.ButtonStruct) {
	var last_floor [config.NUM_BUTTONS][config.NUM_FLOORS]int
	var button config.ButtonType
	for {
		time.Sleep(100 * time.Millisecond)
		for button = config.BUTTON_CALL_UP; button <= config.BUTTON_COMMAND; button++ {
			for floor := 0; floor < config.NUM_FLOORS; floor++ {
				value := GetButtonSignal(button, floor)
				if value != 0 && value != last_floor[button][floor] {
					ch_button_polling <- config.ButtonStruct{Button_type: button, Floor: floor}
				}
				last_floor[button][floor] = value
			}
		}
	}
}

func SetLights() {
	for {
		time.Sleep(100 * time.Millisecond)
		for floor := 0; floor < len(queue.Queue); floor++ {
			for button := 0; button < len(queue.Queue[floor]); button++ {
				if queue.Queue[floor][button].Active {
					SetButtonLamp(config.ButtonType(button), floor, 1)
				} else {
					SetButtonLamp(config.ButtonType(button), floor, 0)
				}
			}
		}
	}
}

func FloorPoller(ch_floor_poll chan<- int) {
	var current_floor int
	prev := -1
	for {
		time.Sleep(100 * time.Millisecond)
		current_floor = GetFloorSensorSignal()
		if current_floor != -1 && current_floor != prev {
			prev = current_floor
			ch_floor_poll <- current_floor
		}
	}
}

func KillDefective() {
	backup_kill := func(){
		config.UDP_BACKUP_MSG = "Kill backup"
		time.Sleep(300*time.Millisecond)
		log.Fatal("Hardware not responding! Killing elevator...")
	}
	config.Local_elev.Timer = time.AfterFunc(config.TIMEOUT_HARDWARE, backup_kill)
	config.Local_elev.Timer.Stop()
	prev_idle_status := true
	for {
		if !config.Local_elev.Is_idle && (prev_idle_status == true) {
			config.Local_elev.Timer.Reset(config.TIMEOUT_HARDWARE)
		} else if config.Local_elev.Is_idle {
			config.Local_elev.Timer.Stop()
		}
		prev_idle_status = config.Local_elev.Is_idle
		time.Sleep(100*time.Millisecond)
	}
}