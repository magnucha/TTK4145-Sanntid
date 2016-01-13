#define F_CPU 4915200

#include <avr/io.h>
#include <util/delay.h>
#include "SPI.h"
#include "SPI_driver.h"


void SPI_wr(char d){
	SPDR = d;
	while(!(SPSR & (1<<SPIF)));
}

char SPI_Controller_Transceive(char cData){
	/* Start transmission */
	
	SPDR = cData;
	/* Wait for transmission complete */
	while(!(SPSR & (1<<SPIF)));
	_delay_ms(1);
	
	return SPDR;
}

char SPI_read(uint8_t address){
	char temp;
	PORTB &= ~(1<<PB4);
	SPI_Controller_Transceive(0x03);
	SPI_Controller_Transceive(address);
	temp = SPI_Controller_Transceive(0xFF);
	PORTB |= (1<<PB4);
	return temp;
}

void SPI_write(uint8_t address, uint8_t data){
	PORTB &= ~(1<<PB4);
	SPI_Controller_Transceive(0x02);
	SPI_Controller_Transceive(address);
	SPI_Controller_Transceive(data);
	PORTB |= (1<<PB4);
}

void SPI_Req_To_Send(uint8_t channel){
	PORTB &= ~(1<<PB4);
	SPI_Controller_Transceive(0x80 + channel); // channel = 0b0xxx
	PORTB |= (1<<PB4);
}

uint8_t SPI_Read_Status(){
	PORTB &= ~(1<<PB4);
	return SPI_Controller_Transceive(0xA0); // kan hende ikke funker
	PORTB |= (1<<PB4);
}

void SPI_bit_modify(uint8_t address, uint8_t mask, uint8_t data){
	PORTB &= ~(1<<PB4);
	SPI_Controller_Transceive(0x05);
	SPI_Controller_Transceive(address);
	SPI_Controller_Transceive(mask);
	SPI_Controller_Transceive(data);
	PORTB |= (1<<PB4);
}

char SPI_RX_read(uint8_t address){
	char temp;
	PORTB &= ~(1<<PB4);
	SPI_Controller_Transceive(0x92);
	temp = SPI_Controller_Transceive(0);
	PORTB |= (1<<PB4);
	return temp;
}