package hardware

import(
	"config"
	"time"
)

//Goes to floor and picks up any orders in the same direction
func Go_To_Floor(floor int, this_elev *ElevState){
	this_elev.Is_idle(0)
	if(this_elev.Last_floor < floor){
		Elev_Set_Motor_Direction(config.DIR_UP)
	}
	else if(this_elev.Last_floor == floor){
		Elev_Set_Motor_Direction(config.DIR_STOP)
	}
	else{
		Elev_Set_Motor_Direction(config.DIR_DOWN)
	}
	for(this_elev.Last_floor != Elev_Get_Floor_Sensor_Signal()){
		if(Elev_Get_Floor_Sensor_Signal() != -1){
			//Passing_Floor(Elev_Get_Floor_Sensor_Signal(), this_elev.Direction); Fiksien
		}
	}
}

func Open_Door(){
	Elev_Set_Door_Open_Lamp(1)
	time.Sleep(2*time.Second)
	Elev_Set_Door_Open_Lamp(0)
}

//Check for orders in same direction as the floor you're passing
func Passing_Floor(floor int, dir MotorDir){
	var button ButtonType
	if(dir == config.DIR_UP){
		button = BUTTON_CALL_UP
	}
	else{
		button = BUTTON_CALL_DOWN
	}

	if(queue[floor][button].addr == config.Laddr){
		Elev_Set_Motor_Direction(config.DIR_STOP)
		Open_Door()
		Delete_Order(floor, button)
		Elev_Set_Motor_Direction(dir)
	}
}