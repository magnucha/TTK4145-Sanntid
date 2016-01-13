#define F_CPU 4915200
#define BAUD 9600
#include "UART.h"
#include <avr/io.h>
#include <util/setbaud.h>

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