/*
 * CAN_driver.c
 *
 * Created: 17.09.2015 15:24:03
 *  Author: oyvindrk
 */ 
#include <stdio.h>
//#include <avr/io.h>
#include <avr/interrupt.h>
#include "CAN_driver.h"
#include "SPI.h"
#include "MCP2515_defs.h"
#include "joystick.h"

volatile uint8_t RX_buffer_0;
volatile uint8_t RX_buffer_1;

#include <stdlib.h>

ISR(INT0_vect){
	RX_buffer_0 = 1;
	printf("INT0\n");
}

//ISR(INT1_vect){
	//RX_buffer_1 = 1;
	//printf("INT1\n");
	//
//}

void CAN_init(){
	SPI_reset();
	//Sett til loopback mode
	SPI_bit_modify(CANCTRL,0xE0,MODE_NORMAL);
	//Setter mask filters off, recieve any message
	SPI_bit_modify(RXB0CTRL, 0b01100000, 0xFF);
	//Aktiverer recieve buffer interrupts
	SPI_write(BFPCTRL, 0b00001111);
	
	RX_buffer_0 = 0;
	RX_buffer_1 = 0;
}

void CAN_transmit(struct CAN_message msg){
	uint8_t reg;
	switch(msg.channel){
		case 1:
			reg = TXB0SIDH;
			break;
		case 2:
			reg = TXB1SIDH;
			break;
		case 3:
			reg = TXB2SIDH;
			break;
		default:
			printf("CAN: Invalid channel");
			return;
	}
	
	SPI_write(reg, msg.id); //Bruker SID(3:10) for å kunne ha en char som ID
	SPI_write(++reg, (0b000 << 5)); // Setter de 3 LSB i ID = 0
	reg += 3; //Fordi vi ikke har extendedID
	
	SPI_write(reg, msg.length); //Setter length bits og RTR = 0 (Data Frame)
	
	for(uint8_t i = 0; i < msg.length;i++){
		SPI_write(++reg, msg.data.u8[i]);
	}
	
	SPI_Req_To_Send(msg.channel);
}

uint8_t CAN_RX_buffer_full() {
	//switch(channel){
		//case 0:
			//return (RX_buffer_0 | RX_buffer_1);
		//case 1:
			//return RX_buffer_0;
		//case 2:
			//return RX_buffer_1;
	//}
	if(RX_buffer_0 | RX_buffer_1){
		return RX_buffer_0 ? 1 : 2;
	}
	return 0;
}

struct CAN_message CAN_recieve(uint8_t channel){
	struct CAN_message temp;
	temp.length = SPI_read(RXB0DLC);
	
	for(uint8_t i = 0; i < temp.length; i++){
		temp.data.u8[i] = SPI_read(RXB0D0 + i);
	}
	SPI_RX_read(0);
	RX_buffer_0 = 0;
	RX_buffer_1 = 0;
	return temp;
}


void joystick_to_CAN(struct CAN_message* msg){
	volatile struct joystick temp;
	temp = joystick_read(1);
	
	msg->data.u8[0] = temp.x_pos;
	msg->data.u8[1] = temp.y_pos;
	
	msg->length = 2;
}
