#define F_CPU 4915200
#include <stdlib.h>
#include <stdio.h>
#include <util/delay.h>
//#include <avr/io.h>
#include <avr/interrupt.h>
#include "joystick.h"
#include "UART.h"
#include "CAN_driver.h"

struct joystick calibration;

//Cut-off frequency slider low pass filer: ~800Hz
//Slope: -20dB/dec (generelt for RC-filtere)
//Voltage = (Vref/255) * ADC
uint8_t ADC_read(uint8_t input_select){
	volatile char *ADC = (char *) 0x1400;
	//INTR_0 = 0;
	
	*ADC = input_select;
	//printf("%d", INTR_0);
	//while(!INTR_0);
	while(PINB &(1<<PINB1));
	//printf("%d", INTR_0);
	//_delay_ms(100); //FIKS DETTE MED HARDWARE! (INTR-PIN)
	return ADC[0];
}

void joystick_init(){
	DDRB &= ~(1 << PB1); // ADC interrupt
	//Calibration saved for node 1
	calibration.x_pos = ADC_read(joystick_x_sel);
	calibration.y_pos = ADC_read(joystick_y_sel);
	//Send calibration to node 2
	struct CAN_message msg;
	msg.channel = 1;
	msg.length = 2;
	msg.id = 'c';
	joystick_to_CAN(&msg);
	CAN_transmit(msg);
}

//Voltage-angle relation: a = (U-2,5) * 40 -- Gir intervall (-100, 100)
struct joystick joystick_read(uint8_t rawInput){
	volatile struct joystick temp;
	
	temp.x_pos = ADC_read(joystick_x_sel);	//(ADC_read(joystick_x_sel)-130/*calibration.x_pos*/)*100/128;
	temp.y_pos = ADC_read(joystick_y_sel);	//(ADC_read(joystick_y_sel)-130/*calibration.y_pos)*/)*100/128;
	
	if (temp.x_pos > 240){
		temp.x_dir = RIGHT;
	}
	else if (temp.x_pos < 15){
		temp.x_dir = LEFT;
	}
	if (temp.y_pos > 240){
		temp.y_dir = UP;
	}
	else if (temp.y_pos < 15){
		temp.y_dir = DOWN;
	}
	if((temp.x_pos > 52) && (temp.x_pos < 180) && (temp.y_pos > 52) && (temp.y_pos < 180)){
		temp.x_dir = NEUTRAL;
		temp.y_dir = NEUTRAL;
	}
	if (!rawInput){
		joystick_correct(&temp);
	}
	return temp;
}

void joystick_correct(volatile struct joystick* input) {
	input->x_pos = (input->x_pos-calibration.x_pos)*100/128;
	input->y_pos = (input->y_pos-calibration.y_pos)*100/128;
}

//struct joystick joystick_getDir() {
	//struct joystick temp = joystick_read();
	//
	//if (temp.x_pos > 80){
		//temp.dir = RIGHT;
		//printf("RIGHT");
	//}
	//if (temp.x_pos < -80){
		//temp.dir = LEFT;
		//printf("LEFT");
	//}
	//if (temp.y_pos > 80){
		//temp.dir = UP;
		//printf("UP");
	//}
	//if (temp.y_pos < -80){
		//temp.dir = DOWN;
		//printf("DOWN");
	//}
	//return stick;
//}

uint8_t SW3_pressed(){
	return (PINB &(1 << PINB0)) == 1 ? 0 : 1;
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
