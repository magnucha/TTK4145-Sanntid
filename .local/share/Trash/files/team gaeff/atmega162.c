#define F_CPU 4915200
#define BAUD 9600

#define set_bit(reg, bit) (reg |= (1 << bit))
#define clear_bit(reg, bit) (reg &= ~(1 << bit))

#include <avr/io.h>
#include <util/delay.h>
#include <util/setbaud.h>
#include <stdlib.h>
#include <stdio.h>
#include "OLED.h"
#include "joystick.h"
#include "XMEM.h"

void UART_Transmit(unsigned char data){
	//Wait for empty transmit buffer
	while(!(UCSR0A & (1 << UDRE0)));
	
	//Send data
	UDR0 = data;
}

unsigned char UART_Receive(){
	//Wait for data to be received
	while(!(UCSR0A & (1 << RXC0)));
	
	//Read received data
	return UDR0;
}

void UART_init(){
	UBRR0L = 31;//UBRR_VALUE;
	UCSR0B = (1 << RXEN0) | (1 << TXEN0);
	fdevopen(&UART_Transmit, &UART_Receive);
}

//Cut-off frequency slider low pass filer: ~800Hz
//Slope: -20dB/dec (generelt for RC-filtere)
//Voltage = (Vref/255) * ADC
uint8_t ADC_read(uint8_t input_select){
	volatile char *ADC = (char *) 0x1400;
	
	*ADC = input_select;
	_delay_us(40); //FIKS DETTE MED HARDWARE! (INTR-PIN)
	return ADC[0];
}


int main(void)
{
	UART_init();
	XMEM_init();
	OLED_init();
	//SRAM_test();
	
	//struct joystick joystick_init_value = joystick_calibrate();
	//struct sliders sliders_init_value = slider_calibrate();
	
	//volatile char *ext_ram = (char *) 0x1800;
	//volatile char *ADC = (char *) 0x1400;
	//volatile char *OLED_data = (char *) 0x1200;
	//volatile char *OLED_command = (char *) 0x1000;
	for (int i=0;i<8;i++) {
		OLED_clear_line(i);
	}
	OLED_pos(4,0);
	OLED_printf("Halla");
    while(1)
    {	
		
	}
}
