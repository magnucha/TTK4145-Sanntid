package hardware

import (
    "config"
    "log"
)


const MOTOR_SPEED = 2800

var lamp_channel_matrix = [config.NUM_FLOORS][config.NUM_BUTTONS] int {
    {LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
    {LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
    {LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
    {LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
};


var button_channel_matrix = [config.NUM_FLOORS][config.NUM_BUTTONS] int {
    {BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
    {BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
    {BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
    {BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
};



func Elev_Init() bool {
    if !IO_Init() {
        log.Printf("Unable to initialize elevator hardware!")
        return false
    }
    
    //Disable lights in all buttons
    var button config.ButtonType
    for floor := 0; floor < config.NUM_FLOORS; floor++ {   
        for button = 0; button < config.NUM_BUTTONS; button++ {
            Elev_Set_Button_Lamp(button, floor, 0)
        }
    }
    Elev_Set_Stop_Lamp(0)
    Elev_Set_Door_Open_Lamp(0)
    Elev_Set_Floor_Indicator(0)

    Elev_Set_Motor_Direction(config.DIR_UP)
    for(Elev_Get_Floor_Sensor_Signal() == -1){}
    Elev_Set_Motor_Direction(config.DIR_STOP)

    config.Active_elevs[config.Laddr] = &config.ElevState{Is_idle: true, Direction: config.DIR_STOP, Last_floor: Elev_Get_Floor_Sensor_Signal(), Destination_floor: -1}


    return true
}

func Elev_Set_Motor_Direction(dirn config.MotorDir) {
    if (dirn == 0){
        IO_Write_Analog(MOTOR, 0)
    } else if (dirn > 0) {
        IO_Clear_Bit(MOTORDIR)
        IO_Write_Analog(MOTOR, MOTOR_SPEED)
    } else if (dirn < 0) {
        IO_Set_Bit(MOTORDIR)
        IO_Write_Analog(MOTOR, MOTOR_SPEED)
    }
}


func Elev_Set_Button_Lamp(button config.ButtonType, floor int, value int) {
    if floor >= 0 && floor < config.NUM_FLOORS && button >= 0 && button < config.NUM_BUTTONS {
        if (value != 0) {
            IO_Set_Bit(lamp_channel_matrix[floor][button])
        } else {
             IO_Clear_Bit(lamp_channel_matrix[floor][button])
        }
    }
}


func Elev_Set_Floor_Indicator(floor int) {
    if !(floor >= 0 && floor < config.NUM_FLOORS) {
        log.Printf("Floor indicator: Invalid floor")
        return
    }
    
    // Binary encoding. One light must always be on.
    if (floor & 0x02 != 0) {
        IO_Set_Bit(LIGHT_FLOOR_IND1)
    } else {
        IO_Clear_Bit(LIGHT_FLOOR_IND1)
    }    

    if (floor & 0x01 != 0) {
        IO_Set_Bit(LIGHT_FLOOR_IND2)
    } else {
        IO_Clear_Bit(LIGHT_FLOOR_IND2)
    }    
}


func Elev_Set_Door_Open_Lamp(value int) {
    if value != 0 {
        IO_Set_Bit(LIGHT_DOOR_OPEN)
    } else {
        IO_Clear_Bit(LIGHT_DOOR_OPEN)
    }
}


func Elev_Set_Stop_Lamp(value int) {
    if value != 0 {
        IO_Set_Bit(LIGHT_STOP)
    } else {
        IO_Clear_Bit(LIGHT_STOP)
    }
}



func Elev_Get_Button_Signal(button config.ButtonType, floor int) int {
    if (IO_Read_Bit(button_channel_matrix[floor][button])) {
        return 1
    } else {
        return 0
    }
}


func Elev_Get_Floor_Sensor_Signal() int {
    if (IO_Read_Bit(SENSOR_FLOOR1)) {
        return 0
    } else if (IO_Read_Bit(SENSOR_FLOOR2)) {
        return 1
    } else if (IO_Read_Bit(SENSOR_FLOOR3)) {
        return 2
    } else if (IO_Read_Bit(SENSOR_FLOOR4)) {
        return 3
    } else {
        return -1
    }
}


func Elev_Get_Stop_Signal() bool {
    return IO_Read_Bit(STOP)
}


func Elev_Get_Obstruction_Signal() bool {
    return IO_Read_Bit(OBSTRUCTION)
}
