#include <avr/io.h>
#include "SPI.h"

void SPI_reset(){
	PORTB &= ~(1 << PB4);
	SPI_Controller_Transceive(0xC0);
	PORTB |= (1 << PB4);
}

void SPI_init(){
	/* Set MOSI, SCK and SS output, all others input */
	DDRB = (1<<PB5)|(1<<PB7)|(1<<PB4);
	/* Enable SPI, Master, set clock rate fck/16 */
	SPCR = (1<<SPE)|(1<<MSTR)|(1<<SPR0);
	PORTB |= (1<<PB4);

	SPI_reset();
	//SPI_Transceive(0x05); // Enable interrupt (Move to CAN interface)
	//SPI_Transceive(0x2B);
	//SPI_Transceive(0x02);
	//SPI_Transceive(0x02);
}

char SPI_Transceive(char cData){
	/* Start transmission */
	char datain;
	
	PORTB &= ~(1 << PB4);
	SPDR = cData;
	/* Wait for transmission complete */
	while(!(SPSR & (1<<SPIF)));
	PORTB |= (1 << PB4);
	
	datain = SPDR;
	
	return datain;
}