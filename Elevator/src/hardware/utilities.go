package hardware

import(
	"config"
	"time"
	"queue"
)

//Goes to floor and picks up any orders in the same direction
/*func Go_To_Floor(floor int){
	local_elev := config.Active_elevs[config.Laddr]
	local_elev.Is_idle = false
	
	Change_Destination(floor)
	Elev_Set_Motor_Direction(local_elev.Direction)
	for(local_elev.Destination_floor != Elev_Get_Floor_Sensor_Signal()){
		if current_floor := Elev_Get_Floor_Sensor_Signal(); current_floor != -1{
			Passing_Floor(current_floor, local_elev.Direction);
		}
	}
	local_elev.Is_idle = true
	Stop_On_Floor(config.DIR_STOP)
}
*/
func Open_Door(){
	Elev_Set_Door_Open_Lamp(1)
	time.Sleep(2*time.Second)
	Elev_Set_Door_Open_Lamp(0)
}

func Stop_On_Floor() {
	Elev_Set_Motor_Direction(config.DIR_STOP)
	Open_Door()
	Elev_Set_Motor_Direction(config.Local_elev.Direction)
}

//Check for orders in same direction as the floor you're passing
/*func Passing_Floor(floor int, dir config.MotorDir){
	
	var button config.ButtonType
	if(dir == config.DIR_UP){
		button = config.BUTTON_CALL_UP
	} else{
		button = config.BUTTON_CALL_DOWN
	}

	if(queue.Queue[floor][button].Addr == config.Laddr && queue.Queue[floor][button].Active){
		Stop_On_Floor(dir)
		queue.Delete_Order(floor, button)
	}
}
*/
func Change_Destination(floor int) {
	local_elev := config.Active_elevs[config.Laddr]
	local_elev.Destination_floor = floor

	if(local_elev.Last_floor < floor){
		local_elev.Direction = config.DIR_UP
	} else if(local_elev.Last_floor == floor){
		local_elev.Direction = config.DIR_STOP
	} else{
		local_elev.Direction = config.DIR_DOWN
	}
}

func Read_Buttons(ch_button_polling chan<- config.ButtonStruct){
	var last_floor[config.NUM_BUTTONS][config.NUM_FLOORS] int
	var button config.ButtonType
	for{
		time.Sleep(100*time.Millisecond)
		for button = config.BUTTON_CALL_UP; button < config.BUTTON_COMMAND; button++ {
			for floor := 0; floor < config.NUM_FLOORS; floor++ {
				value := Elev_Get_Button_Signal(button, floor)
				if(value != 0 && value != last_floor[button][floor]){
					ch_button_polling <- config.ButtonStruct{Button_type: button, Floor: floor}
				}
				last_floor[button][floor] = value
			}
		}
	}
}

func Set_Lights() {
	for {
		time.Sleep(100*time.Millisecond)
		for floor:=0; floor<len(queue.Queue); floor++ {
			for button:=0; button<len(queue.Queue[floor]); button++ {
				if(queue.Queue[floor][button].Active){
					Elev_Set_Button_Lamp(config.ButtonType(button), floor, 1)
				} else{
					Elev_Set_Button_Lamp(config.ButtonType(button), floor, 0)
				}
			}
		}
	}
}

func Floor_Poller(ch_floor_poll chan<- int) {
	var current_floor int
	var prev int
	for {
		time.Sleep(100*time.Millisecond)
		current_floor = Elev_Get_Floor_Sensor_Signal()
		if (current_floor != -1 && current_floor != prev) {
			ch_floor_poll <- current_floor
		}
	}
}

