
void OLED_printInvertedChar(char data){ //Data according to ASCII table
	for (int i=0;i<FONT_WIDTH;i++) {
		write_d(~font[data-FONT_OFFSET][i]) //Invert the font, and write
	}
}
void OLED_printInverted(char* data){
	int i = -1;
	while (data[++i]){
		OLED_printInvertedChar(data[i]);
	}
}
/*
Menyene er ordnet som integers:
	Main menu: 0
	Valg x i main: x
	Valg y i x = x*10+y
	Eks: Main->1->3 er meny 13
*/

uint8_t OLED_menu(uint8_t menu){
	OLED_pos(0,0);
	switch (menu){
		case 0:
			char* items = ["Main Menu", "Return", "Item1", "Item2", "Item3", "Item4", "Item5"];
			break;
		default:
			OLED_clear();
			OLED_pos(0,0);
			OLED_print("Menu selection");
			OLED_pos(1,0);
			OLED_print("failed!");
			break;
	}
	
	uint8_t  line;
	uint8_t cursorPos = 1;
	while (true) {
		line = 0;
		while items[++line] {
			line ? OLED_pos(line,2); //Add a 2 character indentation
			(cursorPos == line) ? OLED_printInverted(items[line]) : OLED_print(items[line]);
		}
		while (!(joystick_dir() || joystick_pressed()));
		if (joystick_dir().dir == DOWN) {
			cursorPos++;
		}
		else if (joystick_dir().dir == UP) {
			cursorPos--;
		}
		else if (joystick_pressed()) {
			return (cursorPos == 1) ? floor(menu/10) : menu*10 + cursorPos;
		}
	}
}



struct joystick joystick_dir() {
	struct joystick stick;
	switch (stick.pos(joystick_x_sel){
		case >95:
			stick.dir = RIGHT;
			break;
		case <-95:
			stick.dir = LEFT;
			break;
		default:
			stick.dir = NEUTRAL;
	}
	switch(stick.pos(jostick_y_sel){
		case >95:
			stick.dir = UP;
			break;
		case <-95:
			stick.dir = DOWN;
			break;
		default:
			stick.dir = NEUTRAL;
	}
	return stick;
}








