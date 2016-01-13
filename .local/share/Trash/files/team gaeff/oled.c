
#include "font8x8.h"
#include "OLED.h"
#define FONT_OFFSET 32
#define FONT_WIDTH 8

#include <stdarg.h>
#include <stdio.h>
#include <avr/pgmspace.h>
#include "joystick.h"

void OLED_init(){
	write_c(0xae);    // display off
	write_c(0xa1);    //segment remap
	write_c(0xda);    //common pads hardware: alternative
	write_c(0x12);
	write_c(0xc8);    //common output scan direction:com63~com0
	write_c(0xa8);    //multiplex ration mode:63
	write_c(0x3f);
	write_c(0xd5);    //display divide ratio/osc. freq. mode
	write_c(0x80);
	write_c(0x81);    //contrast control
	write_c(0x50);
	write_c(0xd9);    //set pre-charge period
	write_c(0x21);
	write_c(0x20);    //Set Memory Addressing Mode
	write_c(0x02);
	write_c(0xdb);    //VCOM deselect level mode
	write_c(0x30);
	write_c(0xad);    //master configuration
	write_c(0x00);
	write_c(0xa4);    //out follows RAM content
	write_c(0xa6);    //set normal display
	write_c(0xaf);    // display on
}

void OLED_home(){
	write_c(0x21); //set column address
	write_c(0x00); 
	write_c(0xFF); 
	write_c(0x22); //Set page address
	write_c(0x00);
	write_c(0x00);
}

void OLED_goto_line(uint8_t line){
	write_c(0x22); //Set page address
	write_c(line);
	write_c(line);
}

void OLED_clear_line(uint8_t line){
	OLED_goto_line(line);
	for(int i = 0; i < 128; i++){
		OLED_D[i] = 0;
	}
}

void OLED_clear_screen(){
	for (int i=0;i<8;i++){
		OLED_clear_line(i);
	}
}

void OLED_pos(uint8_t row,uint8_t column){
	write_c(0x21); //set column address
	write_c(column*8); //16 chars per row
	write_c(0xFF);
	write_c(0x22); //Set page address
	write_c(row);
	write_c(row);
}

void OLED_printchar(char c){
	for(uint8_t i = 0; i < FONT_WIDTH; i++){
		//write_d(font[c - FONT_OFFSET][i]);
		write_d(pgm_read_byte((void*)font + (c - FONT_OFFSET)*FONT_WIDTH + i));
	}
}

static FILE OLED_stdout = FDEV_SETUP_STREAM(OLED_printchar, NULL, _FDEV_SETUP_WRITE);

void OLED_printf(char* data, ...){
	va_list argp;
	va_start(argp, data);
	vfprintf(&OLED_stdout, data, argp);
	va_end(argp);
}

void OLED_printInvertedChar(char data){
	for (int i=0;i<FONT_WIDTH;i++) {
		//write_d(~font[data-FONT_OFFSET][i]) 
		write_d(~pgm_read_byte((void*)font + (c - FONT_OFFSET)*FONT_WIDTH + i));
	}
}

static FILE OLED_inverted_stdout = FEDEV_SETUP_STREAM(OLED_printInvertedChar, NULL, _FDEV_SETUP_WRITE);

void OLED_printfInverted(char* data, ...){
	va_list argp;
	va_start(argp, data);
	vfprintf(&OLED_inverted_stdout, data, argp);
	va_end(argp);
}

//---------------------------------------------------------------------------------------
//				START OLED MENY
//---------------------------------------------------------------------------------------

/*
Menyene er ordnet som integers:
	Main menu: 0
	Valg x i main: x
	Valg y i x = x*10+y
	Eks: Main->1->3 er meny 13
	
	Menyformat: char* items = ["Navn", "Return", "Item1", "Item2", "Item3", "Item4", "Item5", "Item6"];
*/

uint8_t OLED_menu(uint8_t menu){
	OLED_clear_screen();
	char* items;
	uint8_t numItems;
	switch (menu){
		case 0:
			items = ["Main Menu", "    ", "Play", "Highscores", "Reset scores", "Sound ON", "Sound OFF"];
			numItems = 6;
			break;
		default:
			OLED_clear();
			OLED_pos(0,0);
			OLED_printf("Menu selection");
			OLED_pos(1,0);
			OLED_printf("failed!");
			break;
	}
	
	uint8_t  line;
	uint8_t cursorPos = 2; //Starter på første valgalternativ
	while (true) {
		line = -1;
		while items[++line] {
			OLED_pos(line,line ? 2 : 0); //Add a 2 character indentation on all menu items
			(cursorPos == line) ? OLED_printfInverted(items[line]) : OLED_printf(items[line]);
		}
		while (joystick_dir().dir);	//Only lets you move one option at a time
		while (!(joystick_dir().dir || joystick_pressed()));
		if ((joystick_dir().dir == DOWN) && (cursorPos < numItems)) {
			cursorPos++;
		}
		else if ((joystick_dir().dir == UP) && (cursorPos > 1)) {
			cursorPos--;
		}
		else if (joystick_pressed()) {
			return (cursorPos == 1) ? floor(menu/10) : menu*10 + cursorPos;
		}
	}
}











