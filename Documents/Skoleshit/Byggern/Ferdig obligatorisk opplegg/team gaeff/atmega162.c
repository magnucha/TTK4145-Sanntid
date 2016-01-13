#define set_bit(reg, bit) (reg |= (1 << bit))
#define clear_bit(reg, bit) (reg &= ~(1 << bit))

#define F_CPU 4915200

#include <stdlib.h>
#include <stdio.h>
//#include <avr/io.h>
#include <util/delay.h>
#include <avr/interrupt.h>

#include "oled.h"
#include "joystick.h"
#include "XMEM.h"
#include "UART.h"
#include "SPI.h"
#include "SPI_driver.h"
#include "CAN_driver.h"
#include "MCP2515_defs.h"



void INTR_init(){
	GICR |= (1 << INT0)|(1 << INT1) | (1 << INT2); // Enable external interrupt request 0 and 1
	MCUCR |= (1 << ISC01)|(1 << ISC11); // Interrupt on falling edge on INT0 and INT1
	EMCUCR &= ~(1 << ISC2); // Sets interrupt on falling edge on INT2
	SREG |= 0b10000000; // Global interrupt enable
}


int main(void)
{
	UART_init();
	XMEM_init();
	OLED_init();
	INTR_init();
	SPI_init();
	CAN_init();
	joystick_init();
	//SRAM_test();
	
	//struct joystick joystick_init_value = joystick_calibrate();
	//struct sliders sliders_init_value = slider_calibrate();
	
	//volatile char *ext_ram = (char *) 0x1800;
	//volatile char *ADC = (char *) 0x1400;
	//volatile char *OLED_data = (char *) 0x1200;
	//volatile char *OLED_command = (char *) 0x1000;
	
	//OLED_clear_screen();
	//OLED_home();
	//OLED_printf("Halla");	 
	
	//SPI_loopback();
	uint8_t menu = 0;
	//printf("%d", SPI_Read_Status());
	
	//struct CAN_message halla;
	//halla.channel = 1;
	//halla.id = 'j';
	//halla.data = "hall";
	//halla.length = 4;
	//CAN_transmit(halla);
	
	//char* msg = CAN_recieve(1);
	
	//printf("%c%c%c%c\n", msg[0], msg[1], msg[2], msg[3]);
	
	
	//free(msg);
    while(1)
    {
		menu = OLED_menu(menu);
		
		//joystick_to_CAN(&halla);
		////printf("%d%d\n", halla.data[0], halla.data[1]);
		//CAN_transmit(halla);
		
		//if(CAN_RX_buffer_full()){ //Interrupt fra CAN controller?
			//halla = CAN_recieve(CAN_RX_buffer_full()); //Les fra bufferet som trigget interrupt
			//for (int i=0;i<halla.length;i++) {
				//printf("%c", halla.data.u8[i]);
			//}
			//printf("\n");
			////free(msg);
		//}
//
		//_delay_ms(100);
		
	}
}
