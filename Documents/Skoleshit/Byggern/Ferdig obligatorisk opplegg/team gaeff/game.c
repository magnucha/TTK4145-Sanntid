/*
 * game.c
 *
 * Created: 22.10.2015 10:42:06
 *  Author: magnucha
 */ 
#include "CAN_driver.h"
#include <avr/interrupt.h>
#include <util/delay.h>
#include "oled.h"
//#include <avr/io.h>


volatile uint8_t solenoid_pressed = 0;
volatile uint16_t score = 0;

ISR(INT2_vect){
	solenoid_pressed = 1;
}

ISR(timer0_vect){
	score++;
	oled_clear_screen();
	oled_pos(1,0);
	oled_printf("Score: ")
	oled_printf(score);
}

void scoreTimer_start() { //Score: TIMER0
	TCCR0 |= (1<<CS02)|(1<<CS00); //Timer clock = F_CPU/1024 = 4809Hz --> score/sec = 19
	TIFR &= ~(1<<TOV0); //Clear overflow interrupt flag
	TIMSK |= (1<<TOIE0); //Enable overflow interrupt
}

void scoreTimer_stop() {
	TIMSK &= ~(1<<TOIE0); //Disable overflow interrupt
}


void play_game_backup(){
	struct CAN_message received;
	
	struct CAN_message solenoid;
	solenoid.channel = 1;
	solenoid.id = 's';
	//solenoid.length = 1;
	
	struct CAN_message joystick;
	joystick.channel = 1;
	joystick.id = 'j';
	
	struct CAN_message joystick_old;
	joystick_old.data.u8[0] = -1;
	joystick_old.data.u8[1] = -1;
	
	uint8_t game_over = 0;
	
	scoreTimer_start();
	
	while(!game_over){
		if(solenoid_pressed){
			CAN_transmit(solenoid);
			printf("solenoid_pressed\n");
			_delay_ms(20);
			solenoid_pressed = 0;
		}
		joystick_to_CAN(&joystick);
		
		//I stedet for dette, bruk en timer til å sette samplingfrekvens
		if(abs(joystick_old.data.u8[0]-joystick.data.u8[0]) > 1 || abs(joystick_old.data.u8[1]-joystick.data.u8[1]) > 10) { //Only transmit if there is a noticeable difference in input
			CAN_transmit(joystick);
			joystick_old.data.u8[0] = joystick.data.u8[0];
			joystick_old.data.u8[1] = joystick.data.u8[1];
		}
		
		if(CAN_RX_buffer_full()){ //Interrupt fra CAN controller?
			received = CAN_recieve(CAN_RX_buffer_full()); //Les fra bufferet som trigget interrupt
			if(received.id == 'g') {
				game_over = 1;
				scoreTimer_stop();
			}
		}
		_delay_ms(10);
	}
	oled_clear_screen();
	oled_home();
	oled_printf("Game over!");
	oled_pos(2,0);
	oled_printf("Your score was");
	oled_pos(3,0);
	oled_printf(score);
	//Flash score to SRAM!
	
	_delay_ms(5000);
}

/*Spill-flyt
	- Trykk play game
	- Send joystick calibration
	- GAME LOOP
		- Send joystick hver gang den endres
		- hvis knapp -> Send solenoid trigger beskjed
		- Hvis IR brytes, send game stop beskjed (data = score)
	- Stopp motor
	- Print score
	- Legg score til highscoreliste (lagret i RAM)
	- Returner til main menu






















