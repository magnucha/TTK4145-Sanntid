/*
 * joy_slide.c
 *
 * Created: 10.09.2015 17:26:18
 *  Author: magnucha
 */ 
#include <stdlib.h>
#include <avr/io.h>
#include "joystick.h"

struct joystick joystick_calibrate(){
	struct joystick temp;
	temp.x_pos = ADC_read(joystick_x_sel);
	temp.y_pos = ADC_read(joystick_y_sel);
	return temp;
}

//Voltage-angle relation: a = (U-2,5) * 40 -- Gir intervall (-100, 100)
struct joystick joystick_read(struct joystick calibration){
	struct joystick temp;
	
	temp.x_pos = (ADC_read(joystick_x_sel)-calibration.x_pos)*100/128;
	temp.y_pos = (ADC_read(joystick_y_sel)-calibration.y_pos)*100/128;
	
	return temp;
}

struct joystick joystick_dir(struct joystick calibration) {
	struct joystick stick = joystick_read(calibration);
	
	switch (stick){
		case (stick.x_pos > 95):
			stick.dir = RIGHT;
			break;
		case (stick.x_pos < -95):
			stick.dir = LEFT;
			break;
		case (stick.y_pos > 95):
			stick.dir = UP;
			break;
		case (stick.y_pos < -95):
			stick.dir = DOWN;
			break;
		default:
			stick.dir = NEUTRAL;
	}
	return stick;
}

bool joystick_pressed(){//PA0 has to be an input, and needs a pull-up resistor
	return PA0 ? false : true;
}

struct sliders slider_calibrate(){
	struct sliders temp;
	uint8_t slider_left_sel = 0x06;
	uint8_t slider_right_sel = 0x07;
	temp.left_slider = ADC_read(slider_left_sel);
	temp.right_slider = ADC_read(slider_right_sel);
	return temp;
}

struct sliders sliders_read(){
	uint8_t slider_left_sel = 0x06;
	uint8_t slider_right_sel = 0x07;
	struct sliders temp;
	
	temp.left_slider = ADC_read(slider_left_sel) * 100/256;
	temp.right_slider = ADC_read(slider_right_sel) * 100/256;
	
	return temp;
}
