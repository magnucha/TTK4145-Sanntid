package queue

import(
	"config"
)

var Queue = [config.NUM_FLOORS][config.NUM_BUTTONS] config.Order{};

func Delete_Order(floor int, button config.ButtonType){
	Queue[floor][button].Active = false;
	Queue[floor][button].Addr = "";
}

func Add_Order(floor int, button config.ButtonType){
	Queue[floor][button].Active = true;
	Queue[floor][button].Addr = ""; //Use cost function when we get one(i.e. Assign_Order_To_Lift())
}
