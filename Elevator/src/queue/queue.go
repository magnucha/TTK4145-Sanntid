package queue

import(
	"config"
)

var queue = [config.NUM_FLOOR][config.NUM_BUTTONS] Order{};

func Delete_Order(floor int, button ButtonType){
	queue[floor][button].active = false;
	queue[floor][button].addr = "";
}

func Add_Order(floor int, button ButtonType){
	queue[floor][button].active = true;
	queue[floor][button].addr = ""; //Use cost function when we get one(i.e. Assign_Order_To_Lift())
}