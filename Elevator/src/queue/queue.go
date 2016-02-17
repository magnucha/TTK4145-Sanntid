package queue

import(
	"config"
)

var Queue = [config.NUM_FLOORS][config.NUM_BUTTONS] config.Order{};

//Delete all orders on the floor
func Delete_Order(floor int){
	for button := config.BUTTON_CALL_UP; button <= config.BUTTON_COMMAND; button++{
		Queue[floor][button].Active = false;
		Queue[floor][button].Addr = "";
	}
}

func Add_Order(button config.ButtonStruct){
	Queue[button.Floor][button.Button_type].Active = true;
	Queue[button.Floor][button.Button_type].Addr = ""; //Use cost function when we get one(i.e. Assign_Order_To_Lift())
}

func Check_Order(floor int) bool{
	var button config.ButtonType
	if(config.Local_elev.Direction == config.DIR_UP){
		button = config.BUTTON_CALL_UP
	} else{
		button = config.BUTTON_CALL_DOWN
	}
	order := Queue[floor][button]
	pick_up := order.Active /*&& (order.Addr == config.Laddr)*/ //We think this is a double check
	command := Queue[floor][config.BUTTON_COMMAND].Active
	return  pick_up || command
}